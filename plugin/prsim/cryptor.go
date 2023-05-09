/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/crypt"
	"github.com/be-io/mesh/client/golang/prsim"
)

var _ prsim.Cryptor = new(PSRICryptor)

// PSRICryptor
// @MPS(macro.MeshMPS)
// @SPI(macro.MeshSPI)
type PSRICryptor struct {
}

func (that *PSRICryptor) Encrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	if nil == features {
		features = map[string][]byte{}
	}
	if nil == features[string(prsim.PublicKey)] {
		pk, err := that.DefaultPublicKey(ctx)
		if nil != err {
			return nil, cause.Error(err)
		}
		features[string(prsim.PublicKey)] = pk
	}
	return aware.Cryptor.Encrypt(ctx, buff, features)
}

func (that *PSRICryptor) Decrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	if nil == features {
		features = map[string][]byte{}
	}
	if nil == features[string(prsim.PrivateKey)] {
		pk, err := that.DefaultPrivateKey(ctx)
		if nil != err {
			return nil, cause.Error(err)
		}
		features[string(prsim.PrivateKey)] = pk
	}
	return aware.Cryptor.Decrypt(ctx, buff, features)
}

func (that *PSRICryptor) Hash(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	return aware.Cryptor.Hash(ctx, buff, features)
}

func (that *PSRICryptor) Sign(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	if nil == features {
		features = map[string][]byte{}
	}
	if nil == features[string(prsim.PrivateKey)] {
		pk, err := that.DefaultPrivateKey(ctx)
		if nil != err {
			return nil, cause.Error(err)
		}
		features[string(prsim.PrivateKey)] = pk
	}
	return aware.Cryptor.Sign(ctx, buff, features)
}

func (that *PSRICryptor) Verify(ctx context.Context, buff []byte, features map[string][]byte) (bool, error) {
	if nil == features {
		features = map[string][]byte{}
	}
	if nil == features[string(prsim.PublicKey)] {
		pk, err := that.DefaultPublicKey(ctx)
		if nil != err {
			return false, cause.Error(err)
		}
		features[string(prsim.PublicKey)] = pk
	}
	return aware.Cryptor.Verify(ctx, buff, features)
}

func (that *PSRICryptor) DefaultPublicKey(ctx context.Context) ([]byte, error) {
	environ, err := aware.LocalNet.GetEnviron(ctx)
	if nil != err {
		return nil, cause.Error(err)
	}
	nodeCrt, err := crypt.PEM.InformKey(environ.RootCrt)
	if nil != err {
		return nil, cause.Error(err)
	}
	certificate, err := x509.ParseCertificate(nodeCrt)
	if nil != err {
		return nil, cause.Error(err)
	}
	publicKey, ok := certificate.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, cause.Error(err)
	}
	return x509.MarshalPKCS1PublicKey(publicKey), nil
}

func (that *PSRICryptor) DefaultPrivateKey(ctx context.Context) ([]byte, error) {
	environ, err := aware.LocalNet.GetEnviron(ctx)
	if nil != err {
		return nil, cause.Errorc(cause.LicenseFormatError, err)
	}
	return crypt.PEM.InformKey(environ.RootKey)
}
