/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cmd

import (
	"context"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/spf13/cobra"
)

func init() {
	Provide(new(Babel))
}

type Babel struct {
}

func (that *Babel) Home(ctx context.Context) *cobra.Command {
	var node, addr, lang, name string
	babel := &cobra.Command{
		Use:     "babel",
		Aliases: []string{"bb"},
		Version: prsim.Version,
		Short:   "Mesh babel.",
		Long:    "Mesh babel.",
		Run: func(cmd *cobra.Command, args []string) {
			mtx := mpc.ContextWith(cmd.Context())
			if "" == addr || "" == lang || "" == name || "" == node {
				log.Info(mtx, "Mesh babel command must use valid arguments. ")
				return
			}
			mtx.SetAttribute(mpc.AddressKey, addr)
			if types.LocalNodeId != node {
				mtx.GetPrincipals().Push(&types.Principal{NodeId: node, InstId: node})
			}
			r, err := aware.Devops.Distribute(mtx, &types.DistributeOption{
				Set:  name,
				Lang: lang,
				Addr: addr,
			})
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			log.Info(mtx, r)
		},
	}
	babel.Flags().StringVarP(&node, "node", "n", types.LocalNodeId, "Mesh node or inst id.")
	babel.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1:570", "Mesh address.")
	babel.Flags().StringVarP(&lang, "lang", "l", "javascript", "Specify babel language")
	babel.Flags().StringVarP(&name, "set", "s", "", "Specify babel set name")
	return babel
}
