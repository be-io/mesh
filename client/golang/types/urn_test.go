/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import (
	"fmt"
	"github.com/be-io/mesh/client/golang/macro"
	"testing"
)

func TestURNArray(t *testing.T) {
	urn := &URN{
		Name:   "${mesh.name}.process.task.callback",
		NodeId: "lx000000000000000000x",
		Flag: &URNFlag{
			V:       "00",
			Proto:   "00",
			Codec:   "00",
			Version: "1.2.3",
			Zone:    "00",
			Cluster: "00",
			Cell:    "00",
			Group:   "00",
			Address: "127.0.0.1",
			Port:    "4123",
		},
	}
	t.Log(urn.String())
}

func TestURNParse(t *testing.T) {
	ctx := macro.Context()
	urn := &URN{
		Name:   "com.omega.network.edge.accessible",
		NodeId: "lx000000000000000000x",
		Flag: &URNFlag{
			V:       "00",
			Proto:   "00",
			Codec:   "00",
			Version: "1.2.3",
			Zone:    "00",
			Cluster: "00",
			Cell:    "00",
			Group:   "00",
			Address: "127.0.0.1",
			Port:    "4123",
		},
	}
	t.Log(urn.String())

	urn0 := FromURN(ctx, urn.String())
	if urn0.Name != "com.omega.network.edge.accessible" {
		t.Error(fmt.Sprintf("URN name check failed with %s", urn0.Name))
	}
	if urn0.Flag.Address != "127.0.0.1" {
		t.Error(fmt.Sprintf("URN name check failed with %s", urn0.Flag.Address))
	}
	if urn0.Flag.Port != "4123" {
		t.Error(fmt.Sprintf("URN name check failed with %s", urn0.Flag.Port))
	}
	if urn0.Flag.Version != "1.2.3" {
		t.Error(fmt.Sprintf("URN name check failed with %s", urn0.Flag.Version))
	}

	if Substring("0", 0, 0) != "" {
		t.Error("Substring failed")
	}
	if Substring("0", 1, 1) != "" {
		t.Error("Substring failed")
	}
	if Substring("01", 1, 1) != "" {
		t.Error("Substring failed")
	}
	if Substring("0", 1, 2) != "" {
		t.Error("Substring failed")
	}
	if Substring("01", 1, 2) != "1" {
		t.Error("Substring failed")
	}
	if Substring("01", 1, 3) != "1" {
		t.Error("Substring failed")
	}
	if Substring("010", 1, 2) != "1" {
		t.Error("Substring failed")
	}
	if Substring("010", 0, 2) != "01" {
		t.Error("Substring failed")
	}
	t.Log(Reduce("0", 10))
	t.Log(Reduce("0", 1))

	name := "accessible.edge.network.omega.com.0000000102030000000012700000000104123.%s.trustbe.cn"
	if !FromURN(ctx, fmt.Sprintf(name, "lx1101011100010")).MatchNode(ctx, "lx1101011100010") {
		t.Error("Not match")
	}
	if !FromURN(ctx, fmt.Sprintf(name, "lx1101011100010")).MatchInst(ctx, "JG2021010101000001") {
		t.Error("Not match")
	}
	if !FromURN(ctx, fmt.Sprintf(name, "lx1101011100020")).MatchNode(ctx, "lx1101011100020") {
		t.Error("Not match")
	}
	if !FromURN(ctx, fmt.Sprintf(name, "lx1101011100020")).MatchInst(ctx, "JG2021010101000002") {
		t.Error("Not match")
	}
	if FromURN(ctx, fmt.Sprintf(name, "lx1101011100010")).MatchNode(ctx, "lx1101011100020") {
		t.Error("Not match")
	}
	if FromURN(ctx, fmt.Sprintf(name, "lx1101011100010")).MatchInst(ctx, "JG2021010101000002") {
		t.Error("Not match")
	}
	if FromURN(ctx, fmt.Sprintf(name, "lx1101011100020")).MatchNode(ctx, "lx1101011100010") {
		t.Error("Not match")
	}
	if FromURN(ctx, fmt.Sprintf(name, "lx1101011100020")).MatchInst(ctx, "JG2021010101000001") {
		t.Error("Not match")
	}

	releaseProduct := FromURN(ctx, "releaseProduct.server.base.0001100000000000000000000000000000080.lx0000010000010.trustbe.cn")
	t.Log(releaseProduct.MatchNode(ctx, "LX0000010000010"))
	t.Log(releaseProduct.MatchInst(ctx, "JG0100000100000000"))

	if len(FromURN(ctx, "releaseProduct.server.base.000110000000000000000000000.lx0000010000010.trustbe.cn").Flag.String()) != 37 {
		t.Error("Not compatible")
	}

	if len(FromURN(ctx, "releaseProduct.server.base.000110000000000000000000000000000008000.lx0000010000010.trustbe.cn").Flag.String()) != 37 {
		t.Error("Not compatible")
	}
	FromURN(ctx, "list.repository.data.asset.0001100000000000000012700000000139426.jg0100000500000000.trustbe.cn")
}
