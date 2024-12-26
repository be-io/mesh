/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package crypt

import (
	"context"
	"github.com/opendatav/mesh/client/golang/macro"
	"github.com/opendatav/mesh/client/golang/prsim"
)

func init() {
	var _ prsim.Cryptor = new(sm3x)
	macro.Provide(prsim.ICryptor, new(sm3x))
}

const SM3 = "sm3"

type sm3x struct {
}

func (that *sm3x) Encrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm3x) Decrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm3x) Hash(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm3x) Sign(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm3x) Verify(ctx context.Context, buff []byte, features map[string][]byte) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm3x) Att() *macro.Att {
	return &macro.Att{Name: SM3}
}
