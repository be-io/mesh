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
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/spf13/cobra"
)

func init() {
	Provide(CommandFn(Dump))
}

func Dump(ctx context.Context) *cobra.Command {
	dump := &cobra.Command{
		Use:     "dump",
		Version: prsim.Version,
		Short:   "Dump the mesh cluster metadata tables.",
		Long:    "Dump the mesh cluster metadata tables.",
		Run: func(cmd *cobra.Command, args []string) {
			registry, ok := macro.Load(prsim.IRegistry).Get(macro.MeshSPI).(prsim.Registry)
			if !ok {
				log.Warn(ctx, "No registry named %s. ", macro.MeshSPI)
				return
			}
			metadata, err := registry.Export(ctx, types.METADATA)
			if nil != err {
				log.Error(ctx, err.Error())
				return
			}
			log.Info(ctx, "%v", metadata)
		},
	}
	return dump
}
