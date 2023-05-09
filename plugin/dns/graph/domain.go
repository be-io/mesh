/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package graph

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/dsa"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/system"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/miekg/dns"
	"net"
	"strings"
	"time"
)

const (
	AcmeChallenge   = "_acme-challenge"
	NativeDomainKey = "mesh.plugin.dns.domain.native"
)

func (that *graphPlugin) resolve(ctx context.Context, qdn string, remote string) ([]dns.RR, error) {
	if strings.Index(qdn, AcmeChallenge) > -1 {
		return []dns.RR{
			&dns.A{
				A: net.ParseIP(tool.Runtime.Get().Host),
				Hdr: dns.RR_Header{
					Name:   dns.Fqdn(qdn),
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    30,
				}},
		}, nil
	}
	uname := strings.TrimSuffix(qdn, ".")
	if !types.MatchURNDomain(uname) {
		return that.resolveNative(ctx, uname)
	}
	urn := types.FromURN(ctx, uname)
	answers := that.resolveDomain(ctx, urn, uname, prsim.ManuDomain)
	if len(answers) > 0 {
		return answers, nil
	}
	answers = that.resolveDomain(ctx, urn, uname, prsim.AutoDomain)
	if len(answers) > 0 {
		return answers, nil
	}
	answers = that.resolveService(ctx, uname, urn)
	if len(answers) > 0 {
		return answers, nil
	}
	return that.resolveEdges(ctx, uname, urn), nil
}

func (that *graphPlugin) resolveEdges(ctx context.Context, name string, urn *types.URN) []dns.RR {
	var answers []dns.RR
	addresses, err := that.ResolveByRoutes(ctx, urn)
	if nil != err {
		log.Error(ctx, err.Error())
		return nil
	}
	for _, addr := range addresses {
		adds := strings.Split(addr.Addr, ":")
		log.Info(ctx, "DNS resolve %s to %s", name, adds[0])
		answers = append(answers, &dns.A{
			A: net.ParseIP(adds[0]),
			Hdr: dns.RR_Header{
				Name:   dns.Fqdn(name),
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    30,
			}})
	}
	return answers
}

func (that *graphPlugin) resolveService(ctx context.Context, name string, urn *types.URN) []dns.RR {
	var answers []dns.RR
	addresses, err := that.ResolveByServices(ctx, urn, false)
	if nil != err {
		log.Error(ctx, err.Error())
		return nil
	}
	for _, addr := range addresses {
		adds := strings.Split(addr.Addr, ":")
		log.Info(ctx, "DNS resolve %s to %s", name, adds[0])
		answers = append(answers, &dns.A{
			A: net.ParseIP(adds[0]),
			Hdr: dns.RR_Header{
				Name:   dns.Fqdn(name),
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    30,
			}})
	}
	return answers
}

func (that *graphPlugin) resolveNative(ctx context.Context, name string) ([]dns.RR, error) {
	ips := func(retries int) []net.IP {
		var answers []net.IP
		ipKey := fmt.Sprintf("%s.%s", NativeDomainKey, name)
		err := system.GetWithCache(ctx, aware.Cache, ipKey, &answers)
		if nil != err {
			log.Error(ctx, "Cant resolve dns from cache, %s. ", err.Error())
		}
		if len(answers) > 0 {
			return answers
		}
		for index := 0; index < retries; index++ {
			answers, err = dsa.Transform[net.IPAddr, net.IP](func(r net.IPAddr) net.IP {
				return r.IP
			}).Map(net.DefaultResolver.LookupIPAddr(ctx, name))
			if nil == err {
				break
			}
			log.Error(ctx, err.Error())
		}
		if len(answers) > 0 {
			if err = system.PutWithCache(ctx, aware.Cache, ipKey, answers, time.Minute*10); nil != err {
				log.Error(ctx, "Cant cache dns resolve records, %s. ", err.Error())
			}
		}
		return answers
	}(3)
	var answers []dns.RR
	for _, ip := range ips {
		answers = append(answers, &dns.A{
			A: ip,
			Hdr: dns.RR_Header{
				Name:   dns.Fqdn(name),
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    30,
			}})
	}
	return answers, nil
}

func (that *graphPlugin) resolveDomain(ctx context.Context, urn *types.URN, name string, kind string) []dns.RR {
	domains, err := aware.Network.GetDomains(ctx, kind)
	if nil != err {
		log.Error(ctx, err.Error())
		return nil
	}
	var answers []dns.RR
	for _, domain := range domains {
		if "" == domain.Address {
			continue
		}
		if !urn.Match(ctx, types.FromURN(ctx, domain.URN)) {
			continue
		}
		adds := strings.Split(domain.Address, ":")

		log.Info(ctx, "DNS resolve %s to %s", name, adds[0])

		answers = append(answers, &dns.A{
			A: net.ParseIP(adds[0]),
			Hdr: dns.RR_Header{
				Name:   dns.Fqdn(name),
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    30,
			}})
	}
	return answers
}
