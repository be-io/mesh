/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"bytes"
	"context"
	"github.com/opendatav/mesh/client/golang/types"
)

var IPublisher = (*Publisher)(nil)

// Publisher
// Event publisher.
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Publisher interface {

	// Publish
	// @MPI("mesh.queue.publish")
	Publish(ctx context.Context, events []*types.Event) ([]string, error)

	// Broadcast
	// @MPI("mesh.queue.multicast")
	Broadcast(ctx context.Context, events []*types.Event) ([]string, error)
}

type Queue interface {

	// Publish Unicast will async publish the event to queue.
	Publish(ctx context.Context, topic types.Topic, buff *bytes.Buffer) (string, error)

	// Multicast will publish event to principal groups.
	Multicast(ctx context.Context, topic types.Topic, buff *bytes.Buffer) (string, error)

	// Broadcast Synchronized broadcast the event to all subscriber. This maybe timeout with to many subscriber.
	Broadcast(ctx context.Context, topic types.Topic, buff *bytes.Buffer) (string, error)
}
