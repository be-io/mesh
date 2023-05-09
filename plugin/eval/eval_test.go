/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package eval

import (
	"fmt"
	"github.com/be-io/mesh/client/golang/macro"
	"github.com/be-io/mesh/client/golang/mpc"
	"github.com/be-io/mesh/client/golang/plugin"
	"github.com/be-io/mesh/client/golang/prsim"
	"github.com/be-io/mesh/client/golang/types"
	_ "github.com/be-io/mesh/plugin/metabase"
	"testing"
)

func init() {
	macro.Mock(prsim.ISequence, macro.MeshSPI, metabase.Name)
	macro.Mock(prsim.IKV, macro.MeshSPI, metabase.Name)
}

func TestCompile(t *testing.T) {
	cases := []*types.Script{
		{
			Code:       "boolean",
			Name:       "布尔规则",
			Desc:       "布尔规则",
			Kind:       types.VALUE,
			Expr:       "true",
			Attachment: map[string]string{},
		},
		{
			Code:       "int",
			Name:       "整形规则",
			Desc:       "整形规则",
			Kind:       types.VALUE,
			Expr:       "1",
			Attachment: map[string]string{},
		},
		{
			Code:       "string",
			Name:       "字符规则",
			Desc:       "字符规则",
			Kind:       types.VALUE,
			Expr:       "abc",
			Attachment: map[string]string{},
		},
		{
			Code:       "expression",
			Name:       "表达式规则",
			Desc:       "表达式规则",
			Kind:       types.EXPRESSION,
			Expr:       "fmt.Sprintf(\"Name:%s\", args[\"name\"])",
			Attachment: map[string]string{},
		},
		{
			Code:       "expression",
			Name:       "表达式规则更新",
			Desc:       "表达式规则",
			Kind:       types.EXPRESSION,
			Expr:       "fmt.Sprintf(\"NewName:%s\", args[\"name\"])",
			Attachment: map[string]string{},
		},
		{
			Code: "script",
			Name: "脚本规则",
			Desc: "脚本规则",
			Kind: types.SCRIPT,
			Expr: `
			import (
				"fmt"
				"context"
			)
			func Invoke(ctx context.Context, args map[string]string, dft string) (string, error) {
				return fmt.Sprintf("%d", 1), nil
			}
			`,
			Attachment: map[string]string{},
		},
	}
	ctx := mpc.Context()
	dsn := "root:@tcp(127.0.0.1:3306)/mesh"
	container := plugin.LoadC("omega")
	container.Start(ctx, fmt.Sprintf("--dsn=%s", dsn))
	defer container.Stop(ctx)
	for _, c := range cases {
		code, err := engine.Compile(ctx, c)
		if nil != err {
			t.Error(err)
		} else {
			t.Log(fmt.Sprintf("Compile %s(%s) -c %s", c.Name, c.Code, code))
		}
		r, err := engine.Exec(ctx, c.Code, map[string]string{}, "")
		if nil != err {
			t.Error(err)
		} else {
			t.Log(fmt.Sprintf("Invoke %s(%s) -o %s", c.Name, c.Code, r))
		}
	}
	scripts, err := engine.Dump(ctx, map[string]string{})
	if nil != err {
		t.Error(err)
	} else {
		t.Log(scripts)
	}
	page, err := engine.Index(ctx, new(types.Paging).Reset("", 0, 10))
	if nil != err {
		t.Error(err)
	} else {
		t.Log(page)
	}
}
