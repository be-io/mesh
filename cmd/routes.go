/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/types"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"time"
)

func init() {
	Provide(new(Route))
}

type Route struct {
}

func (that *Route) Remove(ctx context.Context) *cobra.Command {
	var addr, node string
	remove := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Version: prsim.Version,
		Short:   "Remove mesh net route.",
		Long:    "Remove mesh net route.",
		Run: func(cmd *cobra.Command, kvs []string) {
			mtx := mpc.ContextWith(cmd.Context())
			if len(kvs) < 1 {
				log.Info(mtx, "Mesh net command must use valid arguments. ")
				return
			}
			mtx.SetAttribute(mpc.AddressKey, addr)
			mtx.GetPrincipals().Push(&types.Principal{NodeId: node, InstId: node})
			defer mtx.GetPrincipals().Pop()
			route, err := aware.Network.GetRoute(mtx, kvs[0])
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			if nil == route {
				log.Info(mtx, "Route %s not present.", kvs[0])
				return
			}
			route.Status = route.Status | int32(types.Removed)
			if err = aware.Network.Refresh(mtx, []*types.Route{route}); nil != err {
				log.Info(mtx, err.Error())
				return
			}
			log.Info(mtx, "OK")
		},
	}
	remove.Flags().StringVarP(&node, "node", "n", types.LocalNodeId, "Mesh node or inst id.")
	remove.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1:570", "Mesh address.")
	return remove
}

func (that *Route) VIP(ctx context.Context) *cobra.Command {
	var addr, node, remove, write string
	vip := &cobra.Command{
		Use:     "vip",
		Aliases: []string{"vp"},
		Version: prsim.Version,
		Short:   "Mesh net vip.",
		Long:    "Mesh net vip.",
		Example: "mesh net vip -w src_id:LX0000000000000=x://127.0.0.1",
		Run: func(cmd *cobra.Command, args []string) {
			mtx := mpc.ContextWith(cmd.Context())
			mtx.SetAttribute(mpc.AddressKey, addr)
			mtx.SetAttribute(mpc.RemoteUname, "mesh.dot.vip")
			mtx.GetPrincipals().Push(&types.Principal{NodeId: node, InstId: node})
			defer mtx.GetPrincipals().Pop()
			r, err := aware.Endpoint.Fuzzy(mtx, []byte(fmt.Sprintf("{\"r\":\"%s\",\"w\":\"%s\"}", remove, write)))
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			var vips []*types.VIP
			if _, err = aware.Codec.Decode(bytes.NewBuffer(r), &vips); nil != err {
				log.Info(mtx, err.Error())
				return
			}
			for _, vip := range vips {
				log.Info(mtx, fmt.Sprintf("%s:%s=%s(%s)", vip.Name, vip.Matcher, vip.Label, strings.Join(vip.Hosts, ",")))
			}
			log.Info(mtx, "OK")
		},
	}
	vip.Flags().StringVarP(&remove, "remove", "r", "", "Mesh vip remove.")
	vip.Flags().StringVarP(&write, "write", "w", "", "Mesh vip write.")
	vip.Flags().StringVarP(&node, "node", "n", types.LocalNodeId, "Mesh node or inst id.")
	vip.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1:570", "Mesh address.")
	return vip
}

func (that *Route) Trick(ctx context.Context) *cobra.Command {
	var addr, node, remove, write string
	vip := &cobra.Command{
		Use:     "trick",
		Aliases: []string{"tk"},
		Version: prsim.Version,
		Short:   "Mesh net trick.",
		Long:    "Mesh net trick.",
		Example: "mesh net trick -w doc=http://doc:7220/studio.route.doc;studio=http://studio:7200/studio.route.asset;socket=http://socket:9904/studio.route.socket;cube=http://cube:9902/studio.route.cube;jupyter=http://jupyter:9906/studio.route.jupyter;ruby=http://ruby:80/http.route.ruby",
		Run: func(cmd *cobra.Command, args []string) {
			mtx := mpc.ContextWith(cmd.Context())
			mtx.SetAttribute(mpc.AddressKey, addr)
			mtx.SetAttribute(mpc.RemoteUname, "mesh.dot.trick")
			mtx.GetPrincipals().Push(&types.Principal{NodeId: node, InstId: node})
			defer mtx.GetPrincipals().Pop()
			r, err := aware.Endpoint.Fuzzy(mtx, []byte(fmt.Sprintf("{\"r\":\"%s\",\"w\":\"%s\"}", remove, write)))
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			var tricks []*types.Trick
			if _, err = aware.Codec.Decode(bytes.NewBuffer(r), &tricks); nil != err {
				log.Info(mtx, err.Error())
				return
			}
			for _, trick := range tricks {
				log.Info(mtx, fmt.Sprintf("%s:%s+%s=%s://%s/%s", trick.Name, trick.Kind, trick.Service, trick.Proto, trick.Address, strings.Join(trick.Patterns, ",")))
			}
			log.Info(mtx, "OK")
		},
	}
	vip.Flags().StringVarP(&remove, "remove", "r", "", "Mesh vip remove.")
	vip.Flags().StringVarP(&write, "write", "w", "", "Mesh vip write.")
	vip.Flags().StringVarP(&node, "node", "n", types.LocalNodeId, "Mesh node or inst id.")
	vip.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1:570", "Mesh address.")
	return vip
}

func (that *Route) Home(ctx context.Context) *cobra.Command {
	var addr, node string
	route := &cobra.Command{
		Use:     "net",
		Version: prsim.Version,
		Short:   "Display a live stream of mesh cluster resource usage statistics.",
		Long:    "Display a live stream of mesh cluster resource usage statistics.",
		Run: func(cmd *cobra.Command, args []string) {
			mtx := mpc.ContextWith(cmd.Context())
			if "" == addr || "" == node {
				log.Info(mtx, "Mesh net command must use valid arguments. ")
				return
			}
			mtx.SetAttribute(mpc.AddressKey, addr)
			mtx.GetPrincipals().Push(&types.Principal{NodeId: node, InstId: node})
			defer mtx.GetPrincipals().Pop()
			routes, err := aware.Network.GetRoutes(mtx)
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			table := pterm.TableData{{"NodeId", "InstId", "Name", "Address", "Status", "ExpireAt", "AuthCode", "Extra", "Group"}}
			for _, route := range routes {
				table = append(table, []string{
					route.NodeId,
					route.InstId,
					route.Name,
					route.URC().String(),
					strconv.FormatInt(int64(route.Status), 10),
					time.UnixMilli(route.ExpireAt).Format(log.DateFormat),
					route.AuthCode,
					route.Extra,
					route.Group,
				})
			}
			if err = pterm.DefaultTable.WithHasHeader().WithBoxed(true).WithData(table).Render(); nil != err {
				log.Info(mtx, err.Error())
				return
			}
			pterm.Println()
		},
	}
	route.AddCommand(that.Remove(ctx))
	route.AddCommand(that.VIP(ctx))
	route.AddCommand(that.Trick(ctx))
	route.Flags().StringVarP(&node, "node", "n", types.LocalNodeId, "Mesh node or inst id.")
	route.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1:7304", "Mesh address.")
	return route
}
