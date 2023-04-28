/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package nsqio

import (
	"bytes"
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/nsqio/go-nsq"
	"io"
	"net"
	"runtime/debug"
	"time"
)

func init() {
	var _ net.Conn = new(nsqSubscriber)
	var _ prsim.Subscriber = new(nsqSubscriber)
}

type nsqSubscriber struct {
	commands      chan *nsq.Command
	topic         *types.Topic
	deadline      time.Time
	readDeadline  time.Time
	writeDeadline time.Time
	subscriber    prsim.Subscriber
	eventReader   io.ReadCloser
	eventWriter   io.WriteCloser
}

func (that *nsqSubscriber) Active(ctx context.Context) {
	that.eventReader, that.eventWriter = io.Pipe()
	that.commands <- nsq.Subscribe(that.topic.Topic, that.topic.Code)
	that.commands <- nsq.Ready(31)
	log.Info(ctx, "Subscriber %s:%s active. ", that.topic.Topic, that.topic.Code)
}

func (that *nsqSubscriber) Read(b []byte) (int, error) {
	select {
	case command := <-that.commands:
		var buff bytes.Buffer
		if _, err := command.WriteTo(&buff); nil != err {
			log.Error0("Send command to nsq - %s", err.Error())
			return 0, cause.Error(err)
		}
		copy(b, buff.Bytes())
		return buff.Len(), nil
	}
}

func (that *nsqSubscriber) Write(b []byte) (n int, err error) {
	return that.eventWriter.Write(b)
}

func (that *nsqSubscriber) readLoop() error {
	for {
		frameType, data, err := nsq.ReadUnpackedResponse(that.eventReader)
		if nil != err {
			log.Error0("Unpack message %s - %s, but cant discard it. ", string(data), err.Error())
			return cause.Error(err)
		}
		switch frameType {
		case nsq.FrameTypeResponse:
			if !bytes.Equal(data, []byte("_heartbeat_")) {
				continue
			}
			that.commands <- nsq.Nop()
		case nsq.FrameTypeMessage:
			message, err := nsq.DecodeMessage(data)
			if nil != err {
				log.Error0("Decode message %s - %s, but cant discard it. ", string(data), err.Error())
				continue
			}
			if err = that.doSubscribe(mpc.Context(), message); nil != err {
				log.Error0("ACK event %s - %s, requeue it. ", string(message.Body), err.Error())
				that.commands <- nsq.Requeue(message.ID, time.Second*5)
				continue
			}
			log.Debug0("ACK finish of event %s -. ", string(message.Body))
			that.commands <- nsq.Finish(message.ID)
		case nsq.FrameTypeError:
			log.Error0("Protocol error - %s", string(data))
		default:
			log.Error0("Unknown frame type %d - %s", frameType, string(data))
		}
	}
}

func (that *nsqSubscriber) Close() error {
	if nil != that.eventReader {
		log.Catch(that.eventReader.Close())
	}
	if nil != that.eventWriter {
		log.Catch(that.eventWriter.Close())
	}
	return nil
}

func (that *nsqSubscriber) LocalAddr() net.Addr {
	return &net.IPAddr{IP: net.ParseIP("127.0.0.1")}
}

func (that *nsqSubscriber) RemoteAddr() net.Addr {
	return &net.IPAddr{IP: net.ParseIP("127.0.0.1")}
}

func (that *nsqSubscriber) SetDeadline(deadline time.Time) error {
	that.deadline = deadline
	return nil
}

func (that *nsqSubscriber) SetReadDeadline(deadline time.Time) error {
	that.readDeadline = deadline
	return nil
}

func (that *nsqSubscriber) SetWriteDeadline(deadline time.Time) error {
	that.writeDeadline = deadline
	return nil
}

func (that *nsqSubscriber) doSubscribe(ctx context.Context, message *nsq.Message) (err error) {
	defer func() {
		if e := recover(); nil != e {
			log.Error(ctx, string(debug.Stack()))
			err = cause.Errorf("%v", e)
		}
	}()
	var ent types.Event
	if _, err = aware.Codec.Decode(bytes.NewBuffer(message.Body), &ent); nil != err {
		return cause.Errorf("Decode message - %s", err.Error())
	}
	rtx := &mpc.MeshContext{
		Context:     ctx,
		TraceId:     ent.Tid,
		SpanId:      ent.Sid,
		Timestamp:   time.Now().UnixMilli(),
		RunMode:     int(prsim.Routine),
		Attachments: map[string]string{},
	}
	mtx := mpc.Context()
	mtx.RewriteContext(rtx)
	mtx.SetAttribute(mpc.TimeoutKey, time.Second*5)
	if err := that.Subscribe(mtx, &ent); nil != err {
		return cause.Errorf("Exchange message - %s", err.Error())
	}
	return nil
}

func (that *nsqSubscriber) Subscribe(ctx context.Context, event *types.Event) error {
	if nil == that.subscriber {
		log.Info(ctx, "No subscriber named mesh exist, discard message %s, %s", event.Eid, event.Mid)
		return nil
	}
	if AllEventCode == event.Binding.Code {
		event.Binding.Code = "*"
	}
	return that.subscriber.Subscribe(ctx, event)
}
