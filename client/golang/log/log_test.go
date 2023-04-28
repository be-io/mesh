/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package log

import (
	"fmt"
	"github.com/be-io/mesh/client/golang/macro"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestLevelDisable(t *testing.T) {
	ctx := macro.Context()
	SetLevel(FATAL)
	t.Log("fatal test")
	Fatal(ctx, "fatal")
	Error(ctx, "error")
	Warn(ctx, "warn")
	Info(ctx, "info")
	Debug(ctx, "debug")
	Stack(ctx, "stack")

	SetLevel(INFO)
	t.Log("info test")
	Fatal(ctx, "fatal")
	Error(ctx, "error")
	Warn(ctx, "warn")
	Info(ctx, "info")
	Debug(ctx, "debug")
	Stack(ctx, "stack")

	SetLevel(ERROR)
	t.Log("error test")
	Fatal(ctx, "fatal")
	Error(ctx, "error")
	Warn(ctx, "warn")
	Info(ctx, "info")
	Debug(ctx, "debug")
	Stack(ctx, "stack")

	SetLevel(WARN)
	t.Log("warn test")
	Fatal(ctx, "fatal")
	Error(ctx, "error")
	Warn(ctx, "warn")
	Info(ctx, "info")
	Debug(ctx, "debug")
	Stack(ctx, "stack")

	SetLevel(INFO)
	t.Log("info test")
	Fatal(ctx, "fatal")
	Error(ctx, "error")
	Warn(ctx, "warn")
	Info(ctx, "info")
	Debug(ctx, "debug")
	Stack(ctx, "stack")

	SetLevel(DEBUG)
	t.Log("debug test")
	Fatal(ctx, "fatal")
	Error(ctx, "error")
	Warn(ctx, "warn")
	Info(ctx, "info")
	Debug(ctx, "debug")
	Stack(ctx, "stack")

	SetLevel(STACK)
	t.Log("stack test")
	Fatal(ctx, "fatal")
	Error(ctx, "error")
	Warn(ctx, "warn")
	Info(ctx, "info")
	Debug(ctx, "debug")
	Stack(ctx, "stack")
}

func TestURLParse(t *testing.T) {
	uri, err := url.Parse("file:///var/log/be/mesh.log?size=100&backups=120&age=28&compress=1")
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(uri)
	home, _ := os.UserHomeDir()
	t.Log(fmt.Sprintf("file://%s", filepath.Join(home, "mesh", "xxx.log?size=100&backups=120&age=28&compress=1")))
}
