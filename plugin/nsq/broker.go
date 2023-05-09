/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package nsqio

import (
	"context"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"strings"
)

func init() {
	var _ prsim.Subscriber = new(exchangeBroker)
}

type exchangeBroker struct {
}

func (that *exchangeBroker) Subscribe(ctx context.Context, event *types.Event) error {
	mtx := mpc.ContextWith(ctx)
	mtx.GetPrincipals().Push(event.Target)
	defer func() {
		mtx.GetPrincipals().Pop()
	}()
	events := make([]*types.Event, 1)
	events[0] = event
	ids, err := aware.RemotePublisher.Publish(mtx, events)
	if len(ids) > 0 {
		log.Info(ctx, "publishedMsgs, %s", strings.Join(ids, ","))
	}
	return err
}
