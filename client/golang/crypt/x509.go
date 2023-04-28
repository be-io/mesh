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
	var _ prsim.Cryptor = new(x509x)
	macro.Provide(prsim.ICryptor, new(x509x))
}

const X509 = "x509"

type x509x struct {
}

func (that *x509x) Att() *macro.Att {
	return &macro.Att{Name: X509}
}

func (that *x509x) Encrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *x509x) Decrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *x509x) Hash(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *x509x) Sign(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *x509x) Verify(ctx context.Context, buff []byte, features map[string][]byte) (bool, error) {
	//TODO implement me
	panic("implement me")
}
