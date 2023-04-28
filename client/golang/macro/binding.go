/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package macro

import (
	"time"
)

func init() {
	var _ Binding = new(BindingAnnotation)
}

type Binding interface {
	Btt() *Btt
}

type Bindings interface {

	// Btt
	// Subscribe bindings.
	Btt() []*Btt
}

type Btt struct {
	// Event topic.
	Topic string
	// Event code.
	Code string
	// Event version.
	Version string
	// Service net/io protocol.
	Proto string
	// Service codec.
	Codec string
	// Event subscribe asyncable.
	Flags int64
	// Service invoke timeout. millions.
	Timeout time.Duration
}

type BindingAnnotation struct {
	Binding *Btt
}

func (that *BindingAnnotation) Btt() *Btt {
	return that.Binding
}
