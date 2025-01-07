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
	"encoding/pem"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/mpc"
	"github.com/opendatav/mesh/ptp"
	"io"
	"net/http"
)

func main() {
	ctx := mpc.Context()
	if err := do(ctx); nil != err {
		log.Error(ctx, "%s", err)
	}
}

func do(ctx context.Context) error {
	v, _ := pem.Decode([]byte(ARootCrt))
	x, err := x509.ParseCertificate(v.Bytes)
	if nil != err {
		return cause.Error(err)
	}
	log.Info(ctx, "%s", x.PermittedDNSDomains)
	keys, err := tls.X509KeyPair([]byte(BAClientCrt), []byte(BAClientKey))
	if nil != err {
		return cause.Error(err)
	}
	ca := x509.NewCertPool()
	ca.AppendCertsFromPEM([]byte(ARootCrt))
	cc := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
		RootCAs:            ca,
		Certificates:       []tls.Certificate{keys},
	}}}
	payload := ptp.PushInbound{
		Topic:    "",
		Payload:  []byte(""),
		Metadata: map[string]string{},
	}
	buf, err := ptp.Encode(payload, "application/json")
	if nil != err {
		return cause.Error(err)
	}
	r, err := http.NewRequest(http.MethodPost, "https://yl.com:7304/v1/interconn/chan/invoke", bytes.NewBuffer(buf))
	if nil != err {
		return cause.Error(err)
	}
	r.Header.Set("x-ptp-target-node-id", "IC075")
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
)
