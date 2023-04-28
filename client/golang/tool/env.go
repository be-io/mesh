/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package tool

import (
	"bytes"
	"fmt"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	Port    = 7304
	MPCPort = 8864
	pandora = "pandora"
)

var (
	IP = new(macro.Once[string]).With(func() string {
		conn, err := net.DialTimeout("udp", "8.8.8.8:80", time.Second*1)
		if nil != err {
			log.Error0(err.Error())
			return "127.0.0.1"
		}
		defer log.Catch(conn.Close())
		return conn.LocalAddr().(*net.UDPAddr).IP.String()
	})
	IPHex = new(macro.Once[string]).With(func() string {
		traceId := bytes.Buffer{}
		for _, feg := range strings.Split(IP.Get(), ".") {
			digit, _ := strconv.Atoi(feg)
			traceId.WriteString(Padding(fmt.Sprintf("%X", digit), 2, "0"))
		}
		return traceId.String()
	})
	Host = new(macro.Once[string]).With(func() string {
		name, err := os.Hostname()
		if nil != err {
			log.Error0(err.Error())
		}
		return name
	})
	AvailablePort = new(macro.Once[string]).With(func() string {
		listener, err := net.Listen("tcp", ":0")
		if nil != err {
			log.Error0(err.Error())
			return "8864"
		}
		defer log.Catch(listener.Close())
		return strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	})
	Runtime = new(macro.Once[*Addr]).With(func() *Addr {
		return NewAddr(Anyone(macro.Runtime(), fmt.Sprintf("%s:%s", IP.Get(), AvailablePort.Get())), MPCPort)
	})
	Address = new(macro.Once[*Addrs]).With(func() *Addrs { return NewAddrs(macro.Address(), Port) })
	Name    = new(macro.Once[string]).With(func() string { return macro.Name() })
	Direct  = new(macro.Once[string]).With(func() string { return macro.Direct() })
	Subset  = new(macro.Once[string]).With(func() string { return macro.Subset() })
	Proxy   = new(macro.Once[*Addr]).With(func() *Addr { return NewAddr(macro.Proxy(), Port) })
	MDC     = new(macro.Once[[]*Addr]).With(func() []*Addr {
		var adds []*Addr
		dcs := strings.Split(macro.MDC(), ",")
		for _, dc := range dcs {
			if "" == strings.TrimSpace(dc) {
				continue
			}

			adds = append(adds, NewAddr(dc, Port))
		}
		return adds
	})
	SPA = new(macro.Once[*Addr]).With(func() *Addr { return NewAddr(Anyone(macro.SPA(), macro.Runtime()), Port) })
)
