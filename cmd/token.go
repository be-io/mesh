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
	Provide(new(Token))
}

type Token struct {
}

func (that *Token) Apply(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "apply",
		Version: prsim.Version,
		Short:   "Mesh apply token in a duration.",
		Long:    "Mesh apply token in a duration.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func (that *Token) Verity(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "verify",
		Version: prsim.Version,
		Short:   "Mesh verify token in a duration.",
		Long:    "Mesh verify token in a duration.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func (that *Token) QuickAuth(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "quickauth",
		Aliases: []string{"qa"},
		Version: prsim.Version,
		Short:   "Mesh apply quick auth token.",
		Long:    "Mesh apply quick auth token.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func (that *Token) Grant(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "grant",
		Version: prsim.Version,
		Short:   "Mesh grant auth token.",
		Long:    "Mesh grant auth token.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func (that *Token) Accept(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "accept",
		Version: prsim.Version,
		Short:   "Mesh accept token auth.",
		Long:    "Mesh accept token auth.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func (that *Token) Reject(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "reject",
		Version: prsim.Version,
		Short:   "Mesh reject token auth.",
		Long:    "Mesh reject token auth.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func (that *Token) Authorize(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "authorize",
		Aliases: []string{"auth"},
		Version: prsim.Version,
		Short:   "Mesh auth token authorize.",
		Long:    "Mesh auth token authorize.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func (that *Token) Authenticate(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "authenticate",
		Version: prsim.Version,
		Short:   "Mesh auth token authenticate.",
		Long:    "Mesh auth token authenticate.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func (that *Token) Refresh(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:     "refresh",
		Version: prsim.Version,
		Short:   "Mesh auth token refresh.",
		Long:    "Mesh auth token refresh.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func (that *Token) Home(ctx context.Context) *cobra.Command {
	token := &cobra.Command{
		Use:     "token",
		Version: prsim.Version,
		Short:   "Mesh token center.",
		Long:    "Mesh token center.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Usage(); nil != err {
				log.Error(cmd.Context(), err.Error())
			}
		},
	}
	token.AddCommand(that.Apply(ctx))
	token.AddCommand(that.Verity(ctx))
	token.AddCommand(that.QuickAuth(ctx))
	token.AddCommand(that.Grant(ctx))
	token.AddCommand(that.Accept(ctx))
	token.AddCommand(that.Reject(ctx))
	token.AddCommand(that.Authorize(ctx))
	token.AddCommand(that.Authenticate(ctx))
	token.AddCommand(that.Refresh(ctx))
	return token
}
