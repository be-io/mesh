/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

import (
	"context"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/plugin"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"time"
)

func init() {
	var _ prsim.Listener = registers
	macro.Provide(prsim.IListener, registers)
}

var registers = new(register)

type register struct {
}

func (that *register) Att() *macro.Att {
	return &macro.Att{Name: prsim.ProxyRegisterEvent.Topic}
}

func (that *register) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.ProxyRegisterEvent}
}

func (that *register) Listen(ctx context.Context, event *types.Event) error {
	version, err := aware.Builtin.Version(ctx)
	if nil != err {
		log.Warn(ctx, "Inspect version, %s", err.Error())
		version = new(types.Versions)
	}
	if plugin.PROXY.Match() {
		that.Proxy(ctx, version)
	}
	if plugin.SERVER.Match() {
		that.Server(ctx, version)
	}
	return nil
}

func (that *register) Server(ctx context.Context, versions *types.Versions) {
	if tool.Proxy.Get().Empty() {
		return
	}
	registration := &types.Registration[any]{
		InstanceId: tool.SPA.Get().String(),
		Name:       tool.Name.Get(),
		Kind:       types.SERVER,
		Address:    tool.SPA.Get().String(),
		Timestamp:  time.Now().Add(time.Minute * 5).UnixMilli(),
		Content: &types.Metadata{
			References: []*types.Reference{},
			Services:   []*types.Service{},
		},
		Attachments: versions.Infos,
	}
	mtx := mpc.ContextWith(ctx).Resume(ctx)
	mtx.SetAttribute(mpc.AddressKey, tool.Proxy.Get())
	if err := aware.RemoteRegistry.Registers(mtx, []*types.Registration[any]{registration}); nil != err {
		log.Error(ctx, "Register server to proxy, %s", err.Error())
	}
}

func (that *register) Proxy(ctx context.Context, versions *types.Versions) {
	registration := &types.Registration[any]{
		InstanceId: tool.SPA.Get().String(),
		Name:       tool.Name.Get(),
		Kind:       types.PROXY,
		Address:    tool.SPA.Get().String(),
		Timestamp:  time.Now().Add(time.Minute * 5).UnixMilli(),
		Content: &types.Metadata{
			References: []*types.Reference{},
			Services:   []*types.Service{},
		},
		Attachments: versions.Infos,
	}
	for _, addr := range tool.Address.Get().All() {
		func() {
			mtx := mpc.ContextWith(ctx).Resume(ctx)
			mtx.SetAttribute(mpc.AddressKey, addr)
			if err := aware.RemoteRegistry.Registers(mtx, []*types.Registration[any]{registration}); nil != err {
				log.Error(ctx, err.Error())
			}
		}()
	}
	var registrations []*types.Registration[any]
	for _, addr := range tool.Address.Get().All() {
		mtx := mpc.ContextWith(ctx).Resume(ctx)
		mtx.SetAttribute(mpc.AddressKey, addr)
		v, err := InspectVersion(mtx, types.LocalNodeId)
		if nil != err {
			log.Error(ctx, err.Error())
			continue
		}
		registrations = append(registrations, &types.Registration[any]{
			InstanceId: addr,
			Name:       tool.Name.Get(),
			Kind:       types.SERVER,
			Address:    addr,
			Timestamp:  time.Now().Add(time.Minute * 5).UnixMilli(),
			Content: &types.Metadata{
				References: []*types.Reference{},
				Services:   []*types.Service{},
			},
			Attachments: v.Infos,
		})
	}
	if len(registrations) > 0 {
		if err := aware.LocalRegistry.Registers(ctx, registrations); nil != err {
			log.Error(ctx, err.Error())
		}
	}
}

func (that *register) KeepAlive(ctx context.Context) {
	log.Info(ctx, "Proxy active in %v mode, register self as proxy or server. ", plugin.Mode)
	topic := &types.Topic{Topic: prsim.ProxyRegisterEvent.Topic, Code: prsim.ProxyRegisterEvent.Code}
	if _, err := aware.Scheduler.Period(ctx, time.Second*30, topic); nil != err {
		log.Error(ctx, err.Error())
	}
}
