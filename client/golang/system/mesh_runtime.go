/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.RuntimeAware = aware
	macro.Provide(prsim.IRuntimeAware, aware)
}

var aware = new(runtimeAware)

type runtimeAware struct {
	KMS       prsim.KMS
	Network   prsim.Network
	Cache     prsim.Cache
	JSON      codec.Codec
	Registry  prsim.Registry
	Sequence  prsim.Sequence
	Publisher prsim.Publisher
	Licenser  prsim.Licenser
	KV        prsim.KV
	Evaluator prsim.Evaluator
	Cluster   prsim.Cluster
	Session   prsim.Session
	Transport prsim.Transport
	Devops    prsim.Devops
}

func (that *runtimeAware) Att() *macro.Att {
	return &macro.Att{Name: "mesh.client.system.runtime"}
}

func (that *runtimeAware) Init() error {
	that.KMS = macro.Load(prsim.IKMS).GetAny(macro.MeshSPI, macro.MeshMPI).(prsim.KMS)
	that.Network = macro.Load(prsim.INetwork).GetAny(macro.MeshSPI, macro.MeshMPI).(prsim.Network)
	that.Cache = macro.Load(prsim.ICache).GetAny(macro.MeshSPI, macro.MeshMPI).(prsim.Cache)
	that.JSON = macro.Load(codec.ICodec).Get(codec.JSON).(codec.Codec)
	that.Registry = macro.Load(prsim.IRegistry).GetAny(macro.MeshSPI, macro.MeshMPI).(prsim.Registry)
	that.Sequence = macro.Load(prsim.ISequence).GetAny(macro.MeshSPI, macro.MeshMPI).(prsim.Sequence)
	that.Publisher = macro.Load(prsim.IPublisher).GetAny(macro.MeshSPI, macro.MeshMPI).(prsim.Publisher)
	that.Licenser = macro.Load(prsim.ILicenser).GetAny(macro.MeshSPI, macro.MeshMPI).(prsim.Licenser)
	that.KV = macro.Load(prsim.IKV).GetAny(macro.MeshSPI, macro.MeshMPI).(prsim.KV)
	that.Evaluator = macro.Load(prsim.IEvaluator).GetAny(macro.MeshSPI, macro.MeshMPI).(prsim.Evaluator)
	that.Cluster = macro.Load(prsim.ICluster).GetAny(macro.MeshSPI, macro.MeshMPI).(prsim.Cluster)
	that.Session = macro.Load(prsim.ISession).GetAny(macro.MeshSPI, macro.MeshMPI).(prsim.Session)
	that.Transport = macro.Load(prsim.ITransport).GetAny(macro.MeshSPI, macro.MeshMPI).(prsim.Transport)
	that.Devops = macro.Load(prsim.IDevops).GetAny(macro.MeshSPI, macro.MeshMPI).(prsim.Devops)
	return nil
}
