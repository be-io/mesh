/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import "context"

var IDispatcher = (*Dispatcher)(nil)

// Dispatcher
// @SPI(macro.MeshSPI)
type Dispatcher interface {

	// Invoke with map param
	// In multi returns, it's an array.
	Invoke(ctx context.Context, urn string, param map[string]interface{}) ([]interface{}, error)

	// Invoke0 with generic param
	// In multi returns, it's an array.
	Invoke0(ctx context.Context, urn string, param interface{}) ([]interface{}, error)

	// InvokeLR with fewer returns
	// In multi returns, it will discard multi returns
	InvokeLR(ctx context.Context, urn string, param map[string]interface{}) (interface{}, error)

	// InvokeLRG with fewer returns in generic mode
	// In multi returns, it will discard multi returns
	InvokeLRG(ctx context.Context, urn string, param interface{}) (interface{}, error)
}
