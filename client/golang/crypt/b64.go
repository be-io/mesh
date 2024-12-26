/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package crypt

import (
	"bytes"
	"encoding/base64"
	"github.com/opendatav/mesh/client/golang/cause"
	"io/ioutil"
)

var Base64 = new(b64)

type b64 struct {
}

func (that *b64) FormKey(bytes []byte, headers ...string) (string, error) {
	return that.FormPlain(bytes, headers...)
}

func (that *b64) InformKey(data string) ([]byte, error) {
	return that.InformPlain(data)
}

func (that *b64) FormSIG(sig []byte, headers ...string) (string, error) {
	return that.FormPlain(sig, headers...)
}

func (that *b64) InformSIG(data string) ([]byte, error) {
	return that.InformPlain(data)
}

func (that *b64) FormCipher(data []byte, headers ...string) (string, error) {
	return that.FormPlain(data, headers...)
}

func (that *b64) InformCipher(data string) ([]byte, error) {
	return that.InformPlain(data)
}

func (that *b64) FormPlain(data []byte, headers ...string) (string, error) {
	var buffer bytes.Buffer
	writer := base64.NewEncoder(base64.StdEncoding, &buffer)
	if _, err := writer.Write(data); nil != err {
		return "", cause.Error(err)
	}
	if err := writer.Close(); nil != err {
		return "", cause.Error(err)
	}
	return buffer.String(), nil
}

func (that *b64) InformPlain(data string) ([]byte, error) {
	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(data))
	return ioutil.ReadAll(reader)
}
