/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/log"
	"net/url"
	"regexp"
	"strings"
)

type URC string

func (that URC) String() string {
	return string(that)
}

func (that URC) URL(ctx context.Context) *URL {
	return ParseURL(ctx, that.String())
}

func (that URC) Addrs() []string {
	uri, err := FormatURL(that.String())
	if nil != err {
		log.Error0("Parse %s, %s", that.String(), err.Error())
		return []string{string(that)}
	}
	return strings.Split(uri.Host, ",")
}

func (that URC) Failover() ([]string, []string) {
	members := that.Members()
	if len(members) > 1 {
		return members[0 : len(members)-1], members[len(members)-1:]
	}
	return members, nil
}

func (that URC) Members() []string {
	uri, err := FormatURL(that.String())
	if nil != err {
		log.Error0("Parse %s, %s", that.String(), err.Error())
		return []string{string(that)}
	}
	var members []string
	for _, addr := range that.Addrs() {
		u := url.URL{
			Scheme:      uri.Scheme,
			Opaque:      uri.Opaque,
			User:        uri.User,
			Host:        addr,
			Path:        uri.Path,
			RawPath:     uri.RawPath,
			OmitHost:    uri.OmitHost,
			ForceQuery:  uri.ForceQuery,
			RawQuery:    uri.RawQuery,
			Fragment:    uri.Fragment,
			RawFragment: uri.RawFragment,
		}
		if strings.Index(that.String(), uri.Scheme) != 0 {
			members = append(members, strings.ReplaceAll(u.String(), fmt.Sprintf("%s://", uri.Scheme), ""))
		} else {
			members = append(members, u.String())
		}
	}
	return members
}

func (that URC) IPv() []string {
	var ips []string
	for _, addr := range that.Addrs() {
		ips = append(ips, strings.Split(addr, ":")[0])
	}
	return ips
}

type URL struct {
	raw     string
	uri     *url.URL // nil able
	subset  string
	proxy   string
	skip    string
	version string
	domain  string
	tensor  string
	http    string
	secure  string
	vip     string
	fido    string
	oneway  string
	length  string
}

func FormatURL(addr string) (*url.URL, error) {
	if !strings.Contains(addr, "://") {
		return url.Parse(fmt.Sprintf("https://%s", addr))
	}
	pair := strings.Split(addr, "://")
	if ok, err := regexp.MatchString("^\\w+$", pair[0]); nil == err && ok {
		return url.Parse(addr)
	}
	return url.Parse(fmt.Sprintf("https://%s", addr))
}

func ParseURL(ctx context.Context, uri string) *URL {
	u, err := FormatURL(uri)
	if nil != err {
		log.Error(ctx, "Parse url %s, %s", uri, err.Error())
		return &URL{raw: uri}
	}
	return FromURL(u)
}

func FromURL(uri *url.URL) *URL {
	if nil == uri {
		return &URL{}
	}
	query := uri.Query()
	return &URL{
		raw:     uri.String(),
		uri:     uri,
		subset:  query.Get("s"),
		proxy:   query.Get("p"),
		skip:    query.Get("k"),
		version: query.Get("v"),
		domain:  query.Get("d"),
		tensor:  query.Get("t"),
		http:    query.Get("h"),
		secure:  query.Get("e"),
		vip:     query.Get("w"),
		fido:    query.Get("f"),
		oneway:  query.Get("o"),
		length:  query.Get("l"),
	}
}

func (that *URL) URL() *url.URL {
	if nil == that.uri {
		return that.uri
	}
	query := that.uri.Query()
	that.Force(query, "s", that.subset)
	that.Force(query, "p", that.proxy)
	that.Force(query, "k", that.skip)
	that.Force(query, "v", that.version)
	that.Force(query, "d", that.domain)
	that.Force(query, "t", that.tensor)
	that.Force(query, "h", that.http)
	that.Force(query, "e", that.secure)
	that.Force(query, "w", that.vip)
	that.Force(query, "f", that.fido)
	that.Force(query, "o", that.oneway)
	that.Force(query, "l", that.length)
	that.uri.RawQuery = query.Encode()
	return that.uri
}

func (that *URL) SetS(subset string) *URL {
	that.subset = subset
	return that
}

func (that *URL) GetS() string {
	return that.subset
}

func (that *URL) SetP(proxy string) *URL {
	that.proxy = proxy
	return that
}

func (that *URL) GetP() string {
	return that.proxy
}

func (that *URL) SetK(skip string) *URL {
	that.skip = skip
	return that
}

func (that *URL) GetK() string {
	return that.skip
}

func (that *URL) SetV(version string) *URL {
	that.version = version
	return that
}

func (that *URL) GetV() string {
	return that.version
}

func (that *URL) SetD(domain string) *URL {
	that.domain = domain
	return that
}

func (that *URL) GetD() string {
	return that.domain
}

func (that *URL) SetT(tensor string) *URL {
	that.tensor = tensor
	return that
}

func (that *URL) GetT() string {
	return that.tensor
}

func (that *URL) SetH(http string) *URL {
	that.http = http
	return that
}

func (that *URL) GetH() string {
	return that.http
}

func (that *URL) SetE(secure string) *URL {
	that.secure = secure
	return that
}

func (that *URL) GetE() string {
	return that.secure
}

func (that *URL) SetW(vip string) *URL {
	that.vip = vip
	return that
}

func (that *URL) GetW() string {
	return that.vip
}

func (that *URL) SetF(fido string) *URL {
	that.fido = fido
	return that
}

func (that *URL) GetF() string {
	return that.fido
}

func (that *URL) SetO(oneway string) *URL {
	that.oneway = oneway
	return that
}

func (that *URL) GetO() string {
	return that.oneway
}

func (that *URL) SetL(length string) *URL {
	that.length = length
	return that
}

func (that *URL) GetL() string {
	return that.length
}

func (that *URL) Force(query url.Values, name string, value string) {
	if "" != value {
		query.Set(name, value)
	} else {
		query.Del(name)
	}
}

func (that *URL) NQ() string {
	uri := that.URL()
	if nil == uri {
		return that.raw
	}
	uri.RawQuery = ""
	return uri.String()
}

func (that *URL) String() string {
	uri := that.URL()
	if nil == uri {
		return that.raw
	}
	return uri.String()
}

func (that *URL) SetHost(host string) *URL {
	that.uri.Host = host
	that.raw = that.uri.String()
	return that
}

func (that *URL) SetSchema(schema string) *URL {
	that.uri.Scheme = schema
	that.raw = that.uri.String()
	return that
}

func (that *URL) URC() URC {
	return URC(that.String())
}
