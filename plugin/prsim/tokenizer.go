/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/system"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"time"
)

var _ prsim.Tokenizer = new(PRSITokenizer)

// PRSITokenizer
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PRSITokenizer struct {
	Name  string
	Kind  string
	Token string
}

func (that *PRSITokenizer) Apply(ctx context.Context, kind string, duration time.Duration) (string, error) {
	token := tool.NextID()
	if err := system.PutWithCache(ctx, aware.Cache, token, token, duration); nil != err {
		return "", cause.Error(err)
	}
	return token, nil
}

func (that *PRSITokenizer) Verify(ctx context.Context, token string) (bool, error) {
	entity, err := aware.Cache.Get(ctx, token)
	if nil != err {
		return false, cause.Error(err)
	}
	if err := aware.Cache.Remove(ctx, token); nil != err {
		return false, cause.Error(err)
	}
	return nil != entity, nil
}

func (that *PRSITokenizer) Quickauth(ctx context.Context, credential *types.Credential) (*types.AccessToken, error) {
	return aware.Tokenizer.Quickauth(ctx, credential)
}

func (that *PRSITokenizer) Grant(ctx context.Context, credential *types.Credential) (*types.AccessGrant, error) {
	return aware.Tokenizer.Grant(ctx, credential)
}

func (that *PRSITokenizer) Accept(ctx context.Context, code string) (*types.AccessCode, error) {
	return aware.Tokenizer.Accept(ctx, code)
}

func (that *PRSITokenizer) Reject(ctx context.Context, code string) error {
	return aware.Tokenizer.Reject(ctx, code)
}

func (that *PRSITokenizer) Authorize(ctx context.Context, code string) (*types.AccessToken, error) {
	return aware.Tokenizer.Authorize(ctx, code)
}

func (that *PRSITokenizer) Authenticate(ctx context.Context, token string) (*types.AccessID, error) {
	return aware.Tokenizer.Authenticate(ctx, token)
}

func (that *PRSITokenizer) Refresh(ctx context.Context, token string) (*types.AccessToken, error) {
	return aware.Tokenizer.Refresh(ctx, token)
}
