/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package metabase

import (
	"fmt"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/plugin"
	"github.com/be-io/mesh/client/golang/types"
	"testing"
	"time"
)

func TestNetworkEdgeRefresh(t *testing.T) {
	ctx := mpc.Context()
	dsn := "root:@tcp(127.0.0.1:3306)/mesh"
	container := plugin.LoadC("metabase")
	container.Start(ctx, fmt.Sprintf("--dsn=%s", dsn))
	defer container.Stop(ctx)
	routes := `
		[
            {
                "node_id": "LX1101011100100",
                "inst_id": "JG0111001000000100",
                "name": "10.12.0.118",
                "address": "10.12.0.118:576"
            },
            {
                "node_id": "LX1101011100200",
                "inst_id": "JG0111002000000100",
                "name": "10.12.0.124",
                "address": "10.12.0.124:576"
            },
            {
                "node_id": "LX1101011100300",
                "inst_id": "JG0111003000000100",
                "name": "10.12.0.122",
                "address": "10.12.0.122:576"
            },
            {
                "node_id": "LX1101011100040",
                "inst_id": "JG0111000400000100",
                "name": "10.12.0.27",
                "address": "10.12.0.27:570"
            },
            {
                "node_id": "LX1101011100030",
                "inst_id": "JG0111000300000100",
                "name": "10.12.0.26",
                "address": "10.12.0.26:570"
            },
            {
                "node_id": "LX1101011100010",
                "inst_id": "JG0111000100000100",
                "name": "10.12.0.83",
                "address": "10.12.0.83:570"
            },
            {
                "node_id": "LX1101011100020",
                "inst_id": "JG0111000200000100",
                "name": "10.12.0.82",
                "address": "10.12.0.82:570"
            },
            {
                "node_id": "LX0000010000010",
                "inst_id": "JG0100000100000000",
                "name": "久弥集团A公司",
                "address": "gaia-mesh.dev-server-101:571"
            },
            {
                "node_id": "LX0000010000020",
                "inst_id": "JG0100000200000000",
                "name": "久弥集团B公司",
                "address": "gaia-mesh.dev-client-102:571"
            },
            {
                "node_id": "LX0000010000030",
                "inst_id": "JG0100000300000000",
                "name": "久弥集团C公司",
                "address": "gaia-mesh.dev-client-103:571"
            },
            {
                "node_id": "LX0000010000040",
                "inst_id": "JG0100000400000000",
                "name": "久弥集团D公司",
                "address": "gaia-mesh:571"
            },
            {
                "node_id": "LX0000010000050",
                "inst_id": "JG0100000500000000",
                "name": "久弥集团E公司",
                "address": "gaia-mesh:571"
            },
            {
                "node_id": "LX0000010000060",
                "inst_id": "JG0100000600000000",
                "name": "久弥集团F公司",
                "address": "gaia-mesh:571"
            }
        ]
`
	var x []*types.Route
	if err := codec.Jsonizer.Unmarshal([]byte(routes), &x); nil != err {
		t.Error(err)
		return
	}
	if err := lan.Refresh(mpc.Context(), x); nil != err {
		t.Error(err)
		return
	}
}

func TestTimeFormat(t *testing.T) {
	datetime, err := time.Parse(time.RFC3339, "2022-02-21T14:46:16.000Z")
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(datetime)
}
