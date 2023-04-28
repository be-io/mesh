/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package macro

import "fmt"

func init() {
	var _ MPI = new(MPIAnnotation)
}

type ServiceInspector interface {
	// Inspect multiple service plex.
	Inspect() []MPI
}

// MPI
// Multi Attr Interface. Mesh provider interface.
type MPI interface {
	// Rtt x
	// Define the reference metadata
	Rtt() *Rtt
}

type MPIAnnotation struct {
	Meta *Rtt
}

func (that *MPIAnnotation) Rtt() *Rtt {
	return that.Meta
}

type Rtt struct {
	// Service name. {@link MPS#name()}
	Name string
	// Service version. {@link MPS#version()}
	Version string
	// Service net/io protocol.
	Proto string
	// Service codec.
	Codec string
	// Service flag 1 asyncable 2 encrypt 4 communal.
	Flags int64
	// Service invoke timeout. millions.
	Timeout int64
	// Invoke retry times.
	Retries int
	// Service node identity.
	Node string
	// Service inst identity.
	Inst string
	// Service zone.
	Zone string
	// Service cluster.
	Cluster string
	// Service cell.
	Cell string
	// Service group.
	Group string
	// Service address.
	Address string
}

func (that *Rtt) String() string {
	return fmt.Sprintf("%s.%s", that.Name, that.Version)
}
