/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package redis

import (
	"github.com/be-io/mesh/plugin/dns/plugin"
	cdp "github.com/coredns/coredns/plugin"
	"time"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

func init() {
	var _ plugin.Plugin = new(Redis)
}

type Redis struct {
	keyPrefix      string
	keySuffix      string
	Ttl            uint32
	Zones          []string
	LastZoneUpdate time.Time
}

func (that *Redis) Name() []string {
	return []string{"redis"}
}

func (that *Redis) Priority() int {
	return 0
}

func (that *Redis) ServeDNS(ctx context.Context, pip plugin.Pip, r *request.Request) ([]dns.RR, error) {
	qname := r.Name()
	qtype := r.Type()

	if time.Since(that.LastZoneUpdate) > zoneUpdateTime {
		that.LoadZones(ctx)
	}

	zone := cdp.Zones(that.Zones).Matches(qname)
	// fmt.Println("zone : ", zone)
	if zone == "" {
		return nil, nil
	}

	z := that.load(ctx, zone)
	if z == nil {
		return nil, nil
	}

	if qtype == "AXFR" {
		records := that.AXFR(ctx, z)
		return records, nil
	}

	location := that.findLocation(qname, z)
	if len(location) == 0 { // empty, no results
		return nil, nil
	}

	answers := make([]dns.RR, 0, 10)

	record := that.get(ctx, location, z)

	switch qtype {
	case "A":
		answers, _ = that.A(qname, z, record)
	case "AAAA":
		answers, _ = that.AAAA(qname, z, record)
	case "CNAME":
		answers, _ = that.CNAME(qname, z, record)
	case "TXT":
		answers, _ = that.TXT(qname, z, record)
	case "NS":
		answers, _ = that.NS(ctx, qname, z, record)
	case "MX":
		answers, _ = that.MX(ctx, qname, z, record)
	case "SRV":
		answers, _ = that.SRV(ctx, qname, z, record)
	case "SOA":
		answers, _ = that.SOA(qname, z, record)
	case "CAA":
		answers, _ = that.CAA(qname, z, record)
	default:
		return nil, nil
	}
	return answers, nil
}
