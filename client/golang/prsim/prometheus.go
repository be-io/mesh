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
)

// Prometheus
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Prometheus interface {

	// Range
	// @MPI("mesh.prom.monitor_query_range")
	Range(ctx context.Context, param *types.PromRangeQuery) (interface{}, error)

	// Tasks
	// @MPI("mesh.prom.task_offline.query_list")
	Tasks(ctx context.Context, param *types.TasksIndex) (interface{}, error)

	// Logs
	// @MPI("mesh.prom.task_offline.query_log")
	Logs(ctx context.Context, param *types.LogsIndex) (interface{}, error)

	// Proxy
	// @MPI("mesh.prom.proxy_get")
	Proxy(ctx context.Context, param *types.PromProxy) (interface{}, error)

	// Query
	// @MPI("mesh.prom.query")
	Query(ctx context.Context, param map[string]interface{}) (interface{}, error)

	// Range0
	// @MPI("mesh.prom.query_range")
	Range0(ctx context.Context, param map[string]interface{}) (interface{}, error)

	// Series
	// @MPI("mesh.prom.series")
	Series(ctx context.Context, param map[string]interface{}) (interface{}, error)

	// Labels
	// @MPI("mesh.prom.labels")
	Labels(ctx context.Context, param map[string]interface{}) (interface{}, error)

	// LabelValues
	// @MPI("mesh.prom.label.values")
	LabelValues(ctx context.Context, param *types.PromLabelValue) (interface{}, error)

	// Exemplars
	// @MPI("mesh.prom.query_exemplars")
	Exemplars(ctx context.Context, param map[string]interface{}) (interface{}, error)

	// Targets
	// @MPI("mesh.prom.targets")
	Targets(ctx context.Context) (interface{}, error)
}
