/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import (
	"github.com/opendatav/mesh/client/golang/cause"
	"time"
)

const CacheVersion = "1.0.0"

type CacheEntity struct {
	Version string `index:"0" json:"version" xml:"version" yaml:"version" comment:"Cache version"`

	Entity *Entity `index:"5" json:"entity" xml:"entity" yaml:"entity" comment:"Cache entity"`

	Timestamp int64 `index:"10" json:"timestamp" xml:"timestamp" yaml:"timestamp" comment:"Cache timestamp"`

	Duration int64 `index:"15" json:"duration" xml:"duration" yaml:"duration" comment:"Cache expired duration"`

	Key string `index:"20" json:"key" xml:"key" yaml:"key" comment:"Cache key"`
}

func (that *CacheEntity) Wrap(key string, value interface{}, duration time.Duration) (*CacheEntity, error) {
	entity, err := new(Entity).Wrap(value)
	if nil != err {
		return nil, cause.Error(err)
	}
	return &CacheEntity{
		Version:   CacheVersion,
		Entity:    entity,
		Timestamp: time.Now().UnixMilli(),
		Duration:  duration.Milliseconds(),
		Key:       key,
	}, nil
}

func (that *CacheEntity) TryRead(ptr interface{}) error {
	if nil == that.Entity || !that.Entity.Present() {
		return nil
	}
	return that.Entity.TryReadObject(ptr)
}
