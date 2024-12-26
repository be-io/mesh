/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cmd

import (
	"context"
	"github.com/opendatav/mesh/client/golang/codec"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/tool"
	"github.com/opendatav/mesh/client/golang/types"
	"github.com/spf13/cobra"
)

func init() {
	Provide(CommandFn(Status))
}

func Status(ctx context.Context) *cobra.Command {
	var addr, node, format, sets string
	status := &cobra.Command{
		Use:     "status",
		Version: prsim.Version,
		Short:   "Display a live stream of mesh cluster resource usage statistics.",
		Long:    "Display a live stream of mesh cluster resource usage statistics.",
		Run: func(cmd *cobra.Command, args []string) {
			mtx := mpc.ContextWith(cmd.Context())
			if "" == addr || "" == node || "" == format {
				log.Info(mtx, "Mesh status command must use valid arguments. ")
				return
			}
			builtin, ok := macro.Load(prsim.IBuiltin).Get(macro.MeshMPI).(prsim.Builtin)
			if !ok {
				log.Info(mtx, "Mesh status command without builtin, please check the mesh version. ")
				return
			}
			encoder, ok := macro.Load(codec.ICodec).Get(tool.Anyone(format, codec.JSON)).(codec.Codec)
			if !ok {
				log.Info(mtx, "Mesh status command without codec named %s. ", tool.Anyone(format, codec.JSON))
				return
			}
			mtx.SetAttribute(mpc.RemoteName, sets)
			mtx.SetAttribute(mpc.AddressKey, addr)
			mtx.GetPrincipals().Push(&types.Principal{NodeId: node, InstId: node})
			defer mtx.GetPrincipals().Pop()
			stats, err := builtin.Stats(mtx, []string{})
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			output, err := encoder.EncodeString(stats)
			if nil != err {
				cmd.Println(err.Error())
				return
			}
			log.Info(mtx, "%s", output)
		},
	}
	status.Flags().StringVarP(&sets, "sets", "s", tool.Name.Get(), "Mesh sets name.")
	status.Flags().StringVarP(&node, "node", "n", types.LocalNodeId, "Mesh node or inst id.")
	status.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1:570", "Mesh address.")
	status.Flags().StringVarP(&format, "output", "o", "json", "Mesh invoke output format.")
	return status
}
