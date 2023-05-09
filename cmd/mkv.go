/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package cmd

import (
	"context"
	"github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/spf13/cobra"
)

func init() {
	Provide(new(MKV))
}

type MKV struct {
}

func (that *MKV) Get(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "get",
		Version: prsim.Version,
		Short:   "Mesh kv get store value.",
		Long:    "Mesh kv get store value.",
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

func (that *MKV) Put(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "put",
		Version: prsim.Version,
		Short:   "Mesh kv put key value in store.",
		Long:    "Mesh kv put key value in store.",
		Run: func(cmd *cobra.Command, kvs []string) {
			if len(kvs) < 2 {
				log.Warn(cmd.Context(), "There is no key and value.")
				return
			}
			if err := aware.KV.Put(ctx, kvs[0], &types.Entity{
				Codec:  codec.JSON,
				Schema: "",
				Buffer: []byte(kvs[1]),
			}); nil != err {
				log.Error(ctx, err.Error())
			} else {
				log.Info(ctx, "OK")
			}
		},
	}
}

func (that *MKV) Remove(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm"},
		Version: prsim.Version,
		Short:   "Mesh kv remove store value.",
		Long:    "Mesh kv remove store value.",
		Run: func(cmd *cobra.Command, kvs []string) {
			if len(kvs) < 1 {
				log.Warn(cmd.Context(), "There is no key.")
				return
			}
			for _, key := range kvs {
				if err := aware.KV.Remove(ctx, key); nil != err {
					log.Error(ctx, err.Error())
				} else {
					log.Info(ctx, "OK")
				}
			}
		},
	}
}

func (that *MKV) Keys(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "keys",
		Version: prsim.Version,
		Short:   "Mesh kv key set.",
		Long:    "Mesh kv key set.",
		Run: func(cmd *cobra.Command, kvs []string) {

		},
	}
}

func (that *MKV) Home(ctx context.Context) *cobra.Command {
	mkv := &cobra.Command{
		Use:     "kv",
		Version: prsim.Version,
		Short:   "Mesh kv store.",
		Long:    "Mesh kv store.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Usage(); nil != err {
				log.Error(cmd.Context(), err.Error())
			}
		},
	}
	mkv.AddCommand(that.Get(ctx))
	mkv.AddCommand(that.Put(ctx))
	mkv.AddCommand(that.Remove(ctx))
	mkv.AddCommand(that.Keys(ctx))
	return mkv
}
