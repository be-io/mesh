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
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"strings"
	"time"
)

const RegistryPrefix = "mesh.registry"

var _ prsim.Registry = new(PRSIRegistry)

// PRSIRegistry
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PRSIRegistry struct {
}

func (that *PRSIRegistry) Register(ctx context.Context, registration *types.Registration[any]) error {
	if "" == registration.Address {
		registration.Address = registration.InstanceId
	}
	if registration.Timestamp < time.Minute.Milliseconds()*30 {
		registration.Timestamp = registration.Timestamp + time.Now().UnixMilli()
	}
	text, _ := aware.Codec.EncodeString(registration)
	log.Debug(ctx, "Register event=%s", text)
	if nil == registration || "" == registration.InstanceId {
		return nil
	}
	isSelf := tool.Name.Get() == registration.Name && strings.Contains(registration.Address, "127.0.0.1")
	if !tool.CheckAvailable(ctx, registration.Address) && !isSelf {
		log.Warn(ctx, "Instance %s:%s unavailable", registration.Name, registration.InstanceId)
		return nil
	}
	log.Debug(ctx, "Instance %s:%s has been registered", registration.Name, registration.InstanceId)
	key := fmt.Sprintf("%s.%s.%s", registration.Kind, registration.Name, registration.InstanceId)
	entity, err := new(types.CacheEntity).Wrap(key, registration, time.Hour*24*30)
	if nil != err {
		return cause.Error(err)
	}
	if err = aware.Cache.HSet(ctx, RegistryPrefix, entity); nil != err {
		return cause.Error(err)
	}
	return cause.Error(registryCaster.Notify(ctx, false))
}

func (that *PRSIRegistry) Registers(ctx context.Context, registrations []*types.Registration[any]) error {
	for _, registration := range registrations {
		if err := that.Register(ctx, registration); nil != err {
			return cause.Error(err)
		}
	}
	return nil
}

func (that *PRSIRegistry) Unregister(ctx context.Context, registration *types.Registration[any]) error {
	key := fmt.Sprintf("%s.%s.%s", registration.Kind, registration.Name, registration.InstanceId)
	return cause.Error(aware.Cache.HDel(ctx, RegistryPrefix, key))
}

func (that *PRSIRegistry) Export(ctx context.Context, kind string) ([]*types.Registration[any], error) {
	keys, err := aware.Cache.HKeys(ctx, RegistryPrefix)
	if nil != err {
		return nil, cause.Error(err)
	}
	var registrations []*types.Registration[any]
	for _, key := range keys {
		if strings.Index(key, kind) != 0 {
			continue
		}
		rt, err := aware.Cache.HGet(ctx, RegistryPrefix, key)
		if nil != err {
			log.Error(ctx, err.Error())
			continue
		}
		if nil == rt || nil == rt.Entity {
			continue
		}
		registration := new(types.Registration[any])
		if err = rt.Entity.TryReadObject(registration); nil != err {
			log.Error(ctx, err.Error())
			continue
		}
		registrations = append(registrations, registration)
	}
	return registrations, nil
}
