/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"context"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/plugin"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"time"
)

func (that *meshProxy) Register(ctx context.Context) {
	log.Info(ctx, "Proxy active in %v mode, register self as proxy or server. ", plugin.Mode)
	topic := &types.Topic{Topic: prsim.ProxyRegisterEvent.Topic, Code: prsim.ProxyRegisterEvent.Code}
	if _, err := aware.Scheduler.Period(ctx, time.Second*30, topic); nil != err {
		log.Error(ctx, err.Error())
	}
}
