/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package tool

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Addrs struct {
	servers        []*StatefulServer
	availableAddrs []string
}

type StatefulServer struct {
	available bool
	address   string
}

func NewAddrs(addrs string, port int) *Addrs {
	var servers []*StatefulServer
	var availableAddrs []string
	for _, addr := range strings.Split(addrs, ",") {
		server := strings.TrimSpace(addr)
		if "" != server && !Contains(availableAddrs, server) {
			availableAddr := Ternary(strings.Contains(server, ":"), server, fmt.Sprintf("%s:%d", server, port))
			availableAddrs = append(availableAddrs, availableAddr)
			servers = append(servers, &StatefulServer{
				available: true,
				address:   availableAddr,
			})
		}
	}
	return &Addrs{servers: servers, availableAddrs: availableAddrs}
}

func (that *Addrs) Any() string {
	addrs := that.availableAddrs
	if len(addrs) < 1 {
		return ""
	}
	return addrs[rand.Intn(len(addrs))]
}

func (that *Addrs) Many() []string {
	addrs := that.availableAddrs
	if len(addrs) < 1 {
		return nil
	}
	adds := make([]string, len(addrs))
	copy(adds, addrs)
	return adds
}

func (that *Addrs) All() []string {
	if len(that.servers) < 1 {
		return nil
	}
	var adds []string
	for _, server := range that.servers {
		adds = append(adds, server.address)
	}
	return adds
}

func (that *Addrs) Servers() []*StatefulServer {
	return that.servers
}

func (that *Addrs) Available(addr string, available bool) {
	for _, server := range that.servers {
		if addr == server.address && server.available != available {
			server.available = available
			// Need update
			var availableAddrs []string
			for _, serv := range that.servers {
				if serv.available {
					availableAddrs = append(availableAddrs, serv.address)
				}
			}
			if len(availableAddrs) > 0 {
				that.availableAddrs = availableAddrs
			} else {
				that.availableAddrs = append(availableAddrs, that.servers[rand.Intn(len(that.servers))].address)
			}
			return
		}
	}

}

func NewAddr(addr string, dft int) *Addr {
	if "" == strings.TrimSpace(addr) {
		return new(Addr)
	}
	pair := strings.Split(addr, ":")
	if len(pair) < 2 {
		return &Addr{Host: pair[0], Port: dft}
	}
	if port, err := strconv.Atoi(pair[1]); nil != err {
		return &Addr{Host: pair[0], Port: dft}
	} else {
		return &Addr{Host: pair[0], Port: port}
	}
}

type Addr struct {
	Host string
	Port int
}

func (that Addr) String() string {
	return fmt.Sprintf("%s:%d", that.Host, that.Port)
}

func (that Addr) Empty() bool {
	return "" == that.Host && 0 == that.Port
}
