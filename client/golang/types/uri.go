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
	query := u.Query()
	return &URL{
		raw:     uri,
		uri:     u,
		subset:  query.Get("s"),
		proxy:   query.Get("p"),
		skip:    query.Get("k"),
		version: query.Get("v"),
		domain:  query.Get("d"),
		tensor:  query.Get("t"),
		http:    query.Get("h"),
		secure:  query.Get("e"),
	}
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
	}
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

func (that *URL) URL() *url.URL {
	if nil == that.uri {
		return that.uri
	}
	query := that.uri.Query()
	if "" != that.subset {
		query.Set("s", that.subset)
	} else {
		query.Del("s")
	}
	if "" != that.proxy {
		query.Set("p", that.proxy)
	} else {
		query.Del("p")
	}
	if "" != that.skip {
		query.Set("k", that.skip)
	} else {
		query.Del("k")
	}
	if "" != that.version {
		query.Set("v", that.version)
	} else {
		query.Del("v")
	}
	if "" != that.domain {
		query.Set("d", that.domain)
	} else {
		query.Del("d")
	}
	if "" != that.tensor {
		query.Set("t", that.tensor)
	} else {
		query.Del("t")
	}
	if "" != that.http {
		query.Set("h", that.http)
	} else {
		query.Del("h")
	}
	that.uri.RawQuery = query.Encode()
	return that.uri
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

func (that *URL) URC() URC {
	return URC(that.String())
}

func (that *URL) Host() string {
	if nil == that.uri {
		return ""
	}
	return that.uri.Host
}

func (that *URL) SetSchema(schema string) *URL {
	that.uri.Scheme = schema
	that.raw = that.uri.String()
	return that
}
func (that *URL) Username() string {
	if nil == that.uri || nil == that.uri.User {
		return ""
	}
	return that.uri.User.Username()
}

func (that *URL) Password() string {
	if nil == that.uri || nil == that.uri.User {
		return ""
	}
	x, _ := that.uri.User.Password()
	return x
}
