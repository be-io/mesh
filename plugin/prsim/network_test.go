/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"fmt"
	_ "github.com/opendatav/mesh/client/golang/grpc"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/plugin"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/tool"
	"github.com/opendatav/mesh/client/golang/types"
	"github.com/opendatav/mesh/plugin/metabase"
	"testing"
	"time"
)

func TestEmpty(t *testing.T) {
	ns := ""
	subscribers := map[string]map[string]bool{}
	subscribers[ns] = map[string]bool{}
	subscribers[ns]["alliance.ctrl.pboc"] = true
	t.Log(subscribers)
}

func TestRefreshEdge(t *testing.T) {
	ctx := mpc.Context()
	container := plugin.LoadC(plugin.PRSIM)
	container.Start(ctx, fmt.Sprintf("--dsn=%s", "root:@tcp(127.0.0.1:3306)/mesh"))
	defer func() {
		container.Stop(ctx)
	}()

	environ, err := aware.LocalNet.GetEnviron(ctx)
	if nil != err {
		t.Error(err)
		return
	}
	for index := 0; index < 3; index++ {
		if err = aware.LocalNet.Refresh(ctx, []*types.Route{
			{
				NodeId:   environ.NodeId,
				InstId:   environ.InstId,
				Name:     environ.InstName,
				Address:  tool.IP.Get(),
				Describe: environ.InstName,
				AuthCode: "",
				ExpireAt: time.Now().UnixMilli(),
			},
		}); nil != err {
			t.Error(err)
		}
	}

	if err = aware.LocalNet.Refresh(ctx, []*types.Route{
		{
			NodeId:   "LX0000010000010",
			InstId:   "JG0100000100000000",
			Name:     "久弥集团A公司",
			Address:  "10.99.31.33:570",
			Describe: "久弥集团A公司",
			AuthCode: "",
			ExpireAt: time.Now().UnixMilli(),
		},
	}); nil != err {
		t.Error(err)
	}

}

func TestGetEdges(t *testing.T) {
	ctx := mpc.Context()
	container := plugin.LoadC(plugin.PRSIM)
	container.Plugins(metabase.Name)
	container.Start(ctx, fmt.Sprintf("--dsn=%s", "root:@tcp(127.0.0.1:3306)/mesh"))
	defer func() {
		container.Stop(ctx)
	}()
	edges, err := aware.LocalNet.GetRoutes(ctx)
	if nil != err {
		t.Error(err)
		return
	}
	for _, edge := range edges {
		t.Log(edge)
	}
}

func TestAccessible(t *testing.T) {
	ctx := mpc.Context()
	container := plugin.LoadC(plugin.PRSIM)
	container.Plugins(metabase.Name)
	container.Start(ctx, fmt.Sprintf("--dsn=%s", "root:@tcp(127.0.0.1:3306)/mesh"))
	defer func() {
		container.Stop(ctx)
	}()
	for _, nodeId := range []string{"LX0000010000010", "LX0000010000090"} {
		access, err := aware.LocalNet.Accessible(ctx, &types.Route{NodeId: nodeId})
		if nil != err {
			t.Error(err)
			return
		}
		t.Log(access)
	}
}

func TestWeave(t *testing.T) {
	ctx := mpc.Context()
	container := plugin.LoadC(plugin.PRSIM)
	container.Plugins(metabase.Name)
	container.Start(ctx, fmt.Sprintf("--dsn=%s", "root:@tcp(127.0.0.1:3306)/mesh"))
	defer func() {
		container.Stop(ctx)
	}()

	err := aware.LocalNet.Weave(ctx, &types.Route{
		NodeId:  "LX0000010000010",
		InstId:  "JG0100000100000000",
		Address: "10.99.1.33:570",
	})
	if nil != err {
		t.Error(err)
	}
}

func TestWeaveACK(t *testing.T) {
	ctx := mpc.Context()
	container := plugin.LoadC(plugin.PRSIM)
	container.Plugins(metabase.Name)
	container.Start(ctx, fmt.Sprintf("--dsn=%s", "root:@tcp(127.0.0.1:3306)/mesh"))
	defer func() {
		container.Stop(ctx)
	}()

	err := aware.LocalNet.Ack(ctx, &types.Route{
		NodeId: "LX0000010000010", InstId: "JG0100000100000000",
	})
	if nil != err {
		t.Error(err)
	}
}

func TestDisable(t *testing.T) {
	ctx := mpc.Context()
	container := plugin.LoadC(plugin.PRSIM)
	container.Plugins(metabase.Name)
	container.Start(ctx, fmt.Sprintf("--dsn=%s", "root:@tcp(127.0.0.1:3306)/mesh"))
	defer func() {
		container.Stop(ctx)
	}()

	err := aware.LocalNet.Disable(ctx, "LX0000010000010")
	if nil != err {
		t.Error(err)
	}
}

func TestEnable(t *testing.T) {
	ctx := mpc.Context()
	container := plugin.LoadC(plugin.PRSIM)
	container.Plugins(metabase.Name)
	container.Start(ctx, fmt.Sprintf("--dsn=%s", "root:@tcp(127.0.0.1:3306)/mesh"))
	defer func() {
		container.Stop(ctx)
	}()

	err := aware.LocalNet.Enable(ctx, "LX0000010000010")
	if nil != err {
		t.Error(err)
	}
}

func TestIndex(t *testing.T) {
	ctx := mpc.Context()
	container := plugin.LoadC(plugin.PRSIM)
	container.Plugins(metabase.Name)
	container.Start(ctx, fmt.Sprintf("--dsn=%s", "root:@tcp(127.0.0.1:3306)/mesh"))
	defer func() {
		container.Stop(ctx)
	}()

	page, err := aware.LocalNet.Index(ctx, &types.Paging{Index: 0, Limit: 10, Factor: map[string]interface{}{}})
	if nil != err {
		t.Error(err)
	}
	t.Log(page)

	for index := 0; index < 5; index++ {
		x, err := aware.LocalNet.Index(ctx, &types.Paging{Index: int64(index), Limit: 1, Factor: map[string]interface{}{}})
		if nil != err {
			t.Error(err)
		}
		t.Log(x)
	}
}

func TestVersion(t *testing.T) {
	mtx := mpc.ContextWith(mpc.Context())
	mtx.SetAttribute(mpc.TimeoutKey, time.Second*3)
	mtx.SetAttribute(mpc.AddressKey, "10.99.1.33:570")
	network, ok := macro.Load(prsim.INetwork).Get(macro.MeshMPI).(prsim.Network)
	if !ok {
		t.Errorf("No network named %s exist. ", macro.MeshMPI)
		return
	}
	version, err := network.Version(mtx, "LX0000010000010")
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(version)
}

func TestGetEdge(t *testing.T) {
	ctx := mpc.Context()
	container := plugin.LoadC(plugin.PRSIM)
	container.Plugins(metabase.Name)
	container.Start(ctx, fmt.Sprintf("--dsn=%s", "root:@tcp(127.0.0.1:3306)/mesh"))
	defer func() {
		container.Stop(ctx)
	}()

	edges, err := aware.LocalNet.GetRoute(ctx, "")
	if nil != err {
		t.Error(err)
	}
	t.Log(edges)
}
