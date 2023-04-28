/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package macro

import (
	"fmt"
	"time"
)

func init() {
	var _ MPS = new(MPSAnnotation)
}

// MPS
// Multi Attr Service. Mesh provider service.
type MPS interface {
	// Stt
	// Define service metadata.
	Stt() *Stt
}

type MPSAnnotation struct {
	Meta *Stt
}

func (that *MPSAnnotation) String() string {
	return fmt.Sprintf("%p", that)
}

func (that *MPSAnnotation) Stt() *Stt {
	return that.Meta
}

type Stt struct {
	// Service name. As alias topic.
	Name string
	// Service version.
	Version string
	// Service net/io protocol.
	Proto string
	// Service codec.
	Codec string
	// Service flag 1 asyncable 2 encrypt 4 communal.
	Flags int64
	// Service invoke timeout. millions.
	Timeout time.Duration
}
