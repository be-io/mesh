/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"context"
	"fmt"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/dsa"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"strings"
	"time"
)

func init() {
	eden := &MeshEden{
		providers:  dsa.NewStringMap[*provider](),
		consumers:  dsa.NewAnyMap[any, dsa.Map[any, *consumer]](),
		indies:     dsa.NewAnyMap[any, dsa.Map[string, dsa.Map[string, Execution]]](),
		services:   dsa.NewStringMap[Execution](),
		references: dsa.NewStringMap[Execution](),
	}
	var _ Execution = new(instance)
	var _ Eden = eden
	var _ prsim.RuntimeHook = eden
	var _ prsim.Listener = eden
	macro.Provide(IEden, eden)
	macro.Provide(prsim.IRuntimeHook, eden)
	macro.Provide(prsim.IListener, eden)
}

type consumer struct {
	mpi       macro.MPI
	reference any
	proxy     any
}

type provider struct {
	making   bool
	mps      macro.MPS
	bindings []macro.Binding
	kind     any
	service  any
}

type MeshEden struct {
	providers  dsa.Map[string, *provider]
	consumers  dsa.Map[any, dsa.Map[any, *consumer]]
	indies     dsa.Map[any, dsa.Map[string, dsa.Map[string, Execution]]]
	services   dsa.Map[string, Execution]
	references dsa.Map[string, Execution]
}

func (that *MeshEden) Att() *macro.Att {
	return &macro.Att{
		Name: macro.MeshSPI,
	}
}

func (that *MeshEden) makeConsumer(mpi macro.MPI, reference any) *consumer {
	return that.consumers.PutIfy(reference, func(key any) dsa.Map[any, *consumer] {
		return dsa.NewAnyMap[any, *consumer]()
	}).PutIfy(mpi, func(key any) *consumer {
		return &consumer{
			mpi:       mpi,
			reference: reference,
			proxy:     macro.Load(reference).Get(macro.MeshMPI),
		}
	})
}

func (that *MeshEden) choosePresent(fn func(mpi macro.MPI) any, defaults any, as ...macro.MPI) any {
	for _, mac := range as {
		if nil != mac && nil != fn(mac) && "" != fn(mac) {
			return fn(mac)
		}
	}
	return defaults
}

func (that *MeshEden) makeReference(mpi macro.MPI, reference any, proxy any, env *types.Environ) dsa.Map[string, Execution] {
	return that.indies.PutIfy(reference, func(k any) dsa.Map[string, dsa.Map[string, Execution]] {
		return dsa.NewStringMap[dsa.Map[string, Execution]]()
	}).PutIfy(mpi.Rtt().String(), func(k string) dsa.Map[string, Execution] {
		executions := dsa.NewStringMap[Execution]()
		for kind, methods := range macro.Load(reference).Methods(macro.MeshMPI) {
			for _, method := range methods {
				refer := that.refMethod(env, mpi, kind, method)
				inst := &instance{
					urn:       refer.URN,
					kind:      reference,
					target:    proxy,
					resource:  refer,
					inspector: method,
					invoker:   &ServiceHandler{inspector: method},
				}
				that.references.Put(refer.URN, inst)
				executions.Put(method.String(), inst)
			}
		}
		return executions
	})
}

func (that *MeshEden) makeSchema(pv *provider) *types.Service {
	return &types.Service{
		URN:       "",
		Namespace: "",
		Name:      "",
		Version:   "",
		Proto:     MeshFlag.GRPC.Code(),
		Codec:     MeshFlag.JSON.Code(),
		Flags:     0,
		Timeout:   3000,
		Retries:   3,
		Node:      "",
		Inst:      "",
		Zone:      "",
		Cluster:   "",
		Cell:      "",
		Group:     "",
		Sets:      tool.Name.Get(),
		Address:   tool.Runtime.Get().String(),
		Kind:      tool.Ternary(nil == pv.mps, types.Binding, types.MPS),
		Lang:      "Golang",
		Attrs:     map[string]string{},
	}
}

func (that *MeshEden) makeServices(ctx context.Context, environ *types.Environ) dsa.Map[string, Execution] {
	for _, pv := range that.providers.Values() {
		if pv.making {
			continue
		}
		inspector, ok := macro.Load(pv.kind).Get(macro.MeshSPI).(macro.ServiceInspector)
		for _, methods := range macro.Load(pv.kind).Methods(macro.MeshMPI) {
			for _, method := range methods {
				var metas []macro.MPI
				if ok {
					metas = append(metas, inspector.Inspect()...)
				}
				if nil != pv.mps && !ok {
					metas = append(metas, method.GetMPI())
				}
				var schemas []*types.Service
				for _, binding := range pv.bindings {
					name := fmt.Sprintf("%s.%s", binding.Btt().Topic, binding.Btt().Code)
					schema := that.makeSchema(pv)
					schema.Namespace = binding.Btt().Topic
					schema.Name = binding.Btt().Code
					schema.Version = binding.Btt().Version
					schema.Flags = binding.Btt().Flags
					schema.URN = that.getURN(name, that.getServiceURNFlag(schema), environ)
					schemas = append(schemas, schema)
				}
				for _, meta := range metas {
					name := strings.ReplaceAll(meta.Rtt().Name, "${mesh.name}", tool.Name.Get())
					schema := that.makeSchema(pv)
					schema.Namespace = method.GetTName()
					schema.Name = method.GetName()
					schema.Version = pv.mps.Stt().Version
					schema.Flags = pv.mps.Stt().Flags | meta.Rtt().Flags
					schema.URN = that.getURN(name, that.getServiceURNFlag(schema), environ)
					schemas = append(schemas, schema)
				}
				that.makeService(ctx, schemas, method, pv)
			}
		}
		pv.making = true
	}
	return that.services
}

func (that *MeshEden) makeService(ctx context.Context, schemas []*types.Service, method *macro.Method, pv *provider) {
	for _, schema := range schemas {
		inspector := &macro.Method{
			DeclaredKind: method.GetDeclaredKind(),
			TName:        method.TName,
			Name:         method.Name,
			Bindings:     pv.bindings,
			Intype:       method.Intype,
			Retype:       method.Retype,
			Inbound:      method.Inbound,
			Outbound:     method.Outbound,
			MPS:          pv.mps,
			SPI:          method.SPI,
			MPI:          method.MPI,
		}
		that.services.PutIfy(types.FromURN(ctx, schema.URN).Name, func(k string) Execution {
			return &instance{
				urn:       schema.URN,
				kind:      method.GetDeclaredKind(),
				target:    pv.service,
				resource:  schema,
				inspector: inspector,
				invoker:   composite(&ServiceHandler{inspector: inspector}, PROVIDER),
			}
		})
	}
}

func (that *MeshEden) refMethod(environ *types.Environ, mpi macro.MPI, kind any, method *macro.Method) *types.Reference {
	om := method.GetMPI()
	reference := &types.Reference{
		URN:       "",
		Namespace: method.TName,
		Name:      method.Name,
		Version:   tool.Anyone(mpi.Rtt().Version, "", om.Rtt().Version),
		Proto:     tool.Anyone(mpi.Rtt().Proto, om.Rtt().Proto, MeshFlag.GRPC.Name()),
		Codec:     tool.Anyone(mpi.Rtt().Codec, om.Rtt().Codec, MeshFlag.JSON.Name()),
		Flags:     mpi.Rtt().Flags,
		Timeout:   mpi.Rtt().Timeout,
		Retries:   mpi.Rtt().Retries,
		Node:      mpi.Rtt().Node,
		Inst:      mpi.Rtt().Inst,
		Zone:      mpi.Rtt().Zone,
		Cluster:   mpi.Rtt().Cluster,
		Cell:      mpi.Rtt().Cell,
		Group:     mpi.Rtt().Group,
		Address:   mpi.Rtt().Address,
	}
	reference.URN = that.getURN(method.MPI.Rtt().Name, that.getReferenceURNFlag(reference), environ)
	return reference
}

func (that *MeshEden) getURN(alias string, flag *types.URNFlag, environ *types.Environ) string {
	urn := &types.URN{
		Domain: types.MeshDomain,
		NodeId: strings.ToLower(environ.NodeId),
		Name:   alias,
		Flag:   flag,
	}
	return urn.String()
}

func (that *MeshEden) getReferenceURNFlag(reference *types.Reference) *types.URNFlag {
	authority := strings.Split(reference.Address, ":")
	flag := &types.URNFlag{
		V:       "00",
		Proto:   MeshFlag.OfName(reference.Proto).Code(),
		Codec:   MeshFlag.OfName(reference.Codec).Code(),
		Version: reference.Version,
		Zone:    reference.Zone,
		Cluster: reference.Cluster,
		Cell:    reference.Cell,
		Group:   reference.Group,
	}
	if len(authority) > 0 {
		flag.Address = authority[0]
	}
	if len(authority) > 1 {
		flag.Port = authority[1]
	}
	return flag
}

func (that *MeshEden) getServiceURNFlag(reference *types.Service) *types.URNFlag {
	return &types.URNFlag{
		V:       "00",
		Proto:   MeshFlag.OfName(reference.Proto).Code(),
		Codec:   MeshFlag.OfName(reference.Codec).Code(),
		Version: reference.Version,
		Zone:    reference.Zone,
		Cluster: reference.Cluster,
		Cell:    reference.Cell,
		Group:   reference.Group,
		Address: "",
		Port:    "",
	}
}

func (that *MeshEden) newRegistration(ctx context.Context, metadata *types.Metadata) *types.Registration[any] {
	attachments := map[string]string{string(prsim.MeshSubset): tool.Subset.Get()}
	version, err := aware.Builtin.Version(ctx)
	if nil != err {
		log.Warn(ctx, "Inspect version, %s", err.Error())
		version = new(types.Versions)
	}
	for k, v := range version.Infos {
		attachments[k] = v
	}
	return &types.Registration[any]{
		InstanceId:  tool.Runtime.Get().String(),
		Name:        tool.Name.Get(),
		Kind:        types.METADATA,
		Address:     tool.Runtime.Get().String(),
		Timestamp:   time.Now().UnixMilli(),
		Content:     metadata,
		Attachments: attachments,
	}
}

func (that *MeshEden) Define(ctx context.Context, mpi macro.MPI, reference any) error {
	that.makeConsumer(mpi, reference)
	return nil
}

func (that *MeshEden) Refer(ctx context.Context, mpi macro.MPI, reference any, method macro.Inspector) (Execution, error) {
	environ, err := aware.Network.GetEnviron(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	cos := that.makeConsumer(mpi, reference)
	executions := that.makeReference(mpi, reference, cos.proxy, environ)
	execution, _ := executions.Get(method.String())
	return execution, nil
}

func (that *MeshEden) Store(ctx context.Context, kind any, service any) error {
	if binding, ok := service.(macro.Binding); ok {
		key := fmt.Sprintf("%s.%s.%s", binding.Btt().Topic, binding.Btt().Code, binding.Btt().Version)
		that.providers.PutIfy(key, func(k string) *provider {
			return &provider{
				making:   false,
				bindings: []macro.Binding{binding},
				kind:     kind,
				service:  service,
			}
		})
	}
	if mps, ok := service.(macro.MPS); ok {
		key := fmt.Sprintf("%s.%s", macro.Inspect(kind), tool.Anyone(mps.Stt().Version, "*"))
		that.providers.PutIfy(key, func(k string) *provider {
			return &provider{
				making:  false,
				mps:     mps,
				kind:    kind,
				service: service,
			}
		})
	}
	return nil
}

func (that *MeshEden) Infer(ctx context.Context, urn string) (Execution, error) {
	environ, err := aware.Network.GetEnviron(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	name := types.FromURN(ctx, urn).Name
	if service, ok := that.services.Get(name); ok {
		return service, nil
	}
	if service, ok := that.makeServices(ctx, environ).Get(name); ok {
		return service, nil
	}
	return &InspectExecution{URN: types.FromURN(ctx, urn)}, nil
}

func (that *MeshEden) ReferTypes(ctx context.Context) ([]macro.MPI, error) {
	var ms []macro.MPI
	for _, key := range that.consumers.Keys() {
		ms = append(ms, key.(macro.MPI))
	}
	return ms, nil
}

func (that *MeshEden) InferTypes(ctx context.Context) ([]macro.MPS, error) {
	var ms []macro.MPS
	for _, v := range that.providers.Values() {
		ms = append(ms, v.kind.(macro.MPS))
	}
	return ms, nil
}

func (that *MeshEden) Start(ctx context.Context, runtime prsim.Runtime) error {
	for _, providers := range macro.List().Values() {
		for _, pv := range providers.List() {
			if x, ok := pv.(macro.Interface); ok {
				for _, method := range x.GetMethods() {
					that.makeConsumer(method.MPI, method.GetDeclaredKind())
					for _, service := range macro.Load(method.GetDeclaredKind()).List() {
						log.Catch(that.Store(ctx, method.GetDeclaredKind(), service))
					}
				}
			}
		}
	}
	for _, subscriber := range macro.Load(prsim.ISubscriber).List() {
		log.Catch(that.Store(ctx, prsim.ISubscriber, subscriber))
	}
	log.Catch(that.Listen(ctx, &types.Event{}))
	topic := &types.Topic{Topic: prsim.MetadataRegisterEvent.Topic, Code: prsim.MetadataRegisterEvent.Code}
	_, err := aware.Scheduler.Period(ctx, time.Second*30, topic)
	return cause.Error(err)
}

func (that *MeshEden) Stop(ctx context.Context, runtime prsim.Runtime) error {
	if err := aware.Scheduler.Shutdown(ctx, time.Second*3); nil != err {
		log.Error(ctx, err.Error())
	}
	if err := aware.Registry.Unregister(ctx, that.newRegistration(ctx, &types.Metadata{})); nil != err {
		log.Error(ctx, err.Error())
	}
	return nil
}

func (that *MeshEden) Refresh(ctx context.Context, runtime prsim.Runtime) error {
	if err := that.Start(ctx, runtime); nil != err {
		log.Error(ctx, err.Error())
	}
	return nil
}

func (that *MeshEden) Btt() []*macro.Btt {
	return []*macro.Btt{prsim.MetadataRegisterEvent}
}

func (that *MeshEden) Listen(ctx context.Context, event *types.Event) error {
	environ, err := aware.Network.GetEnviron(ctx)
	if nil != err {
		return cause.Error(err)
	}
	that.makeServices(ctx, environ)
	for _, consumers := range that.consumers.Values() {
		for _, c := range consumers.Values() {
			that.makeReference(c.mpi, c.reference, c.proxy, environ)
		}
	}
	var references []*types.Reference
	for _, reference := range that.references.Values() {
		references = append(references, reference.Schema().(*types.Reference))
	}
	var services []*types.Service
	for _, service := range that.services.Values() {
		services = append(services, service.Schema().(*types.Service))
	}
	metadata := &types.Metadata{References: references, Services: services}
	return cause.Error(aware.Registry.Registers(ctx, []*types.Registration[any]{that.newRegistration(ctx, metadata)}))
}

type instance struct {
	urn       string
	kind      any
	target    any
	resource  Generic
	inspector macro.Inspector
	invoker   Invoker
}

func (that *instance) Inspect() macro.Inspector {
	return that.inspector
}

func (that *instance) Schema() Generic {
	return that.resource
}

func (that *instance) Invoke(ctx context.Context, invocation Invocation) (any, error) {
	return that.invoker.Invoke(ctx, invocation)
}
