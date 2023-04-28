/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/stats"
	"math"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

func init() {
	var _ mpc.Provider = provider
	macro.Provide(mpc.IProvider, provider)
}

const (
	// MaxSize is the maximum possible size for a gRPC message.
	MaxSize = math.MaxInt32
	Name    = "grpc"
)

var provider = new(grpcProvider)

type grpcProvider struct {
	listener net.Listener
	server   *grpc.Server
}

func (that *grpcProvider) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

func (that *grpcProvider) HandleRPC(ctx context.Context, stats stats.RPCStats) {

}

func (that *grpcProvider) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

func (that *grpcProvider) HandleConn(ctx context.Context, stats stats.ConnStats) {

}

func (that *grpcProvider) Att() *macro.Att {
	return &macro.Att{Name: Name, Prototype: true, Constructor: func() macro.SPI {
		return new(grpcProvider)
	}}
}

func (that *grpcProvider) Start(ctx context.Context, address string, tc *tls.Config) error {
	if err := Service.Start(); nil != err {
		return cause.Error(err)
	}
	if "" == address {
		return cause.Errorf("GRPC address cant be empty.")
	}
	var protos []RPCService
	for _, proto := range macro.Load(IRPCService).List() {
		if p, ok := proto.(RPCService); ok {
			protos = append(protos, p)
		}
	}
	options := []grpc.ServerOption{
		grpc.ForceServerCodec(Codec),
		grpc.MaxRecvMsgSize(MaxSize),
		grpc.MaxSendMsgSize(MaxSize),
		grpc.MaxConcurrentStreams(uint32(runtime.NumGoroutine())),
		grpc.StatsHandler(that),
		grpc.UnaryInterceptor(Interceptors.ServerUnary),
		grpc.StreamInterceptor(Interceptors.ServerStream),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:              1 * time.Minute,
			Timeout:           15 * time.Second,
			MaxConnectionIdle: 2 * time.Minute,
		}),
		grpc.UnknownServiceHandler(func(srv interface{}, stream grpc.ServerStream) error {
			defer func() {
				if err := recover(); nil != err {
					log.Error(ctx, "%v", err)
				}
			}()
			if err := Service.OnNext(srv, stream); nil != err {
				log.Error(ctx, err.Error())
				return err
			}
			return nil
		}),
	}
	if nil != tc {
		options = append(options, grpc.Creds(credentials.NewTLS(tc)))
	}
	// serviceDesc is the grpc.ServiceDesc for Greeter service.
	// It's only intended for direct use with grpc.RegisterService,
	// and not to be introspected or modified (even as a copy)

	log.Info(ctx, "Listening and serving HTTP 2.0 on %s", address)
	listener, err := net.Listen("tcp", address)
	if nil != err {
		return cause.Error(err)
	}
	that.listener = listener
	that.server = grpc.NewServer(options...)
	for _, proto := range protos {
		that.server.RegisterService(proto.Metadata())
	}
	//that.server.RegisterService(serviceDesc, Service)
	if err = that.server.Serve(that.listener); nil != err {
		log.Info(ctx, "GRPC server serve : %v", err)
	}
	return nil
}

func (that *grpcProvider) Close() error {
	// Attempt graceful stop (waits for pending RPCs), but force a stop if
	// it doesn't happen in a reasonable amount of time.
	done := make(chan struct{})
	const timeout = 5 * time.Second
	go func() {
		that.server.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		log.Info0("Stopping grpc gracefully is taking longer than %v. Force stopping now. Pending RPCs will be abandoned.", timeout)
		that.server.Stop()
	}
	if nil != that.listener {
		log.Catch(that.listener.Close())
	}
	return nil
}

// TLSConfig define params used to create a tls.Config
type TLSConfig struct {
	CertRequired     bool
	Cert             string
	Key              string
	ServerName       string
	RootCACert       string
	ClientAuth       string
	UseSystemCACerts bool
}

// GenerateServerTLSConfig creates and returns a new *tls.Config with the
// configuration provided.
func (that *grpcProvider) GenerateServerTLSConfig(config *TLSConfig) (tlsCfg *tls.Config, err error) {
	if config.CertRequired {
		tlsCfg = new(tls.Config)
		cert, err := tls.LoadX509KeyPair(config.Cert, config.Key)
		if err != nil {
			return nil, err
		}
		tlsCfg.Certificates = []tls.Certificate{cert}

		pool, err := that.generateCertPool(config.RootCACert, config.UseSystemCACerts)
		if err != nil {
			return nil, err
		}
		tlsCfg.ClientCAs = pool

		auth, err := that.setupClientAuth(config.ClientAuth)
		if err != nil {
			return nil, err
		}
		tlsCfg.ClientAuth = auth

		tlsCfg.MinVersion = tls.VersionTLS12
		tlsCfg.CipherSuites = []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		}

		return tlsCfg, nil
	}
	return nil, nil
}

func (that *grpcProvider) generateCertPool(certPath string, useSystemCA bool) (*x509.CertPool, error) {
	var pool *x509.CertPool
	if useSystemCA {
		var err error
		if pool, err = x509.SystemCertPool(); err != nil {
			return nil, err
		}
	} else {
		pool = x509.NewCertPool()
	}

	if len(certPath) > 0 {
		caFile, err := os.ReadFile(certPath)
		if err != nil {
			return nil, err
		}
		if !pool.AppendCertsFromPEM(caFile) {
			return nil, cause.Errorf("error reading CA file %q", certPath)
		}
	}

	return pool, nil
}

func (that *grpcProvider) setupClientAuth(authType string) (tls.ClientAuthType, error) {
	auth := map[string]tls.ClientAuthType{
		"REQUEST":          tls.RequestClientCert,
		"REQUIREANY":       tls.RequireAnyClientCert,
		"VERIFYIFGIVEN":    tls.VerifyClientCertIfGiven,
		"REQUIREANDVERIFY": tls.RequireAndVerifyClientCert,
	}

	if len(authType) > 0 {
		if v, has := auth[strings.ToUpper(authType)]; has {
			return v, nil
		}
		return tls.NoClientCert, cause.Errorf("Invalid client auth. Valid values [REQUEST, REQUIREANY, VERIFYIFGIVEN, REQUIREANDVERIFY]")
	}

	return tls.NoClientCert, nil
}
