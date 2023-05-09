/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package plugin

import (
	"context"
	"github.com/be-io/mesh/client/golang/dsa"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

var f = &family{plugins: dsa.NewStringMap[Plugin]()}

type Pip interface {
	Pip(ctx context.Context, r *request.Request) ([]dns.RR, error)
}

type Plugin interface {
	Name() []string

	Priority() int

	ServeDNS(ctx context.Context, pip Pip, r *request.Request) ([]dns.RR, error)
}

func Provide(plugin Plugin) {
	f.Reset(plugin)
}

func Family() Pip {
	return f
}
