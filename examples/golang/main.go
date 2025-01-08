/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/client/golang/prsim"
	"github.com/opendatav/mesh/client/golang/tool"
	"github.com/opendatav/mesh/ptp"
	"github.com/spf13/pflag"
	"io"
	"net/http"
)

func main() {
	var address, nodeId, suit string
	pflag.StringVarP(&address, "address", "a", "127.0.0.1:7304", "")
	pflag.StringVarP(&nodeId, "nodeId", "n", "YL070", "")
	pflag.StringVarP(&suit, "suit", "s", "A", "")
	pflag.Parse()
	ctx := mpc.Context()
	if err := do(ctx, &Suit{
		NodeId: nodeId,
		Addr:   address,
		CA:     tool.Ternary("A" == suit, BRootCrt, ARootCrt),
		CRT:    tool.Ternary("A" == suit, ABClientCrt, BAClientCrt),
		KEY:    tool.Ternary("A" == suit, ABClientKey, BAClientKey),
	}); nil != err {
		log.Error(ctx, "%s", err)
	}
}

type Suit struct {
	NodeId string
	Addr   string
	CA     string
	CRT    string
	KEY    string
}

func do(ctx context.Context, suit *Suit) error {
	keys, err := tls.X509KeyPair([]byte(suit.CRT), []byte(suit.KEY))
	if nil != err {
		return cause.Error(err)
	}
	ca := x509.NewCertPool()
	ca.AppendCertsFromPEM([]byte(suit.CA))
	cc := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
		RootCAs:            ca,
		Certificates:       []tls.Certificate{keys},
	}}}
	push := &ptp.PushInbound{
		Topic:    "1",
		Payload:  []byte("1"),
		Metadata: map[string]string{},
	}
	ib, err := ptp.Encode(push, "application/json")
	if nil != err {
		return cause.Error(err)
	}
	payload := &ptp.Inbound{
		Metadata: map[string]string{},
		Payload:  ib,
	}
	buf, err := ptp.Encode(payload, "application/json")
	if nil != err {
		return cause.Error(err)
	}
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s/v1/interconn/chan/invoke", suit.Addr), bytes.NewBuffer(buf))
	if nil != err {
		return cause.Error(err)
	}
	prsim.ContentType.SetHeader(r.Header, "application/json")
	prsim.MeshTargetNodeId.SetHeader(r.Header, suit.NodeId)
	prsim.MeshURI.SetHeader(r.Header, "/v1/interconn/chan/push")
	s, err := cc.Do(r)
	if nil != err {
		return cause.Error(err)
	}
	defer func() { log.Catch(s.Body.Close()) }()

	b, err := io.ReadAll(s.Body)
	if nil != err {
		return cause.Error(err)
	}
	log.Info(ctx, "%s", string(b))
	return nil
}

const (
	ARootCrt = `
-----BEGIN CERTIFICATE-----
MIICUzCCAfqgAwIBAgIRAMSDp/XH5DlwMCkDNs57w2QwCgYIKoZIzj0EAwIwYDEL
MAkGA1UEBhMCQ04xCzAJBgNVBAgTAlpKMQswCQYDVQQHEwJIWjEPMA0GA1UEChMG
eWwuY29tMQ8wDQYDVQQLEwZ5bC5jb20xFTATBgNVBAMMDOS4reWbvemTtuiBlDAg
Fw0yNTAxMDcwMzQyMDFaGA8yNTI1MDEwNzAzNDIwMVowYDELMAkGA1UEBhMCQ04x
CzAJBgNVBAgTAlpKMQswCQYDVQQHEwJIWjEPMA0GA1UEChMGeWwuY29tMQ8wDQYD
VQQLEwZ5bC5jb20xFTATBgNVBAMMDOS4reWbvemTtuiBlDBZMBMGByqGSM49AgEG
CCqGSM49AwEHA0IABEwt0od9U1kTZbqvpRdeTBbWpY6EwJXjWJ8c8Lu97evRSqMU
VW+Id9CI9bKHXTsgmOsqqE3aoghmdeDnL4P0a+WjgZIwgY8wDgYDVR0PAQH/BAQD
AgK8MCcGA1UdJQQgMB4GCCsGAQUFBwMBBggrBgEFBQcDAgYIKwYBBQUHAwEwDwYD
VR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU1eMlG4OFON+79tma8WlB6cFZt+4wJAYD
VR0RBB0wG4IGeWwuY29tgQtjYUBpY2JjLmNvbYcErBAYUTAKBggqhkjOPQQDAgNH
ADBEAiAi1Xkw9VMNdPuOu22qF60BMTeTI9umAe9TmFTWoXtlFQIgEPiCzfK5xyD0
07Zw+8WDX/2zZiaRza9nJ7WNxOacSjc=
-----END CERTIFICATE-----
	`
	BAClientCrt = `
-----BEGIN CERTIFICATE-----
MIICZTCCAgygAwIBAgIRAJnJYteogHAhGOMF/YB88QYwCgYIKoZIzj0EAwIwZDEL
MAkGA1UEBhMCQ04xCzAJBgNVBAgTAlpKMQswCQYDVQQHEwJIWjERMA8GA1UEChMI
aWNiYy5jb20xETAPBgNVBAsTCGljYmMuY29tMRUwEwYDVQQDDAzlt6XllYbpk7bo
oYwwIBcNMjUwMTA3MDM0MTA0WhgPMjUyNTAxMDcwMzQxMDRaMG0xCzAJBgNVBAYT
AkNOMQswCQYDVQQIEwJaSjELMAkGA1UEBxMCSFoxETAPBgNVBAoTCGljYmMuY29t
MREwDwYDVQQLEwhpY2JjLmNvbTEeMBwGA1UEAwwV5bel5ZWG6ZO26KGM5a6i5oi3
56uvMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAExRD7cUTn2me8vyZiyFn1ylBA
ScgLdP2PwyJFBprzlT10lX7B9o2AE0mF0Z7eZGWiCjs4Tv7Uxf6idZwUiOW0oKOB
kzCBkDAOBgNVHQ8BAf8EBAMCArwwJwYDVR0lBCAwHgYIKwYBBQUHAwEGCCsGAQUF
BwMCBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMB8GA1UdIwQYMBaAFMMsIrLOGJwb
Xe3g2MP6D+6ex1+1MCYGA1UdEQQfMB2CCGljYmMuY29tgQtjYUBpY2JjLmNvbYcE
rBAYUTAKBggqhkjOPQQDAgNHADBEAiBgLlkn3sC57pJSGCywbGvPeRhxqeKgJ5Zk
JSH7tgzhMwIgfsXRUyhrQUbxx1Bx7fBpYSTAXH2q4zjtMMeYdwMVFM8=
-----END CERTIFICATE-----
	`
	BAClientKey = `
-----BEGIN ECDSA PRIVATE KEY-----
MHcCAQEEIETzmtzBBT36x0DENWoNNbw0nEbzTkbsZlQhpGineNi8oAoGCCqGSM49
AwEHoUQDQgAExRD7cUTn2me8vyZiyFn1ylBAScgLdP2PwyJFBprzlT10lX7B9o2A
E0mF0Z7eZGWiCjs4Tv7Uxf6idZwUiOW0oA==
-----END ECDSA PRIVATE KEY-----
	`
	BRootCrt = `
-----BEGIN CERTIFICATE-----
MIICXDCCAgOgAwIBAgIQU4WNU06u8BdJ+uYIWPpDlzAKBggqhkjOPQQDAjBkMQsw
CQYDVQQGEwJDTjELMAkGA1UECBMCWkoxCzAJBgNVBAcTAkhaMREwDwYDVQQKEwhp
Y2JjLmNvbTERMA8GA1UECxMIaWNiYy5jb20xFTATBgNVBAMMDOW3peWVhumTtuih
jDAgFw0yNTAxMDcwMzM2MDhaGA8yNTI1MDEwNzAzMzYwOFowZDELMAkGA1UEBhMC
Q04xCzAJBgNVBAgTAlpKMQswCQYDVQQHEwJIWjERMA8GA1UEChMIaWNiYy5jb20x
ETAPBgNVBAsTCGljYmMuY29tMRUwEwYDVQQDDAzlt6XllYbpk7booYwwWTATBgcq
hkjOPQIBBggqhkjOPQMBBwNCAARnuX9+b8VNDBgnFfOsexxdh7ZUfyfN+y1zRlyC
oFj/+8BqqlN4OIRi3d0Bjz20nHge8PfIstOLo1RpnL28Xg20o4GUMIGRMA4GA1Ud
DwEB/wQEAwICvDAnBgNVHSUEIDAeBggrBgEFBQcDAQYIKwYBBQUHAwIGCCsGAQUF
BwMBMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFMMsIrLOGJwbXe3g2MP6D+6e
x1+1MCYGA1UdEQQfMB2CCGljYmMuY29tgQtjYUBpY2JjLmNvbYcErBAYUTAKBggq
hkjOPQQDAgNHADBEAiB7JjMol4qJPizD8kgEocRyvOMcxJbaFpwLjZsHtxpAtgIg
SvIdFT+xUSTZ8mLB0AYcxIl6stTLN64XvCuHRUaYq9s=
-----END CERTIFICATE-----
	`
	ABClientCrt = `
-----BEGIN CERTIFICATE-----
MIICWzCCAgGgAwIBAgIQJa0atDW0R2vD7KG17PdhwTAKBggqhkjOPQQDAjBgMQsw
CQYDVQQGEwJDTjELMAkGA1UECBMCWkoxCzAJBgNVBAcTAkhaMQ8wDQYDVQQKEwZ5
bC5jb20xDzANBgNVBAsTBnlsLmNvbTEVMBMGA1UEAwwM5Lit5Zu96ZO26IGUMCAX
DTI1MDEwNzAzNDM0MloYDzI1MjUwMTA3MDM0MzQyWjBpMQswCQYDVQQGEwJDTjEL
MAkGA1UECBMCWkoxCzAJBgNVBAcTAkhaMQ8wDQYDVQQKEwZ5bC5jb20xDzANBgNV
BAsTBnlsLmNvbTEeMBwGA1UEAwwV5Lit5Zu96ZO26IGU5a6i5oi356uvMFkwEwYH
KoZIzj0CAQYIKoZIzj0DAQcDQgAE7HRdnKQTFaPUQzN1iJDCmURKHHpAYPMljpoH
MgwBVJ0ZQVBKKEUm2FGYBwt6a/dc3BEIE7UfT9KJ31mfzuZbvKOBkTCBjjAOBgNV
HQ8BAf8EBAMCArwwJwYDVR0lBCAwHgYIKwYBBQUHAwEGCCsGAQUFBwMCBggrBgEF
BQcDATAMBgNVHRMBAf8EAjAAMB8GA1UdIwQYMBaAFNXjJRuDhTjfu/bZmvFpQenB
WbfuMCQGA1UdEQQdMBuCBnlsLmNvbYELY2FAaWNiYy5jb22HBKwQGFEwCgYIKoZI
zj0EAwIDSAAwRQIhAOooxjnapqr82Muu5CPzC5E496HCMYlWPVRDLcbIzWilAiAB
Q/pBCQEwLLg9LJjw6ESilHf96TNkTJSCBJ+CI4TfoQ==
-----END CERTIFICATE-----
	`
	ABClientKey = `
-----BEGIN ECDSA PRIVATE KEY-----
MHcCAQEEIMIYoXzLQQwJNXIS7fyCjJtRaSY383oLmgjdmLrtgxEroAoGCCqGSM49
AwEHoUQDQgAE7HRdnKQTFaPUQzN1iJDCmURKHHpAYPMljpoHMgwBVJ0ZQVBKKEUm
2FGYBwt6a/dc3BEIE7UfT9KJ31mfzuZbvA==
-----END ECDSA PRIVATE KEY-----
	`
)
