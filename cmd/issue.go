/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cmd

import (
	"context"
	"github.com/be-io/mesh/client/golang/boost"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/spf13/cobra"
)

func init() {
	Provide(new(Issue))
}

type Issue struct {
}

func (that *Issue) Home(ctx context.Context) *cobra.Command {
	var addr, name, cname, kind, mid string
	issue := &cobra.Command{
		Use:     "issue",
		Aliases: []string{"is"},
		Version: prsim.Version,
		Short:   "Mesh issue.",
		Long:    "Mesh issue.",
		Run: func(cmd *cobra.Command, args []string) {
			mtx := mpc.ContextWith(cmd.Context())
			if "" != mid {
				env, err := boost.FromMeshID(mtx, mid)
				if nil != err {
					log.Info(mtx, err.Error())
					return
				}
				log.Info(mtx, env.Version)
				log.Info(mtx, env.NodeId)
				log.Info(mtx, env.InstId)
				log.Info(mtx, env.InstName)
				log.Info(mtx, env.RootCrt)
				log.Info(mtx, env.RootKey)
				log.Info(mtx, env.NodeCrt)
				return
			}
			if "" == name {
				log.Info(mtx, "Mesh issue command must use valid arguments. ")
				return
			}
			mtx.SetAttribute(mpc.AddressKey, addr)
			env, err := aware.Commercialize.Issued(mtx, name, kind, cname)
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			log.Info(mtx, env.Cipher)
		},
	}
	issue.Flags().StringVarP(&name, "name", "n", "", "Mesh node name")
	issue.Flags().StringVarP(&kind, "kind", "k", "160", "Mesh node kind")
	issue.Flags().StringVarP(&cname, "tech", "t", "LX", "Mesh node tech provider")
	issue.Flags().StringVarP(&mid, "mid", "m", "", "Mesh identity")
	issue.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1:570", "Mesh address.")
	return issue
}
