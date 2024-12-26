/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/opendatav/mesh/client/golang/types"
	"runtime"
)

var (
	Version  = "1.5.0.0"
	CommitID = "59a06ccd"
	GOOS     = runtime.GOOS
	GOARCH   = runtime.GOARCH
)

var IBuiltin = (*Builtin)(nil)

// Builtin spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Builtin interface {

	// Doc export the documents.
	// @MPI("${mesh.name}.builtin.doc")
	Doc(ctx context.Context, name string, formatter string) (string, error)

	// Version will get the builtin application version.
	// @MPI("${mesh.name}.builtin.version")
	Version(ctx context.Context) (*types.Versions, error)

	// Debug set the application log level.
	// @MPI("${mesh.name}.builtin.debug")
	Debug(ctx context.Context, features map[string]string) error

	// Stats will collect health check stats.
	// @MPI("${mesh.name}.builtin.stats")
	Stats(ctx context.Context, features []string) (map[string]string, error)

	// Fallback is fallback service
	// @MPI("${mesh.name}.builtin.fallback")
	Fallback(ctx context.Context) error
}

var IHodor = (*Hodor)(nil)

type Hodor interface {

	// X only diff with Builtin
	X() string

	// Stats collect the system, application, process or thread status etc.
	Stats(ctx context.Context, features []string) (map[string]string, error)

	// Debug set the debug features.
	Debug(ctx context.Context, features map[string]string) error
}
