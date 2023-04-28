/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package schema

import (
	"fmt"
	"strings"
	"testing"
)

func TestExprTree(t *testing.T) {
	tree := GenericTree("List<List<Map<Map<String,Integer>,String,Map<String, String>,Map<String,Map<Object,String>,String,String,Map<String,String>>>>>")
	x := FormatTree(tree, func(root string, children []string) string {
		if len(children) < 1 {
			return root
		}
		return fmt.Sprintf("[%s, %s]", root, strings.Join(children, ", "))
	})
	t.Log(x)
}
