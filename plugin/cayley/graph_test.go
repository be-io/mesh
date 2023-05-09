/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cayley

import (
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/plugin"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"testing"
)

func TestLink(t *testing.T) {
	ctx := macro.Context()
	container := plugin.LoadC(plugin.GDB)
	container.Start(ctx)
	defer container.Stop(ctx)
	graph := macro.Load(prsim.IGraph).Get(Name).(prsim.Graph)
	if err := graph.Link(ctx, []*types.Quad{{
		Name:      "X",
		Subject:   "I",
		Predicate: "Love",
		Object:    "Go",
		Label:     "V",
	}, {
		Name:      "X",
		Subject:   "Go",
		Predicate: "Like",
		Object:    "Me",
		Label:     "V",
	}}); nil != err {
		t.Error(err)
		return
	}
	r, err := graph.GraphQL(ctx, &types.MeshQL{Name: "X", Expr: `
	graph.V('I')
	   .out('Go')
	   .all()
	`})
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(r)
}
