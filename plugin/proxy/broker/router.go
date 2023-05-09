/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package broker

import (
	"fmt"
	"reflect"

	"google.golang.org/grpc"
)

type HandlerMap map[string]service

var _ grpc.ServiceRegistrar = HandlerMap(nil)

type service struct {
	desc    *grpc.ServiceDesc
	handler interface{}
}

func (that HandlerMap) RegisterService(desc *grpc.ServiceDesc, h interface{}) {
	ht := reflect.TypeOf(desc.HandlerType).Elem()
	st := reflect.TypeOf(h)
	if !st.Implements(ht) {
		panic(fmt.Sprintf("service %s: handler of type %v does not satisfy %v", desc.ServiceName, st, ht))
	}
	if _, ok := that[desc.ServiceName]; ok {
		panic(fmt.Sprintf("service %s: handler already registered", desc.ServiceName))
	}
	that[desc.ServiceName] = service{desc: desc, handler: h}
}

func (that HandlerMap) QueryService(name string) (*grpc.ServiceDesc, interface{}) {
	svc := that[name]
	return svc.desc, svc.handler
}

func (that HandlerMap) GetServiceInfo() map[string]grpc.ServiceInfo {
	ret := make(map[string]grpc.ServiceInfo, len(that))
	for _, svc := range that {
		methods := make([]grpc.MethodInfo, 0, len(svc.desc.Methods)+len(svc.desc.Streams))
		for _, mtd := range svc.desc.Methods {
			methods = append(methods, grpc.MethodInfo{Name: mtd.MethodName})
		}
		for _, mtd := range svc.desc.Streams {
			methods = append(methods, grpc.MethodInfo{
				Name:           mtd.StreamName,
				IsClientStream: mtd.ClientStreams,
				IsServerStream: mtd.ServerStreams,
			})
		}
		ret[svc.desc.ServiceName] = grpc.ServiceInfo{
			Methods:  methods,
			Metadata: svc.desc.Metadata,
		}
	}
	return ret
}

func (that HandlerMap) ForEach(fn func(desc *grpc.ServiceDesc, svr interface{})) {
	for _, svc := range that {
		fn(svc.desc, svc.handler)
	}
}
