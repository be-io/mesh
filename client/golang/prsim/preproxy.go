package prsim

import (
	"context"
)

var IPreProxy = (*PreProxy)(nil)

// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type PreProxy interface {

	// @MPI("mesh.glite2omega.api")
	ProxyInvokeGlite2Omega(ctx context.Context, param interface{}) (interface{}, error)
	// @MPI("mesh.omega2glite.api")
	ProxyInvokeOmega2Glite(ctx context.Context, param interface{}) (interface{}, error)
	// @MPI("mesh.glite.connectivity.api")
	InvokeGliteConnectivity(ctx context.Context, param interface{}) (interface{}, error)
	// @MPI("mesh.glite.healthcheck.api")
	InvokeHealthcheck(ctx context.Context, param interface{}) (interface{}, error)
}
