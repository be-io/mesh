/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/be-io/mesh/plugin/redis"
	"github.com/be-io/mesh/plugin/redis/iset"
	"time"
)

var _ prsim.Session = new(PRSISession)

// PRSISession
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PRSISession struct {
}

func (that *PRSISession) Peek(ctx context.Context, topic string) ([]byte, error) {
	mtx := mpc.ContextWith(ctx)
	chanId := fmt.Sprintf("%s.%s", prsim.MeshSessionId.Get(mtx.GetAttachments()), topic)
	log.Info(ctx, "Session %s pop", chanId)
	rc, err := redis.Ref(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	rt, err := rc.LPop(ctx, chanId).Result()
	if len(rt) < 2 || iset.IsNil(err) {
		return nil, nil
	}
	if nil != err {
		return nil, cause.Error(err)
	}
	return []byte(rt), nil
}

func (that *PRSISession) Pop(ctx context.Context, timeout types.Duration, topic string) ([]byte, error) {
	mtx := mpc.ContextWith(ctx)
	chanId := fmt.Sprintf("%s.%s", prsim.MeshSessionId.Get(mtx.GetAttachments()), topic)
	log.Info(ctx, "Session %s pop", chanId)
	rc, err := redis.Ref(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	rt, err := rc.BLPop(ctx, time.Duration(timeout), chanId).Result()
	if len(rt) < 2 || iset.IsNil(err) {
		return nil, nil
	}
	if nil != err {
		return nil, cause.Error(err)
	}
	return []byte(rt[1]), nil
}

func (that *PRSISession) Push(ctx context.Context, payload []byte, metadata map[string]string, topic string) error {
	mtx := mpc.ContextWith(ctx)
	chanId := fmt.Sprintf("%s.%s", prsim.MeshSessionId.Get(mtx.GetAttachments()), topic)
	log.Info(ctx, "Session %s push", chanId)
	rc, err := redis.Ref(ctx)
	if nil != err {
		return cause.Error(err)
	}
	_, err = rc.RPush(ctx, chanId, payload).Result()
	return cause.Error(err)
}

func (that *PRSISession) Release(ctx context.Context, timeout types.Duration, topic string) error {
	mtx := mpc.ContextWith(ctx)
	rc, err := redis.Ref(ctx)
	if nil != err {
		return cause.Error(err)
	}
	if "" != topic {
		chanId := fmt.Sprintf("%s.%s", prsim.MeshSessionId.Get(mtx.GetAttachments()), topic)
		log.Info(ctx, "Session %s release", chanId)
		_, err = rc.Del(ctx, chanId).Result()
		return cause.Error(err)

	}
	sessionId := prsim.MeshSessionId.Get(mtx.GetAttachments())
	keys, err := rc.Keys(ctx, fmt.Sprintf("%s*", sessionId)).Result()
	if nil != err {
		return cause.Error(err)
	}
	for _, key := range keys {
		log.Info(ctx, "Session %s release", key)
		_, err = rc.Del(ctx, key).Result()
		return cause.Error(err)
	}
	return nil
}
