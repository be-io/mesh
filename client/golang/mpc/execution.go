/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package mpc

import (
	"github.com/be-io/mesh/client/golang/macro"
)

type Generic interface {
	GetURN() string
	GetProto() string
	GetCodec() string
	GetTimeout() int64
	GetRetries() int
}

type Execution interface {

	// Schema is the execution schema.
	Schema() Generic

	// Inspect execution.
	Inspect() macro.Inspector

	Invoker
}
