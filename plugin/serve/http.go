/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package serve

import (
	"context"
	"github.com/opendatav/mesh/client/golang/boost"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/plugin"
)

func init() {
	plugin.Provide(new(httpServe))
}

type httpServe struct {
	Serve *boost.Mooter `json:"-"`
}

func (that *httpServe) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Priority: -100, Name: plugin.Http, Flags: httpServe{}, Create: func() plugin.Plugin {
		return new(httpServe)
	}}
}

func (that *httpServe) Start(ctx context.Context, runtime plugin.Runtime) {
	log.Panic(runtime.Parse(that))
	that.Serve = new(boost.Mooter)
	if err := that.Serve.Start(ctx, runtime); nil != err {
		log.Error(ctx, "Http broker up with unexpect error, %s", err.Error())
	}
}

func (that *httpServe) Stop(ctx context.Context, runtime plugin.Runtime) {
	if nil != that.Serve {
		log.Catch(that.Serve.Stop(ctx, runtime))
	}
}
