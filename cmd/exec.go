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
	Provide(CommandFn(Exec))
}

func Exec(ctx context.Context) *cobra.Command {
	var addr, method, input, format, node string
	exec := &cobra.Command{
		Use:     "exec",
		Version: prsim.Version,
		Short:   "Execute the mesh actor, Such mesh exec actor -i '{}'.",
		Long:    "Execute the mesh actor, Such mesh exec actor -i '{}'.",
		Run: func(cmd *cobra.Command, args []string) {
			mtx := mpc.ContextWith(cmd.Context())
			if "" == addr || "" == method || "" == input || "" == node {
				log.Info(mtx, "Mesh exec command must use valid arguments. ")
				return
			}
			dispatcher, ok := macro.Load(prsim.IDispatcher).Get(macro.MeshSPI).(prsim.Dispatcher)
			if !ok {
				log.Info(mtx, "Mesh exec command without dispatcher, please check the mesh version. ")
				return
			}
			encoder, ok := macro.Load(codec.ICodec).Get(tool.Anyone(format, codec.JSON)).(codec.Codec)
			if !ok {
				log.Info(mtx, "Mesh exec command without codec named %s. ", tool.Anyone(format, codec.JSON))
				return
			}
			var parameters map[string]interface{}
			if _, err := encoder.DecodeString(input, &parameters); nil != err {
				log.Info(mtx, err.Error())
				return
			}
			urn := types.FromURN(mtx, types.LocURN(mtx, method))
			urn.NodeId = tool.Anyone(node, types.LocalNodeId)
			mtx.SetAttribute(mpc.AddressKey, addr)
			re, err := dispatcher.InvokeLRG(mtx, urn.String(), parameters)
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			log.Info(mtx, "Mesh exec command %s@%s:", method, node)
			log.Info(mtx, "input:%s", input)
			if nil != re {
				output, err := encoder.EncodeString(re)
				if nil != err {
					cmd.Println(err.Error())
					return
				}
				log.Info(mtx, "output:%s", output)
			}
		},
	}
	exec.Flags().StringVarP(&node, "node", "n", types.LocalNodeId, "Mesh node or inst id.")
	exec.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1:570", "Mesh address.")
	exec.Flags().StringVarP(&method, "method", "m", "mesh.net.environ", "Mesh program interface name.")
	exec.Flags().StringVarP(&input, "input", "i", "{}", "Mesh program interface parameters.")
	exec.Flags().StringVarP(&format, "output", "o", "json", "Mesh invoke output format.")
	return exec
}
