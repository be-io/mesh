/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package crypt

import (
	"bytes"
	"encoding/pem"
	"github.com/opendatav/mesh/client/golang/cause"
)

var PEM = &Pem{Format: Base64}

type Pem struct {
	Format Format
}

func (that *Pem) FormKey(data []byte, headers ...string) (string, error) {
	kind := "CERTIFICATE"
	if len(headers) > 0 {
		kind = headers[0]
	}
	var buffer bytes.Buffer
	if err := pem.Encode(&buffer, &pem.Block{Type: kind, Bytes: data}); nil != err {
		return "", cause.Error(err)
	}
	return buffer.String(), nil
}

func (that *Pem) InformKey(key string) ([]byte, error) {
	block, _ := pem.Decode([]byte(key))
	if nil == block {
		return nil, cause.Errorf("Decode certificate with unexpected result")
	}
	return block.Bytes, nil
}

func (that *Pem) FormSIG(sig []byte, headers ...string) (string, error) {
	return that.Format.FormSIG(sig, headers...)
}

func (that *Pem) InformSIG(sig string) ([]byte, error) {
	return that.Format.InformSIG(sig)
}

func (that *Pem) FormCipher(data []byte, headers ...string) (string, error) {
	return that.Format.FormCipher(data, headers...)
}

func (that *Pem) InformCipher(data string) ([]byte, error) {
	return that.Format.InformCipher(data)
}

func (that *Pem) FormPlain(data []byte, headers ...string) (string, error) {
	return that.Format.FormPlain(data, headers...)
}

func (that *Pem) InformPlain(data string) ([]byte, error) {
	return that.Format.InformPlain(data)
}
