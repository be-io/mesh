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
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/spf13/cobra"
)

func init() {
	Provide(new(MOS))
}

type MOS struct {
}

func (that *MOS) Stat(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "stat",
		Version: prsim.Version,
		Short:   "Mesh OS 🐳 stat.",
		Long:    "Mesh OS 🐳 stat.",
		Run: func(cmd *cobra.Command, kvs []string) {
			if len(kvs) < 1 {
				log.Warn(cmd.Context(), "There is no key.")
				return
			}
			for _, key := range kvs {
				if entity, err := aware.KV.Get(ctx, key); nil != err {
					log.Error(ctx, err.Error())
				} else {
					log.Info(ctx, entity.String())
				}
			}
		},
	}
}

func (that *MOS) Join(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "join",
		Aliases: []string{"j"},
		Version: prsim.Version,
		Short:   "Mesh OS node join 🐳.",
		Long:    "Mesh OS node join 🐳.",
		Run: func(cmd *cobra.Command, kvs []string) {

		},
	}
}

func (that *MOS) Uninstall(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "uninstall",
		Aliases: []string{"u"},
		Version: prsim.Version,
		Short:   "Mesh OS uninstall 🐳.",
		Long:    "Mesh OS uninstall 🐳.",
		Run: func(cmd *cobra.Command, kvs []string) {

		},
	}
}

func (that *MOS) Install(ctx context.Context) *cobra.Command {
	i := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Version: prsim.Version,
		Short:   "Mesh OS install 🐳.",
		Long:    "Mesh OS install 🐳.",
		Run: func(cmd *cobra.Command, kvs []string) {

		},
	}
	return i
}

func (that *MOS) Home(ctx context.Context) *cobra.Command {
	mkv := &cobra.Command{
		Use:     "os",
		Version: prsim.Version,
		Short:   "Mesh OS 🚀🚀🚀.",
		Long:    "Mesh OS 🚀🚀🚀.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Usage(); nil != err {
				log.Error(cmd.Context(), err.Error())
			}
		},
	}
	mkv.AddCommand(that.Install(ctx))
	mkv.AddCommand(that.Uninstall(ctx))
	mkv.AddCommand(that.Join(ctx))
	mkv.AddCommand(that.Stat(ctx))
	return mkv
}
