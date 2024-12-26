/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package nsqio

import (
	"context"
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/nsqio/nsq/nsqd"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/plugin"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/tool"
	"github.com/opendatav/mesh/client/golang/types"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func init() {
	var _ prsim.Listener = daemon
	plugin.Provide(daemon)
	macro.Provide(prsim.IPublisher, publisher)
	macro.Provide(prsim.IListener, daemon)
}

const (
	Name         = "nsq"
	AllEventCode = "-"
)

var daemon = new(nsqDaemon)

type nsqDaemon struct {
	core      *nsqd.NSQD
	consumers chan *nsqSubscriber
	clients   []io.Closer
	topics    map[string]*types.Topic
}

func (that *nsqDaemon) Att() *macro.Att {
	return &macro.Att{Name: Name}
}

func (that *nsqDaemon) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.RegistryEventRefresh}
}

func (that *nsqDaemon) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Name: Name, WaitAny: true, Flags: nsqOption{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *nsqDaemon) Start(ctx context.Context, runtime plugin.Runtime) {
	option := &nsqOption{}
	log.Panic(runtime.Parse(option))
	if _, err := os.Stat(filepath.Dir(option.DataPath)); os.IsNotExist(err) {
		if err = tool.MakeDir(option.DataPath); nil != err {
			log.Error(ctx, "Create data home for queue, %s", err.Error())
		}
	}
	rand.Seed(time.Now().UTC().UnixNano())
	opts := nsqd.NewOptions()
	opts.DataPath = option.DataPath
	opts.LogPrefix = option.LogPrefix
	opts.MemQueueSize = option.MemQueueSize
	opts.MsgTimeout = option.MsgTimeout
	opts.MaxMsgTimeout = option.MaxMsgTimeout
	opts.MaxReqTimeout = option.MaxReqTimeout
	opts.SyncEvery = option.SyncEvery
	opts.SyncTimeout = option.SyncTimeout

	if err := log.PError(ctx, func() error {
		if err := tool.MakeDir(opts.DataPath); nil != err {
			return cause.Error(err)
		}
		instance, err := nsqd.New(opts)
		if nil != err {
			return cause.Errorf("Failed to instantiate nsq - %s", err.Error())
		}
		if nil != err {
			return cause.Errorf("Failed to instantiate nsq - %s", err.Error())
		}
		that.core = instance
		that.topics = map[string]*types.Topic{}
		that.consumers = make(chan *nsqSubscriber, 512)
		err = that.core.LoadMetadata()
		if nil != err {
			return cause.Errorf("Failed to load metadata - %s", err.Error())
		}
		err = that.core.PersistMetadata()
		if nil != err {
			return cause.Errorf("Failed to persist metadata - %s", err.Error())
		}
		return nil
	}); nil != err {
		log.Error(ctx, err.Error())
		return
	}
	runtime.Submit(func() {
		if err := log.PError(ctx, func() error {
			return cause.Error(that.core.Main())
		}); nil != err {
			log.Error(ctx, "Start queue broker, %s", err.Error())
			that.Stop(ctx, runtime)
		}
	})
	runtime.Submit(func() {
		that.SubscribeEvent(ctx, runtime)
	})
	runtime.Submit(func() {
		that.BrokerExchange(ctx)
	})
}

func (that *nsqDaemon) Stop(ctx context.Context, runtime plugin.Runtime) {
	for _, client := range that.clients {
		log.Catch(client.Close())
	}
	if nil != that.core {
		that.core.Exit()
	}
}

// Context returns a context that will be canceled when nsqd initiates the shutdown
func (that *nsqDaemon) Context() context.Context {
	return that.core.Context()
}

// Listen will watch the topic in registration
func (that *nsqDaemon) Listen(ctx context.Context, event *types.Event) error {
	if nil == that.consumers {
		log.Warn(ctx, "NSQ instance dont startup. ")
		return nil
	}
	var registrations types.MetadataRegistrations
	if err := event.TryGetObject(&registrations); nil != err {
		return cause.Error(err)
	}
	for _, binding := range registrations.Of(types.METADATA).InferService() {
		if binding.Kind != types.Binding {
			continue
		}
		if "" == binding.Name || "*" == binding.Name {
			binding.Name = AllEventCode
		}
		key := fmt.Sprintf("%s:%s", binding.Namespace, binding.Name)
		if nil != that.topics[key] {
			continue
		}
		topic := &types.Topic{
			Topic: binding.Namespace,
			Code:  binding.Name,
		}
		log.Info(ctx, "Subscribe %s:%s", binding.Namespace, binding.Name)
		consumer := &nsqSubscriber{
			topic:      topic,
			commands:   make(chan *nsq.Command, 31),
			subscriber: aware.Subscriber,
		}
		consumer.Active(ctx)
		that.consumers <- consumer
		that.topics[key] = topic
	}
	return nil
}

// BrokerExchange will watch the topic to remote.
func (that *nsqDaemon) BrokerExchange(ctx context.Context) {
	broker := new(exchangeBroker)
	log.Info(ctx, "Subscribe %s:%s", exchangeTopic.Topic, exchangeTopic.Code)
	consumer := &nsqSubscriber{
		topic:      exchangeTopic,
		commands:   make(chan *nsq.Command, 31),
		subscriber: broker,
	}
	consumer.Active(ctx)
	that.consumers <- consumer
}

func (that *nsqDaemon) SubscribeEvent(ctx context.Context, runtime plugin.Runtime) {
	for {
		select {
		case consumer, ok := <-that.consumers:
			if !ok {
				log.Error(ctx, "Mesh queue consumers have been closed. ")
				return
			}
			log.Devour(func() {
				log.Info(ctx, "Subscriber %s:%s consuming. ", consumer.topic.Topic, consumer.topic.Code)
				runtime.Submit(func() {
					if err := consumer.readLoop(); nil != err {
						log.Error(ctx, "%s:%s subscribe stopped, %s", consumer.topic.Topic, consumer.topic.Code, err.Error())
					}
				})
				runtime.Submit(func() {
					proto := nsqd.NewProtocol(that.core)
					client := proto.NewClient(consumer)
					that.clients = append(that.clients, client)
					if err := proto.IOLoop(client); nil != err {
						log.Error(ctx, "%s:%s subscribe stopped, %s", consumer.topic.Topic, consumer.topic.Code, err.Error())
					}
				})
			})
		}
	}
}
