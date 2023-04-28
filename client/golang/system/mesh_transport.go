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
	"sync"
)

func init() {
	var _ prsim.Transport = new(MeshTransport)
	var _ prsim.Session = new(MeshSession)
	macro.Provide(prsim.ITransport, new(MeshTransport))
	macro.Provide(prsim.ISession, new(MeshSession))
}

type MeshTransport struct {
	sessions map[string]prsim.Session
	sync.RWMutex
}

func (that *MeshTransport) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys, Alias: []string{macro.MeshSPI}}
}

func (that *MeshTransport) Open(ctx context.Context, sessionId string, metadata map[string]string) (prsim.Session, error) {
	if session := func() prsim.Session {
		that.RLock()
		defer that.RUnlock()
		if nil != that.sessions && nil != that.sessions[sessionId] {
			return that.sessions[sessionId]
		}
		return nil
	}(); nil != session {
		return session, nil
	}
	that.Lock()
	defer that.Unlock()
	if nil == that.sessions {
		that.sessions = map[string]prsim.Session{}
	}
	that.sessions[sessionId] = &MeshSession{SessionId: sessionId, Metadata: metadata}
	return that.sessions[sessionId], nil
}

func (that *MeshTransport) Close(ctx context.Context, timeout types.Duration) error {
	for sessionId, session := range that.sessions {
		if err := session.Release(ctx, timeout, ""); nil != err {
			log.Error(ctx, "Channel session %s release, %s.", sessionId, err.Error())
		}
	}
	that.sessions = map[string]prsim.Session{}
	return nil
}

func (that *MeshTransport) Roundtrip(ctx context.Context, payload []byte, metadata map[string]string) ([]byte, error) {
	return aware.Transport.Roundtrip(ctx, payload, metadata)
}

type MeshSession struct {
	SessionId string            `json:"session_id" xml:"session_id" yaml:"session_id"`
	Metadata  map[string]string `json:"metadata" xml:"metadata" yaml:"metadata"`
}

func (that *MeshSession) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *MeshSession) Peek(ctx context.Context, topic string) ([]byte, error) {
	return aware.Session.Peek(that.withMetadata(ctx), topic)
}

func (that *MeshSession) Pop(ctx context.Context, timeout types.Duration, topic string) ([]byte, error) {
	return aware.Session.Pop(that.withMetadata(ctx), timeout, topic)
}

func (that *MeshSession) Push(ctx context.Context, payload []byte, metadata map[string]string, topic string) error {
	return aware.Session.Push(that.withMetadata(ctx), payload, metadata, topic)
}

func (that *MeshSession) Release(ctx context.Context, timeout types.Duration, topic string) error {
	return aware.Session.Release(that.withMetadata(ctx), timeout, topic)
}

func (that *MeshSession) withMetadata(ctx context.Context) context.Context {
	mtx := mpc.ContextWith(ctx)
	for k, v := range that.Metadata {
		mtx.GetAttachments()[k] = v
	}
	prsim.MeshSessionId.Set(mtx.GetAttachments(), that.SessionId)
	return mtx
}
