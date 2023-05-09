/*
 * Copyright (c) 2000, 2099, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package kms

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/be-io/mesh/client/golang/boost"
	"github.com/be-io/mesh/client/golang/cause"
	_ "github.com/be-io/mesh/client/golang/codec"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/be-io/mesh/client/golang/types"
	"math/big"
	"net"
	"os"
	"time"
)

func init() {
	var _ prsim.KMS = new(localKMS)
	macro.Provide(prsim.IKMS, new(localKMS))
}

const Local = "local"

type localKMS struct {
	env *types.Environ
}

func (that *localKMS) Att() *macro.Att {
	return &macro.Att{Name: Local}
}

func (that *localKMS) Environ(ctx context.Context) (*types.Environ, error) {
	if nil != that.env {
		return that.env, nil
	}
	mid := os.Getenv("MID")
	if "" == mid {
		return nil, cause.Errorf("Mesh key chain cant read from environment, check is present with MID")
	}
	env, err := boost.FromMeshID(ctx, mid)
	if nil != err {
		return nil, cause.Error(err)
	}
	that.env = env
	return env, nil
}

func (that *localKMS) List(ctx context.Context, cno string) ([]*types.Keys, error) {
	return nil, nil
}

func (that *localKMS) ApplyRoot(ctx context.Context, csr *types.KeyCsr) ([]*types.Keys, error) {
	csr.IsCA = true
	certPEM, keyPEM, err := that.KeyPair(csr)
	if nil != err {
		return nil, cause.Error(err)
	}
	var keys []*types.Keys
	keys = append(keys, &types.Keys{
		CNO:     csr.CNO,
		PNO:     csr.PNO,
		KNO:     tool.NextID(),
		Kind:    types.RootCaCrtKey,
		Key:     string(certPEM),
		Version: 0,
	})
	keys = append(keys, &types.Keys{
		CNO:     csr.CNO,
		PNO:     csr.PNO,
		KNO:     tool.NextID(),
		Kind:    types.RootCaPrivateKey,
		Key:     string(keyPEM),
		Version: 0,
	})
	return keys, nil
}

// ApplyIssue creates a self-signed certificate and key for the given host.
// Host may be an IP or a DNS name
// The certificate will be created with file mode 0644. The key will be created with file mode 0600.
// If the certificate or key files already exist, they will be overwritten.
// Any parent directories of the certPath or keyPath will be created as needed with file mode 0755.
func (that *localKMS) ApplyIssue(ctx context.Context, csr *types.KeyCsr) ([]*types.Keys, error) {
	certPEM, keyPEM, err := that.KeyPair(csr)
	if nil != err {
		return nil, cause.Error(err)
	}
	var keys []*types.Keys
	keys = append(keys, &types.Keys{
		CNO:     csr.CNO,
		PNO:     csr.PNO,
		KNO:     tool.NextID(),
		Kind:    types.IssueCrtKey,
		Key:     string(certPEM),
		Version: 0,
	})
	keys = append(keys, &types.Keys{
		CNO:     csr.CNO,
		PNO:     csr.PNO,
		KNO:     tool.NextID(),
		Kind:    types.IssuePrivateKey,
		Key:     string(keyPEM),
		Version: 0,
	})
	return keys, nil
}

// GenCertificate generates random TLS certificates.
func (that *localKMS) GenCertificate(csr *types.KeyCsr) (*tls.Certificate, error) {
	certPEM, keyPEM, err := that.KeyPair(csr)
	if nil != err {
		return nil, cause.Error(err)
	}

	certificate, err := tls.X509KeyPair(certPEM, keyPEM)
	if nil != err {
		return nil, cause.Error(err)
	}

	return &certificate, nil
}

// KeyPair generates cert and key files.
func (that *localKMS) KeyPair(csr *types.KeyCsr) ([]byte, []byte, error) {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, csr.Length)
	if nil != err {
		return nil, nil, cause.Error(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rsaPrivateKey)})

	certPEM, err := that.PemCert(rsaPrivateKey, csr)
	if nil != err {
		return nil, nil, cause.Error(err)
	}
	return certPEM, keyPEM, nil
}

// PemCert generates PEM cert file.
func (that *localKMS) PemCert(privateKey *rsa.PrivateKey, csr *types.KeyCsr) ([]byte, error) {
	derBytes, err := that.derCert(privateKey, csr)
	if nil != err {
		return nil, cause.Error(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes}), nil
}

func (that *localKMS) derCert(privateKey *rsa.PrivateKey, csr *types.KeyCsr) ([]byte, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if nil != err {
		return nil, cause.Error(err)
	}

	expiration := time.Time(csr.ExpireAt)
	if expiration.IsZero() {
		expiration = time.Now().Add(365 * (24 * time.Hour))
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{"CN"},
			Organization:       []string{csr.Domain},
			OrganizationalUnit: []string{csr.Domain},
			Province:           []string{"ZJ"},
			Locality:           []string{"HZ"},
			CommonName:         tool.Anyone(csr.Subject, csr.Domain),
		},
		NotBefore: time.Now(),
		NotAfter:  expiration,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement | x509.KeyUsageDataEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{csr.Domain},
		EmailAddresses:        []string{csr.Mail},

		SignatureAlgorithm: x509.SHA256WithRSA,
		IsCA:               csr.IsCA,
	}
	if len(csr.IPs) > 0 {
		for _, ip := range csr.IPs {
			template.IPAddresses = append(template.IPAddresses, net.ParseIP(ip))
		}
	}
	if !csr.IsCA {
		template.Issuer = pkix.Name{
			Country:            []string{"CN"},
			Organization:       []string{csr.Domain},
			OrganizationalUnit: []string{csr.Domain},
			Locality:           []string{"HZ"},
			Province:           []string{"ZJ"},
			CommonName:         "RootCA",
		}
		caCrtPem, _ := pem.Decode([]byte(csr.CaCert))
		caKeyPem, _ := pem.Decode([]byte(csr.CaKey))
		caCrt, err := x509.ParseCertificate(caCrtPem.Bytes)
		if nil != err {
			return nil, cause.Error(err)
		}
		caKey, err := x509.ParsePKCS1PrivateKey(caKeyPem.Bytes)
		if nil != err {
			return nil, cause.Error(err)
		}
		return x509.CreateCertificate(rand.Reader, &template, caCrt, &privateKey.PublicKey, caKey)
	} else {
		return x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	}
}
