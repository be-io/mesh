/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package macro

import (
	"context"
	"github.com/be-io/mesh/client/golang/dsa"
	"reflect"
	"strings"
)

const (
	MeshSPI = "mesh-spi"
	MeshMPI = "mesh-mpi"
	MeshMPS = "mesh-mps"
	MeshSys = "mesh-sys"
	MeshNop = "mesh-nop"
)

var ark = dsa.NewStringMap[*ServiceLoader]()

func Provide(spi interface{}, provider SPI) {
	kind := Inspect(spi)
	ark.PutIfy(kind, func(k string) *ServiceLoader {
		return &ServiceLoader{name: provider.Att().Name, providers: dsa.NewStringMap[dsa.List[SPI]]()}
	}).Put(provider)
}

func Mock(spi interface{}, name string, mocking string) {
	Load(spi).PutInGroup(name, Load(spi).Get(mocking))
}

func Load(spi interface{}) *ServiceLoader {
	kind := Inspect(spi)
	return ark.PutIfy(kind, func(k string) *ServiceLoader {
		return &ServiceLoader{kind: spi, providers: dsa.NewStringMap[dsa.List[SPI]]()}
	})
}

func List() dsa.Map[string, *ServiceLoader] {
	return ark
}

func Inspect(spi interface{}) string {
	return reflect.TypeOf(spi).Elem().Name()
}

func Active() error {
	for _, providers := range ark.Values() {
		for _, provider := range providers.List() {
			if initializer, ok := provider.(Initializer); ok {
				if err := initializer.Init(); nil != err {
					return err
				}
			}
		}
	}
	return nil
}

type ServiceLoader struct {
	kind      any
	name      string
	providers dsa.Map[string, dsa.List[SPI]]
}

func (that *ServiceLoader) List() []SPI {
	var spx []SPI
	that.providers.ForEach(func(k string, providers dsa.List[SPI]) {
		if providers.Len() > 0 {
			spx = append(spx, providers.Gt(0))
		}
	})
	return spx
}

func (that *ServiceLoader) Default() any {
	return that.Get(that.name)
}

func (that *ServiceLoader) Name() string {
	return that.name
}

func (that *ServiceLoader) Methods(name string) map[any]map[string]*Method {
	if classify, ok := that.Get(name).(Interface); ok {
		return map[any]map[string]*Method{that.kind: classify.GetMethods()}
	}
	return nil
}

func (that *ServiceLoader) Get(name string) SPI {
	providers, ok := that.providers.Get(name)
	if !ok || providers.Len() < 1 {
		return nil
	}
	if providers.Peek().Att().Prototype {
		return providers.Peek().Att().Constructor()
	}
	return providers.Peek()
}

func (that *ServiceLoader) GetAny(name ...string) SPI {
	for _, n := range name {
		if !EableSPIFirst.Enable() && n == MeshSPI {
			continue
		}
		if sp := that.Get(n); nil != sp {
			return sp
		}
	}
	return nil
}

func (that *ServiceLoader) StartWith(name string) SPI {
	if ps := that.providers.FindAny(func(n string, v dsa.List[SPI]) bool {
		return strings.Index(name, n) == 0
	}); nil != ps {
		return ps.Peek()
	}
	return nil
}

func (that *ServiceLoader) Put(provider SPI) {
	that.PutInGroup(provider.Att().Name, provider)
	for _, name := range provider.Att().Alias {
		that.PutInGroup(name, provider)
	}
}

func (that *ServiceLoader) PutInGroup(name string, provider SPI) {
	that.providers.PutIfy(name, func(k string) dsa.List[SPI] {
		return dsa.NewSortList(func(i, j SPI) bool {
			return i.Att().Priority < j.Att().Priority
		})
	}).Ad(provider)
}

func Context() context.Context {
	return context.Background()
}
