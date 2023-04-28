/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package crypt

import (
	"strings"
	"testing"
)

func TestB64(t *testing.T) {
	explain, err := Base64.FormPlain([]byte(strings.Repeat("1234567890", 1000)))
	if nil != err {
		t.Error(err)
		return
	}
	bytes, err := Base64.InformPlain(explain)
	if nil != err {
		t.Error(err)
		return
	}
	if strings.Repeat("1234567890", 1000) != string(bytes) {
		t.Error("Base64 codec failed.")
	}
}

func TestGz(t *testing.T) {
	explain, err := GZIP.FormPlain([]byte(strings.Repeat("1234567890", 1000)))
	if nil != err {
		t.Error(err)
		return
	}
	bytes, err := GZIP.InformPlain(explain)
	if nil != err {
		t.Error(err)
		return
	}
	if strings.Repeat("1234567890", 1000) != string(bytes) {
		t.Error("GZIP codec failed.")
	}
}
