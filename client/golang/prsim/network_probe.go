package prsim

import (
	"context"
	"github.com/opendatav/mesh/client/golang/types"
)

var INetworkProbe = (*NetworkProbe)(nil)

// NetworkProbe spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type NetworkProbe interface {
	// Ping
	// @MPI("mesh.probe.v1.ping")
	Ping(ctx context.Context) (string, error)

	// HealthCheck
	// @MPI("mesh.probe.v1.healthcheck")
	HealthCheck(ctx context.Context) (*types.ProbeResponse, error)
}
