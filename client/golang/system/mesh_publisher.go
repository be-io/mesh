/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package system

import (
	"context"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
)

func init() {
	var _ prsim.Publisher = new(systemPublisher)
	macro.Provide(prsim.IPublisher, new(systemPublisher))
}

type systemPublisher struct {
}

func (that *systemPublisher) Att() *macro.Att {
	return &macro.Att{Name: macro.MeshSys}
}

func (that *systemPublisher) Publish(ctx context.Context, events []*types.Event) ([]string, error) {
	return aware.Publisher.Publish(ctx, events)
}

func (that *systemPublisher) Broadcast(ctx context.Context, events []*types.Event) ([]string, error) {
	return aware.Publisher.Broadcast(ctx, events)
}
