/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package crypt

import (
	"context"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.Cryptor = new(sm4x)
	macro.Provide(prsim.ICryptor, new(sm4x))
}

const SM4 = "sm4"

type sm4x struct {
}

func (that *sm4x) Encrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm4x) Decrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm4x) Hash(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm4x) Sign(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm4x) Verify(ctx context.Context, buff []byte, features map[string][]byte) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm4x) Att() *macro.Att {
	return &macro.Att{Name: SM4}
}
