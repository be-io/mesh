/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package crypt

import (
	"bytes"
	"compress/gzip"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"io/ioutil"
)

var GZIP = &gz{Format: Base64}
var GZP = &gz{Format: PEM}

type gz struct {
	Format Format
}

func (that *gz) FormKey(bytes []byte, headers ...string) (string, error) {
	return that.Compress(bytes, func(data []byte) (string, error) { return that.Format.FormKey(data, headers...) })
}

func (that *gz) InformKey(key string) ([]byte, error) {
	return that.Decompress(key, func(data string) ([]byte, error) { return that.Format.InformKey(data) })
}

func (that *gz) FormSIG(bytes []byte, headers ...string) (string, error) {
	return that.Compress(bytes, func(data []byte) (string, error) { return that.Format.FormSIG(data, headers...) })
}

func (that *gz) InformSIG(bytes string) ([]byte, error) {
	return that.Decompress(bytes, func(data string) ([]byte, error) { return that.Format.InformSIG(data) })
}

func (that *gz) FormCipher(bytes []byte, headers ...string) (string, error) {
	return that.Compress(bytes, func(data []byte) (string, error) { return that.Format.FormCipher(data, headers...) })
}

func (that *gz) InformCipher(bytes string) ([]byte, error) {
	return that.Decompress(bytes, func(data string) ([]byte, error) { return that.Format.InformCipher(data) })
}

func (that *gz) FormPlain(bytes []byte, headers ...string) (string, error) {
	return that.Compress(bytes, func(data []byte) (string, error) { return that.Format.FormPlain(data, headers...) })
}

func (that *gz) InformPlain(bytes string) ([]byte, error) {
	return that.Decompress(bytes, func(data string) ([]byte, error) { return that.Format.InformPlain(data) })
}

func (that *gz) Compress(data []byte, format func(data []byte) (string, error)) (string, error) {
	var buff bytes.Buffer
	compressor := gzip.NewWriter(&buff)
	if _, err := compressor.Write(data); nil != err {
		return "", cause.Error(err)
	}
	if err := compressor.Close(); nil != err {
		return "", cause.Error(err)
	}
	return format(buff.Bytes())
}

func (that *gz) Decompress(data string, informat func(data string) ([]byte, error)) ([]byte, error) {
	cata, err := informat(data)
	if nil != err {
		return nil, cause.Error(err)
	}
	decompressor, err := gzip.NewReader(bytes.NewBuffer(cata))
	if nil != err {
		return nil, cause.Error(err)
	}
	defer func() { log.Catch(decompressor.Close()) }()
	return ioutil.ReadAll(decompressor)
}
