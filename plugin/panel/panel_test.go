/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package panel

import "testing"

func TestStaticFiles(t *testing.T) {
	fs, err := home.ReadDir("static")
	if nil != err {
		t.Error(err)
		return
	}
	for _, f := range fs {
		t.Log(f.Name())
	}
}
