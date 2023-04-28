/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package log

import (
	"context"
	"flag"
	"fmt"
	"runtime"
	"strings"
)

var store = map[string]Logger{}

func Provide(log Logger) {
	store[log.Name()] = log
}

func Load(names ...string) Logger {
	if len(names) > 0 && nil != store[names[0]] {
		return store[names[0]]
	}
	return zer
}

func KV(key string, value interface{}) *Pair {
	return &Pair{Key: key, Value: value}
}

func Split(args []interface{}) ([]interface{}, []*Pair) {
	var fields []*Pair
	var params []interface{}
	for _, arg := range args {
		if param, ok := arg.(Pair); ok {
			fields = append(fields, KV(param.Key, param.Value))
		} else {
			params = append(params, arg)
		}
	}
	return params, fields
}

func Split0(ctx context.Context, args []interface{}) ([]interface{}, []*Pair) {
	params, fields := Split(args)
	variables := ctx.Value("vars")
	if nil == variables {
		return params, fields
	}
	kv, ok := variables.(map[string]interface{})
	if !ok {
		return params, fields
	}
	for key, value := range kv {
		fields = append(fields, KV(key, value))
	}
	return params, fields
}

func Caller(skip int) string {
	_, name, line, _ := runtime.Caller(skip)
	if IsTest() {
		return fmt.Sprintf("%s:%d", name, line)
	}
	return fmt.Sprintf("%s:%d", name[strings.LastIndex(name, "/")+1:], line)
}

// IsTest is golang testing.
func IsTest() bool {
	return nil != flag.Lookup("test.v")
}
