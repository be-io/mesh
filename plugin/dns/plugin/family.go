/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package plugin

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/dsa"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"math"
	"sort"
)

var _ Pip = new(family)

type family struct {
	plugins dsa.Map[string, Plugin]
	pip     Pip
}

func (that *family) Pip(ctx context.Context, r *request.Request) ([]dns.RR, error) {
	return that.pip.Pip(ctx, r)
}

func (that *family) ServeDNS(ctx context.Context, pip Pip, r *request.Request) ([]dns.RR, error) {
	rr, err := that.pip.Pip(ctx, r)
	if nil != err {
		return nil, cause.Error(err)
	}
	if len(rr) > 0 {
		return rr, nil
	}
	return pip.Pip(ctx, r)
}

func (that *family) Name() []string {
	var names []string
	that.plugins.ForEach(func(key string, v Plugin) {
		names = append(names, v.Name()...)
	})
	return names
}

func (that *family) Priority() int {
	return math.MaxInt
}

func (that *family) Reset(plugin Plugin) {
	if len(plugin.Name()) > 0 {
		that.plugins.Put(plugin.Name()[0], plugin)
	}
	plugins := that.plugins.Values()
	sort.SliceStable(plugins, func(i, j int) bool {
		return plugins[i].Priority() < plugins[j].Priority()
	})
	var pip Pip = &puppet{}
	for _, p := range plugins {
		pip = &puppet{plugin: p, pip: pip}
	}
	that.pip = pip
}

var _ Pip = new(puppet)

type puppet struct {
	plugin Plugin
	pip    Pip
}

func (that *puppet) Name() []string {
	return that.plugin.Name()
}

func (that *puppet) Priority() int {
	return that.plugin.Priority()
}

func (that *puppet) Pip(ctx context.Context, r *request.Request) ([]dns.RR, error) {
	if nil == that.plugin {
		return nil, nil
	}
	return that.plugin.ServeDNS(ctx, that.pip, r)
}
