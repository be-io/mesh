/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import "github.com/be-io/mesh/client/golang/types"

type Routable interface {

	// With
	// Route with attachments.
	With(key string, value string) Routable

	// WithMap
	// Route with attachments.
	WithMap(attachments map[string]string) Routable

	// Local
	// Invoke the service in local network.
	Local() interface{}

	// Any
	// Invoke the service in a network, it may be local or others.
	Any(principal *types.Principal) interface{}

	// Many
	// Invoke the service in many network, it may be local or others. Broadcast mode.
	Many(principals ...*types.Principal) []interface{}
}

type MpcStream struct {
	reference   interface{}
	attachments map[string]string
}

func (that *MpcStream) With(key string, value string) Routable {
	return that.WithMap(map[string]string{key: value})
}

func (that *MpcStream) WithMap(attachments map[string]string) Routable {
	if len(attachments) < 1 {
		return that
	}
	kvs := map[string]string{}
	for key, value := range that.attachments {
		kvs[key] = value
	}
	for key, value := range attachments {
		kvs[key] = value
	}
	return &MpcStream{reference: that.reference, attachments: kvs}
}

func (that *MpcStream) Local() interface{} {
	return that.reference
}

func (that *MpcStream) Any(principal *types.Principal) interface{} {
	if nil == principal || ("" == principal.NodeId && "" == principal.InstId) {

	}
	return that
}

func (that *MpcStream) Many(principals ...*types.Principal) []interface{} {
	var routes []interface{}
	for _, principal := range principals {
		routes = append(routes, that.Any(principal))
	}
	return routes
}

// Of
// Wrap a service with streamable ability.
func Of(reference interface{}) Routable {
	return &MpcStream{reference: reference}
}
