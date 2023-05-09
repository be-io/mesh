/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/plugin"
	_ "github.com/be-io/mesh/client/golang/proxy"
	"github.com/be-io/mesh/client/golang/prsim"
	_ "github.com/be-io/mesh/client/golang/system"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"strconv"
	"time"
)

func init() {
	var _ prsim.Listener = registryCaster
	macro.Provide(prsim.IListener, registryCaster)
	var _ prsim.Listener = tabledataCaster
	macro.Provide(prsim.IListener, tabledataCaster)
	plugin.Provide(new(PRSMPlugin))
}

var (
	registryCaster  = &registryNotifier{Binding: &macro.Btt{Topic: "mesh.plugin.prsim.registration", Code: "refresh"}}
	tabledataCaster = &tabledataNotifier{Binding: &macro.Btt{Topic: "mesh.plugin.prsim.tabledata", Code: "refresh"}}
)

type PRSMPlugin struct {
}

func (that *PRSMPlugin) Ptt() *plugin.Ptt {
	return &plugin.Ptt{Priority: -99, Name: plugin.PRSIM, Flags: PRSMPlugin{}, Create: func() plugin.Plugin {
		return that
	}}
}

func (that *PRSMPlugin) Start(ctx context.Context, runtime plugin.Runtime) {
	for _, spi := range macro.Load(prsim.IRuntimeHook).List() {
		if hook, ok := spi.(prsim.RuntimeHook); ok {
			if err := hook.Refresh(ctx, runtime); nil != err {
				log.Error(ctx, err.Error())
			}
		}
	}
	registryCaster.Alive(ctx)
	tabledataCaster.Alive(ctx)
}

func (that *PRSMPlugin) Stop(ctx context.Context, runtime plugin.Runtime) {

}

type registryNotifier struct {
	Binding *macro.Btt
}

func (that *registryNotifier) Att() *macro.Att {
	return &macro.Att{Name: that.Binding.Topic}
}

func (that *registryNotifier) Btt() []*macro.Btt {
	return []*macro.Btt{that.Binding}
}

func (that *registryNotifier) Listen(ctx context.Context, event *types.Event) error {
	return that.Notify(ctx, true)
}

func (that *registryNotifier) Notify(ctx context.Context, check bool) error {
	var registrations types.MetadataRegistrations
	registrations = append(registrations, that.exportRegistrations(ctx, types.METADATA, check)...)
	registrations = append(registrations, that.exportRegistrations(ctx, types.PROXY, check)...)
	registrations = append(registrations, that.exportRegistrations(ctx, types.SERVER, check)...)
	registrations = append(registrations, that.exportRegistrations(ctx, types.COMPLEX, check)...)
	return cause.Error(notify(ctx, prsim.RegistryEventRefresh, registrations))
}

func (that *registryNotifier) exportRegistrations(ctx context.Context, kind string, check bool) types.MetadataRegistrations {
	registrations, err := tool.Ternary(plugin.PROXY.Match() && kind == types.METADATA, aware.RemoteRegistry, aware.Registry).Export(ctx, kind)
	if nil != err {
		log.Error(ctx, "Export service %s registration failed, %s", kind, err.Error())
		return nil
	}
	if nil == registrations {
		log.Debug(ctx, "No service %s registration exist", kind)
		return nil
	}
	var mrs types.MetadataRegistrations
	buff, err := aware.Codec.Encode(registrations)
	if nil != err {
		log.Error(ctx, "Registration %s has unexpected structure, %s", kind, err.Error())
		return nil
	}
	if _, err = aware.Codec.Decode(buff, &mrs); nil != err {
		log.Error(ctx, "Registration %s has unexpected structure, %s", kind, err.Error())
		return nil
	}
	var instances types.MetadataRegistrations
	for _, registration := range mrs {
		timestamp := time.UnixMilli(registration.Timestamp)
		if timestamp.Before(time.Now().Add(-time.Minute * 5)) {
			log.Info(ctx, "Registration outdated in 5 minutes, discard it, %s/%d/%s", registration.Address, registration.Timestamp, registration.Kind)
			if err = aware.Registry.Unregister(ctx, registration.Any()); nil != err {
				log.Error(ctx, err.Error())
			}
			continue
		}
		if !plugin.PROXY.Match() && check && !tool.CheckAvailable(ctx, registration.Address) {
			if err = aware.Registry.Unregister(ctx, registration.Any()); nil != err {
				log.Error(ctx, err.Error())
			}
			continue
		}
		instances = append(instances, registration)
	}
	return instances
}

func (that *registryNotifier) Alive(ctx context.Context) {
	topic := &types.Topic{Topic: that.Binding.Topic, Code: that.Binding.Code}
	if _, err := aware.Scheduler.Period(ctx, time.Second*10, topic); nil != err {
		log.Error(ctx, err.Error())
	}
}

type tabledataNotifier struct {
	Binding *macro.Btt
}

func (that *tabledataNotifier) Att() *macro.Att {
	return &macro.Att{Name: that.Binding.Topic}
}

func (that *tabledataNotifier) Btt() []*macro.Btt {
	return []*macro.Btt{that.Binding}
}

func (that *tabledataNotifier) Listen(ctx context.Context, event *types.Event) error {
	if err := that.RouteAutoRefresh(ctx); nil != err {
		log.Error(ctx, err.Error())
	}
	return nil
}

func (that *tabledataNotifier) RouteAutoRefresh(ctx context.Context, rs ...*types.Route) error {
	routes, err := tool.Ternary(plugin.PROXY.Match(), aware.RemoteNet, aware.LocalNet).GetRoutes(ctx)
	if nil != err {
		return cause.Error(err)
	}
	routes = append(routes, rs...)
	return cause.Error(notify(ctx, prsim.NetworkRouteRefresh, routes))
}

func (that *tabledataNotifier) RouteManRefresh(ctx context.Context) error {
	routes, err := tool.Ternary(plugin.PROXY.Match(), aware.RemoteNet, aware.LocalNet).GetRoutes(ctx)
	if nil != err {
		return cause.Error(err)
	}
	return cause.Error(notify(ctx, &macro.Btt{Topic: prsim.NetworkRouteRefresh.Topic, Code: "man"}, routes))
}

func (that *tabledataNotifier) Alive(ctx context.Context) {
	topic := &types.Topic{Topic: that.Binding.Topic, Code: that.Binding.Code}
	if _, err := aware.Scheduler.Period(ctx, time.Minute*1, topic); nil != err {
		log.Error(ctx, err.Error())
	}
}

func proxyFlush(ctx context.Context) error {
	return aware.Scheduler.Emit(ctx, &types.Topic{Topic: prsim.RoutePeriodRefresh.Topic, Code: prsim.RoutePeriodRefresh.Code})
}

func notify(ctx context.Context, binding *macro.Btt, content interface{}) error {
	if nil == content {
		return nil
	}
	environ, err := aware.LocalNet.GetEnviron(ctx)
	if nil != err {
		return cause.Error(err)
	}
	buff, err := aware.Codec.Encode(content)
	if nil != err {
		return cause.Error(err)
	}
	if nil == buff {
		return nil
	}
	mtx := mpc.ContextWith(ctx)
	return aware.Listener.Listen(ctx, &types.Event{
		Version:   types.MessageVersion,
		Tid:       mtx.GetTraceId(),
		Sid:       mtx.GetSpanId(),
		Eid:       tool.NextID(),
		Mid:       tool.NextID(),
		Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		Source: &types.Principal{
			NodeId: environ.NodeId,
			InstId: environ.InstId,
		},
		Target: &types.Principal{
			NodeId: environ.NodeId,
			InstId: environ.InstId,
		},
		Binding: &types.Topic{
			Topic: binding.Topic,
			Code:  binding.Code,
		},
		Entity: &types.Entity{
			Codec:  codec.JSON,
			Schema: types.MessageVersion,
			Buffer: buff.Bytes(),
		},
	})
}
