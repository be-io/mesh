/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cmd

import (
	"context"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/plugin"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/spf13/cobra"
)

func init() {
	Provide(CommandFn(Inspect))
}

func Inspect(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "inspect",
		Version: prsim.Version,
		Short:   "Return low-level information on mesh plugin objects.",
		Long:    "Return low-level information on mesh plugin objects.",
		Run: func(cmd *cobra.Command, args []string) {
			yaml, err := plugin.Inspect()
			if nil != err {
				log.Error(ctx, err.Error())
				return
			}
			cmd.Print(string(yaml))
		},
	}
}
