/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/dsa"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/be-io/mesh/plugin/proxy/tbucket"
	tx "github.com/traefik/traefik/v2/pkg/server/router/tcp"
	"github.com/traefik/traefik/v2/pkg/server/service"
	"github.com/traefik/traefik/v2/pkg/tcp"
	px "golang.org/x/net/proxy"
	"io"
	"net"
	"strings"
	"time"
)

func init() {
	var _ tcp.WriteCloser = new(rateLimitConnection)
	var _ net.Conn = new(rateLimitConnection)
	var _ prsim.Listener = limiters
	tx.Provide(new(tcpRateLimiter))
	service.ProvideDialer(new(httpRateLimiters))
	macro.Provide(prsim.IListener, limiters)
}

const safeBandwidth = 1024 * 1024 * 1 / 8

var limiters = &rateLimiters{streams: dsa.NewStringMap[*limitStream]()}

type httpRateLimiter struct {
	serverName string
	nodeId     string
	dialer     *net.Dialer
}

func (that *httpRateLimiter) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	if !proxy.RateLimitEnable {
		return that.dialer.DialContext(ctx, network, address)
	}

	stream := limiters.MatchStreamLimiter(strings.ToUpper(that.nodeId))
	if nil == stream {
		log.Info0("HTTP connect out %s without rate limit", that.nodeId)
		return that.dialer.DialContext(ctx, network, address)
	}
	log.Info0("HTTP connect out %s(%dMbps/s-%dMbps/s)", that.nodeId, stream.upstream, stream.downstream)
	conn, err := that.dialer.DialContext(ctx, network, address)
	if nil != err {
		return nil, err
	}
	return &rateLimitConnection{
		conn: conn,
		lr:   stream.Reader(conn),
		lw:   stream.Writer(conn),
		closeWrite: func() error {
			return nil
		},
	}, nil
}

type httpRateLimiters struct {
}

func (that *httpRateLimiters) Name() string {
	return service.DefaultName
}

func (that *httpRateLimiters) New(serverName string, dialer *net.Dialer) px.ContextDialer {
	if "" == serverName || !types.MatchURNDomain(serverName) {
		return &httpRateLimiter{serverName: serverName, nodeId: "", dialer: dialer}
	}
	names := strings.Split(serverName, ".")
	if len(names) < 3 {
		return &httpRateLimiter{serverName: serverName, nodeId: "", dialer: dialer}
	}
	return &httpRateLimiter{serverName: serverName, nodeId: names[len(names)-3], dialer: dialer}
}

type tcpRateLimiter struct {
	hello tx.Hello
	conn  tcp.WriteCloser
}

func (that *tcpRateLimiter) Name() string {
	return tx.DefaultName
}

func (that *tcpRateLimiter) Accept(hello tx.Hello, conn tcp.WriteCloser) tcp.WriteCloser {
	if !proxy.RateLimitEnable {
		return conn
	}
	if nil == hello || !types.MatchURNDomain(hello.ServerName()) {
		return conn
	}
	names := strings.Split(hello.ServerName(), ".")
	if len(names) < 4 {
		return conn
	}
	nodeId := names[len(names)-4]
	stream := limiters.MatchStreamLimiter(strings.ToUpper(nodeId))
	if nil == stream {
		return conn
	}
	log.Info0("TCP connect in %s(%dMbps/s)-(%dMbps/s)[%s-%s]", nodeId, stream.upstream, stream.downstream, conn.LocalAddr().String(), conn.RemoteAddr().String())
	return &rateLimitConnection{
		conn:       conn,
		lr:         stream.Reader(conn),
		lw:         stream.Writer(conn),
		closeWrite: conn.CloseWrite,
	}
}

type rateLimiters struct {
	streams dsa.Map[string, *limitStream]
}

func (that *rateLimiters) Att() *macro.Att {
	return &macro.Att{Name: "mesh.plugin.proxy.limiter"}
}

func (that *rateLimiters) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.NetworkRouteRefresh}
}

func (that *rateLimiters) Listen(ctx context.Context, event *types.Event) error {
	var routes []*types.Route
	if err := event.TryGetObject(&routes); nil != err {
		return cause.Error(err)
	}
	for _, route := range routes {
		bandwidthUp := 1024 * 1024 * route.Upstream / 8
		bandwidthDown := 1024 * 1024 * route.Downstream / 8
		stream := that.MatchStreamLimiter(strings.ToUpper(route.NodeId))
		if nil != stream && !stream.Diff(route) {
			continue
		}
		log.Info(ctx, "Ratelimit %s upstream:%dMbps#%s downstream:%dMbps#%s. ", route.NodeId, route.Upstream, route.URC().String(), route.Downstream, route.StaticIP)
		if nil == stream {
			xstream := &limitStream{
				nodeId:       strings.ToUpper(route.NodeId),
				upstream:     route.Upstream,
				downstream:   route.Downstream,
				readLimiter:  tbucket.NewFreshRateLimiter(float64(bandwidthDown), bandwidthDown, safeBandwidth, safeBandwidth),
				writeLimiter: tbucket.NewFreshRateLimiter(float64(bandwidthUp), bandwidthUp, safeBandwidth, safeBandwidth),
			}
			that.streams.Put(strings.ToUpper(route.NodeId), xstream)
			that.streams.Put(strings.ToUpper(route.InstId), xstream)
			continue
		}
		stream.Refresh(bandwidthUp, bandwidthDown, route)
	}
	return nil
}

func (that *rateLimiters) MatchStreamLimiter(ip string) *limitStream {
	s, ok := that.streams.Get(ip)
	if !ok {
		return nil
	}
	return s
}

type limitStream struct {
	nodeId       string
	upstream     int64
	downstream   int64
	readLimiter  tbucket.FreshRateLimiter
	writeLimiter tbucket.FreshRateLimiter
}

func (that *limitStream) Diff(route *types.Route) bool {
	return that.upstream != route.Upstream || that.downstream != route.Downstream
}

func (that *limitStream) Refresh(bandwidthUp int64, bandwidthDown int64, route *types.Route) {
	that.nodeId = route.NodeId
	that.upstream = route.Upstream
	that.downstream = route.Downstream
	that.readLimiter.Refresh(float64(bandwidthDown), bandwidthDown, safeBandwidth, safeBandwidth)
	that.writeLimiter.Refresh(float64(bandwidthUp), bandwidthUp, safeBandwidth, safeBandwidth)
}

func (that *limitStream) Reader(reader io.Reader) io.Reader {
	return tbucket.Reader(reader, that.readLimiter)
}

func (that *limitStream) Writer(writer io.Writer) io.Writer {
	return tbucket.Writer(writer, that.writeLimiter)
}

type rateLimitConnection struct {
	conn       net.Conn
	lr         io.Reader
	lw         io.Writer
	closeWrite func() error
}

func (that *rateLimitConnection) CloseWrite() error {
	if nil != that.closeWrite {
		return that.closeWrite()
	}
	return nil
}

func (that *rateLimitConnection) Read(b []byte) (int, error) {
	return that.lr.Read(b)
}

func (that *rateLimitConnection) Write(b []byte) (int, error) {
	return that.lw.Write(b)
}

func (that *rateLimitConnection) Close() error {
	return that.conn.Close()
}

func (that *rateLimitConnection) LocalAddr() net.Addr {
	return that.conn.LocalAddr()
}

func (that *rateLimitConnection) RemoteAddr() net.Addr {
	return that.conn.RemoteAddr()
}

func (that *rateLimitConnection) SetDeadline(t time.Time) error {
	return that.conn.SetDeadline(t)
}

func (that *rateLimitConnection) SetReadDeadline(t time.Time) error {
	return that.conn.SetReadDeadline(t)
}

func (that *rateLimitConnection) SetWriteDeadline(t time.Time) error {
	return that.conn.SetWriteDeadline(t)
}
