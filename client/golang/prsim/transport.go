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

var ITransport = (*Transport)(nil)

// Transport
// Private compute data channel in async and blocking mode.
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Transport interface {

	// Open a channel session.
	// @MPI("mesh.chan.open")
	Open(ctx context.Context, sessionId string, metadata map[string]string) (Session, error)

	// Close the channel.
	// @MPI("mesh.chan.close")
	Close(ctx context.Context, timeout types.Duration) error

	// Roundtrip with the channel.
	// @MPI("mesh.chan.roundtrip")
	Roundtrip(ctx context.Context, payload []byte, metadata map[string]string) ([]byte, error)
}

var ISession = (*Session)(nil)

// Session
// Remote queue in async and blocking mode.
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Session interface {

	// Peek
	// Retrieves, but does not remove, the head of this queue,
	// or returns None if this queue is empty.
	// @MPI("mesh.chan.peek")
	Peek(ctx context.Context, topic string) ([]byte, error)

	// Pop
	// Retrieves and removes the head of this queue,
	// or returns None if this queue is empty.
	// @MPI("mesh.chan.pop")
	Pop(ctx context.Context, timeout types.Duration, topic string) ([]byte, error)

	// Push
	// Inserts the specified element into this queue if it is possible to do
	// so immediately without violating capacity restrictions.
	// When using a capacity-restricted queue, this method is generally
	// preferable to add, which can fail to insert an element only
	// by throwing an exception.
	// @MPI("mesh.chan.push")
	Push(ctx context.Context, payload []byte, metadata map[string]string, topic string) error

	// Release the channel session.
	// @MPI("mesh.chan.release")
	Release(ctx context.Context, timeout types.Duration, topic string) error
}
