/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package redis

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/miekg/dns"
	"strings"
	"time"
)

func (that *Redis) LoadZones(ctx context.Context) error {
	zones, err := aware.Cache.Keys(ctx, that.keyPrefix+"*"+that.keySuffix)
	if nil != err {
		return cause.Error(err)
	}
	for i, _ := range zones {
		zones[i] = strings.TrimPrefix(zones[i], that.keyPrefix)
		zones[i] = strings.TrimSuffix(zones[i], that.keySuffix)
	}
	that.LastZoneUpdate = time.Now()
	that.Zones = zones
	return nil
}

func (that *Redis) A(name string, z *Zone, record *Record) (answers, extras []dns.RR) {
	for _, a := range record.A {
		if a.Ip == nil {
			continue
		}
		r := new(dns.A)
		r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeA,
			Class: dns.ClassINET, Ttl: that.minTtl(a.Ttl)}
		r.A = a.Ip
		answers = append(answers, r)
	}
	return
}

func (that Redis) AAAA(name string, z *Zone, record *Record) (answers, extras []dns.RR) {
	for _, aaaa := range record.AAAA {
		if aaaa.Ip == nil {
			continue
		}
		r := new(dns.AAAA)
		r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeAAAA,
			Class: dns.ClassINET, Ttl: that.minTtl(aaaa.Ttl)}
		r.AAAA = aaaa.Ip
		answers = append(answers, r)
	}
	return
}

func (that *Redis) CNAME(name string, z *Zone, record *Record) (answers, extras []dns.RR) {
	for _, cname := range record.CNAME {
		if len(cname.Host) == 0 {
			continue
		}
		r := new(dns.CNAME)
		r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeCNAME,
			Class: dns.ClassINET, Ttl: that.minTtl(cname.Ttl)}
		r.Target = dns.Fqdn(cname.Host)
		answers = append(answers, r)
	}
	return
}

func (that *Redis) TXT(name string, z *Zone, record *Record) (answers, extras []dns.RR) {
	for _, txt := range record.TXT {
		if len(txt.Text) == 0 {
			continue
		}
		r := new(dns.TXT)
		r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeTXT,
			Class: dns.ClassINET, Ttl: that.minTtl(txt.Ttl)}
		r.Txt = split255(txt.Text)
		answers = append(answers, r)
	}
	return
}

func (that *Redis) NS(ctx context.Context, name string, z *Zone, record *Record) (answers, extras []dns.RR) {
	for _, ns := range record.NS {
		if len(ns.Host) == 0 {
			continue
		}
		r := new(dns.NS)
		r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeNS,
			Class: dns.ClassINET, Ttl: that.minTtl(ns.Ttl)}
		r.Ns = ns.Host
		answers = append(answers, r)
		extras = append(extras, that.hosts(ctx, ns.Host, z)...)
	}
	return
}

func (that *Redis) MX(ctx context.Context, name string, z *Zone, record *Record) (answers, extras []dns.RR) {
	for _, mx := range record.MX {
		if len(mx.Host) == 0 {
			continue
		}
		r := new(dns.MX)
		r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeMX,
			Class: dns.ClassINET, Ttl: that.minTtl(mx.Ttl)}
		r.Mx = mx.Host
		r.Preference = mx.Preference
		answers = append(answers, r)
		extras = append(extras, that.hosts(ctx, mx.Host, z)...)
	}
	return
}

func (that *Redis) SRV(ctx context.Context, name string, z *Zone, record *Record) (answers, extras []dns.RR) {
	for _, srv := range record.SRV {
		if len(srv.Target) == 0 {
			continue
		}
		r := new(dns.SRV)
		r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeSRV,
			Class: dns.ClassINET, Ttl: that.minTtl(srv.Ttl)}
		r.Target = srv.Target
		r.Weight = srv.Weight
		r.Port = srv.Port
		r.Priority = srv.Priority
		answers = append(answers, r)
		extras = append(extras, that.hosts(ctx, srv.Target, z)...)
	}
	return
}

func (that *Redis) SOA(name string, z *Zone, record *Record) (answers, extras []dns.RR) {
	r := new(dns.SOA)
	if record.SOA.Ns == "" {
		r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeSOA,
			Class: dns.ClassINET, Ttl: that.Ttl}
		r.Ns = "ns1." + name
		r.Mbox = "hostmaster." + name
		r.Refresh = 86400
		r.Retry = 7200
		r.Expire = 3600
		r.Minttl = that.Ttl
	} else {
		r.Hdr = dns.RR_Header{Name: dns.Fqdn(z.Name), Rrtype: dns.TypeSOA,
			Class: dns.ClassINET, Ttl: that.minTtl(record.SOA.Ttl)}
		r.Ns = record.SOA.Ns
		r.Mbox = record.SOA.MBox
		r.Refresh = record.SOA.Refresh
		r.Retry = record.SOA.Retry
		r.Expire = record.SOA.Expire
		r.Minttl = record.SOA.MinTtl
	}
	r.Serial = that.serial()
	answers = append(answers, r)
	return
}

func (that *Redis) CAA(name string, z *Zone, record *Record) (answers, extras []dns.RR) {
	if record == nil {
		return
	}
	for _, caa := range record.CAA {
		if caa.Value == "" || caa.Tag == "" {
			continue
		}
		r := new(dns.CAA)
		r.Hdr = dns.RR_Header{Name: dns.Fqdn(name), Rrtype: dns.TypeCAA, Class: dns.ClassINET}
		r.Flag = caa.Flag
		r.Tag = caa.Tag
		r.Value = caa.Value
		answers = append(answers, r)
	}
	return
}

func (that *Redis) AXFR(ctx context.Context, z *Zone) (records []dns.RR) {
	//soa, _ := redis.SOA(z.Name, z, record)
	soa := make([]dns.RR, 0)
	answers := make([]dns.RR, 0, 10)
	extras := make([]dns.RR, 0, 10)

	// Allocate slices for rr Records
	records = append(records, soa...)
	for key := range z.Locations {
		if key == "@" {
			location := that.findLocation(z.Name, z)
			record := that.get(ctx, location, z)
			soa, _ = that.SOA(z.Name, z, record)
		} else {
			fqdnKey := dns.Fqdn(key) + z.Name
			var as []dns.RR
			var xs []dns.RR

			location := that.findLocation(fqdnKey, z)
			record := that.get(ctx, location, z)

			// Pull all zone records
			as, xs = that.A(fqdnKey, z, record)
			answers = append(answers, as...)
			extras = append(extras, xs...)

			as, xs = that.AAAA(fqdnKey, z, record)
			answers = append(answers, as...)
			extras = append(extras, xs...)

			as, xs = that.CNAME(fqdnKey, z, record)
			answers = append(answers, as...)
			extras = append(extras, xs...)

			as, xs = that.MX(ctx, fqdnKey, z, record)
			answers = append(answers, as...)
			extras = append(extras, xs...)

			as, xs = that.SRV(ctx, fqdnKey, z, record)
			answers = append(answers, as...)
			extras = append(extras, xs...)

			as, xs = that.TXT(fqdnKey, z, record)
			answers = append(answers, as...)
			extras = append(extras, xs...)
		}
	}

	records = soa
	records = append(records, answers...)
	records = append(records, extras...)
	records = append(records, soa...)

	fmt.Println(records)
	return
}

func (that *Redis) hosts(ctx context.Context, name string, z *Zone) []dns.RR {
	var (
		record  *Record
		answers []dns.RR
	)
	location := that.findLocation(name, z)
	if location == "" {
		return nil
	}
	record = that.get(ctx, location, z)
	a, _ := that.A(name, z, record)
	answers = append(answers, a...)
	aaaa, _ := that.AAAA(name, z, record)
	answers = append(answers, aaaa...)
	cname, _ := that.CNAME(name, z, record)
	answers = append(answers, cname...)
	return answers
}

func (that *Redis) serial() uint32 {
	return uint32(time.Now().Unix())
}

func (that *Redis) minTtl(ttl uint32) uint32 {
	if that.Ttl == 0 && ttl == 0 {
		return defaultTtl
	}
	if that.Ttl == 0 {
		return ttl
	}
	if ttl == 0 {
		return that.Ttl
	}
	if that.Ttl < ttl {
		return that.Ttl
	}
	return ttl
}

func (that *Redis) findLocation(query string, z *Zone) string {
	var (
		ok                                 bool
		closestEncloser, sourceOfSynthesis string
	)

	// request for zone records
	if query == z.Name {
		return query
	}

	query = strings.TrimSuffix(query, "."+z.Name)

	if _, ok = z.Locations[query]; ok {
		return query
	}

	closestEncloser, sourceOfSynthesis, ok = splitQuery(query)
	for ok {
		ceExists := keyMatches(closestEncloser, z) || keyExists(closestEncloser, z)
		ssExists := keyExists(sourceOfSynthesis, z)
		if ceExists {
			if ssExists {
				return sourceOfSynthesis
			} else {
				return ""
			}
		} else {
			closestEncloser, sourceOfSynthesis, ok = splitQuery(closestEncloser)
		}
	}
	return ""
}

func (that *Redis) get(ctx context.Context, key string, z *Zone) *Record {
	var (
		err error
		val string
	)
	var label string
	if key == z.Name {
		label = "@"
	} else {
		label = key
	}

	reply, err := aware.Cache.HGet(ctx, that.keyPrefix+z.Name+that.keySuffix, label)
	if nil != err {
		return nil
	}
	if err = reply.Entity.TryReadObject(&val); nil != err {
		return nil
	}
	r := new(Record)
	err = codec.Jsonizer.Unmarshal([]byte(val), r)
	if nil != err {
		fmt.Println("parse error : ", val, err)
		return nil
	}
	return r
}

func keyExists(key string, z *Zone) bool {
	_, ok := z.Locations[key]
	return ok
}

func keyMatches(key string, z *Zone) bool {
	for value := range z.Locations {
		if strings.HasSuffix(value, key) {
			return true
		}
	}
	return false
}

func splitQuery(query string) (string, string, bool) {
	if query == "" {
		return "", "", false
	}
	var (
		splits            []string
		closestEncloser   string
		sourceOfSynthesis string
	)
	splits = strings.SplitAfterN(query, ".", 2)
	if len(splits) == 2 {
		closestEncloser = splits[1]
		sourceOfSynthesis = "*." + closestEncloser
	} else {
		closestEncloser = ""
		sourceOfSynthesis = "*"
	}
	return closestEncloser, sourceOfSynthesis, true
}

func (that *Redis) save(ctx context.Context, zone string, subdomain string, value string) error {
	entity, err := new(types.CacheEntity).Wrap(subdomain, value, time.Minute*30)
	if nil != err {
		return cause.Error(err)
	}
	return aware.Cache.HSet(ctx, that.keyPrefix+zone+that.keySuffix, entity)
}

func (that *Redis) load(ctx context.Context, zone string) *Zone {
	var (
		err  error
		vals []string
	)
	z := new(Zone)
	z.Name = zone
	vals, err = aware.Cache.HKeys(ctx, that.keyPrefix+zone+that.keySuffix)
	if nil != err {
		return nil
	}
	z.Locations = make(map[string]struct{})
	for _, val := range vals {
		z.Locations[val] = struct{}{}
	}

	return z
}

func split255(s string) []string {
	if len(s) < 255 {
		return []string{s}
	}
	sx := []string{}
	p, i := 0, 255
	for {
		if i <= len(s) {
			sx = append(sx, s[p:i])
		} else {
			sx = append(sx, s[p:])
			break

		}
		p, i = p+255, i+255
	}

	return sx
}

const (
	defaultTtl     = 360
	hostmaster     = "hostmaster"
	zoneUpdateTime = 10 * time.Minute
	transferLength = 1000
)
