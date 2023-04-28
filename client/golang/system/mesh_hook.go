/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"context"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"os"
	"strconv"
	"time"
)

func init() {
	var _ prsim.RuntimeHook = new(devKitsHook)
	macro.Provide(prsim.IRuntimeHook, new(devKitsHook))
}

type devKitsHook struct {
}

func (that *devKitsHook) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *devKitsHook) Start(ctx context.Context, runtime prsim.Runtime) error {
	for _, spi := range macro.Load(mpc.IConsumer).List() {
		if consumer, ok := spi.(mpc.Consumer); ok {
			if err := consumer.Start(); nil != err {
				log.Error(ctx, err.Error())
			}
		}
	}
	log.AddProcessShutdownHook(func() error {
		for _, spi := range macro.Load(prsim.IRuntimeHook).List() {
			if hook, ok := spi.(prsim.RuntimeHook); ok {
				if err := hook.Stop(ctx, runtime); nil != err {
					log.Error(ctx, "Shutdown hook exec with error, %s", err.Error())
				}
			}
		}
		return nil
	})
	clock := time.Second * 2
	if x := os.Getenv("MESH_SYSTEM_CLOCK"); "" != x {
		xi, _ := strconv.Atoi(x)
		if xi > 0 {
			clock = time.Second * time.Duration(xi)
		}
	}
	topic := &types.Topic{Topic: prsim.SystemClock.Topic, Code: prsim.SystemClock.Code}
	if _, err := sss.Period(macro.Context(), clock, topic); nil != err {
		log.Error(macro.Context(), err.Error())
	}
	return nil
}

func (that *devKitsHook) Stop(ctx context.Context, runtime prsim.Runtime) error {
	return nil
}

func (that *devKitsHook) Refresh(ctx context.Context, runtime prsim.Runtime) error {
	return nil
}
