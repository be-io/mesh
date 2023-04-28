/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package kms

import (
	"fmt"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/types"
	"testing"
)

func TestApplyRoot(t *testing.T) {
	keys, err := new(localKMS).ApplyRoot(macro.Context(), &types.KeyCsr{
		Domain: "trustbe.cn",
		IsCA:   true,
		Length: 4096,
	})
	if nil != err {
		t.Error(err)
		return
	}
	for _, key := range keys {
		fmt.Println(key.Kind)
		fmt.Println(key.Key)
	}
}
