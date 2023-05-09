/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package plugin

import (
	"context"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"runtime/debug"
)

func init() {
	plugin.Register(mdns.Name(), mdns.Setup)
}

func Server() plugin.Handler {
	return mdns
}

var mdns = new(meshDNSServer)

type meshDNSServer struct {
	Next plugin.Handler
}

func (that *meshDNSServer) Name() string {
	return "mesh"
}

func (that *meshDNSServer) Setup(ctrl *caddy.Controller) error {
	dnsserver.GetConfig(ctrl).AddPlugin(func(next plugin.Handler) plugin.Handler {
		that.Next = next
		return that
	})
	ctrl.OnStartup(func() error {
		log.Info0("Setup mesh DNS server.")
		return nil
	})
	return nil
}

func (that *meshDNSServer) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	defer func() {
		if err := recover(); nil != err {
			log.Info(ctx, "DNS resolve with error, %s %s", state.QName(), state.RemoteAddr())
			log.Error(ctx, "%v", err)
			log.Error(ctx, string(debug.Stack()))
		}
	}()
	if state.Type() == "AXFR" {
		return that.errorResponse(state, dns.RcodeNotImplemented, nil)
	}
	answers, err := Family().Pip(ctx, &state)
	if nil != err {
		log.Error(ctx, "%s, %s, %s.", state.QName(), state.RemoteAddr(), err.Error())
		return that.errorResponse(state, dns.RcodeServerFailure, err)
	}
	if len(answers) > 0 {
		msg := new(dns.Msg)
		msg.SetReply(r)
		msg.Authoritative = true
		msg.RecursionAvailable = false
		msg.Compress = true
		msg.Answer = answers

		if err = w.WriteMsg(msg); nil != err {
			log.Error(ctx, err.Error())
		}
		return dns.RcodeSuccess, nil
	}
	if nil != that.Next {
		rt, e := that.Next.ServeDNS(ctx, w, r)
		if nil != e {
			log.Error(ctx, "%s, %s, %s.", state.QName(), state.RemoteAddr(), e.Error())
		}
		return rt, e
	}

	return that.errorResponse(state, dns.RcodeNameError, nil)
}

func (that *meshDNSServer) errorResponse(state request.Request, rCode int, err error) (int, error) {
	msg := new(dns.Msg)
	msg.SetRcode(state.Req, rCode)
	msg.Authoritative, msg.RecursionAvailable, msg.Compress = true, false, true

	state.SizeAndDo(msg)
	_ = state.W.WriteMsg(msg)
	// Return success as the rCode to signal we have written to the client.
	return dns.RcodeSuccess, err
}
