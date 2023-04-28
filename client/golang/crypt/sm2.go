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
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	"github.com/tjfoc/gmsm/sm2"
	"math/big"
	"time"
)

func init() {
	var _ prsim.Cryptor = new(sm2x)
	macro.Provide(prsim.ICryptor, new(sm2x))
}

const SM2 = "sm2"

type sm2x struct {
}

func (that *sm2x) Att() *macro.Att {
	return &macro.Att{Name: SM2}
}

func (that *sm2x) Encrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm2x) Decrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm2x) Hash(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm2x) Sign(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm2x) Verify(ctx context.Context, buff []byte, features map[string][]byte) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (that *sm2x) GenerateKeyPair(domain string) ([]*types.Keys, error) {
	privateKey, err := sm2.GenerateKey(rand.Reader)
	if nil != err {
		return nil, cause.Error(err)
	}
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: domain,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365 * 100),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	template.DNSNames = append(template.DNSNames, domain)

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if nil != err {
		return nil, cause.Error(err)
	}

	ceb := bytes.Buffer{}
	if err := pem.Encode(&ceb, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); nil != err {
		return nil, cause.Error(err)
	}

	pkb := bytes.Buffer{}
	//if err := pem.Encode(&pkb, &pem.Block{Type: "sm2x PRIVATE KEY", Bytes: x509.MarshalECPrivateKey(privateKey)}); nil != err {
	//	return nil, cause.Error(err)
	//}

	pub := bytes.Buffer{}
	//if err := pem.Encode(&pub, &pem.Block{Type: "sm2x PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)}); nil != err {
	//	return nil, cause.Error(err)
	//}

	var keys []*types.Keys
	keys = append(keys, &types.Keys{
		Key: pkb.String(),
	})
	keys = append(keys, &types.Keys{
		Key: pub.String(),
	})
	keys = append(keys, &types.Keys{
		Key: ceb.String(),
	})
	return keys, nil
}
