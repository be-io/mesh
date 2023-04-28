/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/go-co-op/gocron"
	_ "github.com/robfig/cron/v3"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

func init() {
	var _ prsim.Scheduler = sss
	macro.Provide(prsim.IScheduler, sss)
	scheduler.StartAsync()
}

var (
	scheduler = gocron.NewScheduler(time.Local)
	sss       = &systemScheduler{
		ref:    scheduler,
		topics: &sync.Map{},
	}
)

type systemScheduler struct {
	ref    *gocron.Scheduler
	topics *sync.Map
}

func (that *systemScheduler) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *systemScheduler) Timeout(ctx context.Context, timeout *types.Timeout, duration time.Duration) (string, error) {
	taskId, err := aware.Sequence.Next(ctx, prsim.SCSEQ, 8)
	if nil != err {
		return "", cause.Error(err)
	}
	if _, err = that.ref.Tag(taskId).At(time.UnixMilli(time.Now().UnixMilli() + duration.Milliseconds())).Do(that.doEmit(ctx, timeout.Binding)); nil != err {
		return "", cause.Error(err)
	}
	return taskId, nil
}

func (that *systemScheduler) Cron(ctx context.Context, cron string, topic *types.Topic) (string, error) {
	taskId, err := aware.Sequence.Next(ctx, prsim.SCSEQ, 8)
	if nil != err {
		return "", cause.Error(err)
	}
	if _, err = that.ref.Tag(taskId).CronWithSeconds(cron).Do(that.doEmit(ctx, topic)); nil != err {
		return "", cause.Error(err)
	}
	return taskId, nil
}

func (that *systemScheduler) Period(ctx context.Context, duration time.Duration, topic *types.Topic) (string, error) {
	key := fmt.Sprintf("%s:%s", topic.Topic, topic.Code)
	if _, ok := that.topics.Load(key); ok {
		return key, nil
	}
	that.topics.Store(key, true)
	taskId := tool.NextID()
	if _, err := that.ref.Tag(taskId, key).Every(duration).Do(that.doEmit(ctx, topic)); nil != err {
		return "", cause.Error(err)
	}
	return taskId, nil
}

func (that *systemScheduler) Dump(ctx context.Context) ([]string, error) {
	var taskIds []string
	for _, job := range that.ref.Jobs() {
		if len(job.Tags()) > 0 {
			taskIds = append(taskIds, job.Tags()...)
		}
	}
	return taskIds, nil
}

func (that *systemScheduler) Cancel(ctx context.Context, taskId string) (bool, error) {
	err := that.ref.RemoveByTag(taskId)
	return nil != err, cause.Error(err)
}

func (that *systemScheduler) Stop(ctx context.Context, taskId string) (bool, error) {
	return that.Cancel(ctx, taskId)
}

func (that *systemScheduler) Emit(ctx context.Context, topic *types.Topic) error {
	that.doEmit(ctx, topic)()
	return nil
}

func (that *systemScheduler) Shutdown(ctx context.Context, duration time.Duration) error {
	that.ref.Clear()
	that.ref.Stop()
	return nil
}

func (that *systemScheduler) doEmit(ctx context.Context, topic *types.Topic) func() {
	return func() {
		mtx := mpc.Context()
		defer func() {
			if err := recover(); nil != err {
				log.Error(mtx, "%v", err)
				log.Error(mtx, string(debug.Stack()))
			}
		}()
		log.Debug(mtx, "Scheduler emit %s.%s", topic.Topic, topic.Code)
		env, err := aware.Network.GetEnviron(mtx)
		if nil != err {
			log.Error(mtx, err.Error())
			return
		}
		for _, provider := range macro.Load(prsim.IListener).List() {
			if listener, ok := provider.(prsim.Listener); ok && topic.Match(listener.Btt()...) {
				if err = listener.Listen(mtx, &types.Event{
					Version:   types.MessageVersion,
					Tid:       tool.NewTraceId(),
					Sid:       tool.NewSpanId("", 0),
					Eid:       tool.NewTraceId(),
					Mid:       tool.NewTraceId(),
					Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
					Source: &types.Principal{
						NodeId: env.NodeId,
						InstId: env.InstId,
					},
					Target: &types.Principal{
						NodeId: env.NodeId,
						InstId: env.InstId,
					},
					Binding: topic,
					Entity:  new(types.Entity).AsEmpty(),
				}); nil != err {
					log.Error(mtx, "Pub %s:%s, %s", topic.Topic, topic.Code, err.Error())
				}
			}
		}
	}
}
