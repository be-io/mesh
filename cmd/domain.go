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
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	Provide(new(Domain))
}

type Domain struct {
}

func (that *Domain) Pub(ctx context.Context) *cobra.Command {
	var addr, node, kind string
	pub := &cobra.Command{
		Use:     "pub",
		Version: prsim.Version,
		Short:   "Mesh net domains pub.",
		Long:    "Mesh net domains pub.",
		Run: func(cmd *cobra.Command, args []string) {
			mtx := mpc.ContextWith(cmd.Context())
			if "" == addr || "" == node {
				log.Info(mtx, "Mesh net command must use valid arguments. ")
				return
			}
			mtx.SetAttribute(mpc.AddressKey, addr)
			mtx.GetPrincipals().Push(&types.Principal{NodeId: node, InstId: node})
			defer mtx.GetPrincipals().Pop()
			err := aware.Network.PutDomains(mtx, kind, []*types.Domain{{
				URN:     "",
				Address: "",
			}})
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			log.Info(mtx, "OK")
		},
	}
	pub.Flags().StringVarP(&node, "node", "n", types.LocalNodeId, "Mesh node or inst id.")
	pub.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1:570", "Mesh address.")
	pub.Flags().StringVarP(&kind, "kind", "k", prsim.AutoDomain, "Mesh domain mode.")
	return pub
}

func (that *Domain) Remove(ctx context.Context) *cobra.Command {
	var addr, node, kind string
	remove := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Version: prsim.Version,
		Short:   "Mesh net domains remove.",
		Long:    "Mesh net domains remove.",
		Run: func(cmd *cobra.Command, args []string) {
			mtx := mpc.ContextWith(cmd.Context())
			if "" == addr || "" == node {
				log.Info(mtx, "Mesh net command must use valid arguments. ")
				return
			}
			mtx.SetAttribute(mpc.AddressKey, addr)
			mtx.GetPrincipals().Push(&types.Principal{NodeId: node, InstId: node})
			defer mtx.GetPrincipals().Pop()
			err := aware.Network.PutDomains(mtx, kind, []*types.Domain{{
				URN:     "",
				Address: "",
			}})
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			log.Info(mtx, "OK")
		},
	}
	remove.Flags().StringVarP(&node, "node", "n", types.LocalNodeId, "Mesh node or inst id.")
	remove.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1:570", "Mesh address.")
	remove.Flags().StringVarP(&kind, "kind", "k", prsim.AutoDomain, "Mesh domain mode.")
	return remove
}

func (that *Domain) Home(ctx context.Context) *cobra.Command {
	var addr, node, kind string
	domain := &cobra.Command{
		Use:     "domain",
		Version: prsim.Version,
		Short:   "Remove mesh net domains.",
		Long:    "Remove mesh net domains.",
		Run: func(cmd *cobra.Command, kvs []string) {
			mtx := mpc.ContextWith(cmd.Context())
			if len(kvs) < 1 {
				log.Info(mtx, "Mesh domain command must use valid arguments. ")
				return
			}
			mtx.SetAttribute(mpc.AddressKey, addr)
			mtx.GetPrincipals().Push(&types.Principal{NodeId: node, InstId: node})
			defer mtx.GetPrincipals().Pop()
			domains, err := aware.Network.GetDomains(mtx, kind)
			if nil != err {
				log.Info(mtx, err.Error())
				return
			}
			table := pterm.TableData{{"Domain", "A"}}
			for _, domain := range domains {
				table = append(table, []string{
					domain.URN,
					domain.Address,
				})
			}
			if err = pterm.DefaultTable.WithHasHeader().WithBoxed(true).WithData(table).Render(); nil != err {
				log.Info(mtx, err.Error())
				return
			}
			pterm.Println()
		},
	}
	domain.AddCommand(that.Pub(ctx))
	domain.Flags().StringVarP(&node, "node", "n", types.LocalNodeId, "Mesh node or inst id.")
	domain.Flags().StringVarP(&addr, "addr", "a", "127.0.0.1:570", "Mesh address.")
	domain.Flags().StringVarP(&kind, "kind", "k", prsim.AutoDomain, "Mesh domain mode.")
	return domain
}
