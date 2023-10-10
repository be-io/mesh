/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package crypt

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"io"
)

func init() {
	var _ prsim.Cryptor = new(rsa2)
	macro.Provide(prsim.ICryptor, new(rsa2))
}

const RSA2 = "RSA2"

const (
	KeySize             = 2048
	KeyAlgorithm        = "RSA"
	SignatureAlgorithm  = "SHA256withRSA"
	RsaType             = "RSA/ECB/PKCS1Padding"
	MaxEncryptBlockSize = 244 // RSA2最大加密明文大小(2048/8-11=244)
	MaxDecryptBlockSize = 256 // RSA2最大解密密文大小(2048/8=256)
)

type rsa2 struct{}

func (that *rsa2) Att() *macro.Att {
	return &macro.Att{Name: RSA2}
}

func (that *rsa2) Encrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	pk, err := x509.ParsePKCS1PublicKey(features[string(prsim.PublicKey)])
	if nil != err {
		return nil, cause.Error(err)
	}
	var output bytes.Buffer
	r := bytes.NewBuffer(buff)
	for {
		buffer, err := io.ReadAll(io.LimitReader(r, MaxEncryptBlockSize))
		if nil != err {
			return nil, cause.Error(err)
		}
		if len(buffer) < 1 {
			break
		}
		cipher, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pk, buffer, nil)
		if nil != err {
			return nil, cause.Error(err)
		}
		if _, err := output.Write(cipher); nil != err {
			return nil, cause.Error(err)
		}
	}
	return output.Bytes(), nil
}

func (that *rsa2) Decrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	pk, err := x509.ParsePKCS1PrivateKey(features[string(prsim.PrivateKey)])
	if nil != err {
		return nil, cause.Error(err)
	}
	var output bytes.Buffer
	r := bytes.NewBuffer(buff)
	for {
		buffer, err := io.ReadAll(io.LimitReader(r, MaxDecryptBlockSize))
		if nil != err {
			return nil, cause.Error(err)
		}
		if len(buffer) < 1 {
			break
		}
		explain, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, pk, buffer, nil)
		if nil != err {
			return nil, cause.Error(err)
		}
		if _, err = output.Write(explain); nil != err {
			return nil, cause.Error(err)
		}
	}
	return output.Bytes(), nil
}

func (that *rsa2) Hash(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	return nil, nil
}

func (that *rsa2) Sign(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	pk, err := x509.ParsePKCS1PrivateKey(features[string(prsim.PrivateKey)])
	if nil != err {
		return nil, cause.Error(err)
	}
	hash := sha256.New()
	hash.Write(buff)
	hashed := hash.Sum(nil)
	signature, err := pk.Sign(rand.Reader, hashed, crypto.SHA256)
	if nil != err {
		return nil, cause.Error(err)
	}
	return signature, nil
}

func (that *rsa2) Verify(ctx context.Context, buff []byte, features map[string][]byte) (bool, error) {
	pk, err := x509.ParsePKCS1PublicKey(features[string(prsim.PublicKey)])
	if nil != err {
		return false, cause.Error(err)
	}
	sha2560 := sha256.New()
	sha2560.Write(buff)
	hashed := sha2560.Sum(nil)
	result := rsa.VerifyPKCS1v15(pk, crypto.SHA256, hashed, features[string(prsim.Signature)])
	return result == nil, nil
}
