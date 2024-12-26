/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package redis

import (
	"github.com/opendatav/mesh/client/golang/mpc"
	"net/url"
	"testing"
)

func TestParseServers(t *testing.T) {
	uri, err := url.Parse("redis://user:password@127.0.0.1:3306,127.0.0.1:3307,127.0.0.1:3308")
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(uri.Scheme)
	t.Log(uri.User.Username())
	t.Log(uri.User.Password())
	t.Log(uri.Host)
}

func TestRedisClusterGet(t *testing.T) {
	ctx := mpc.Context()
	client, err := new(redisAccessLayer).NewClient(ctx, &redisOption{})
	if nil != err {
		t.Error(err)
		return
	}
	cmd := client.Get(ctx, "")
	if err := cmd.Err(); nil != err {
		t.Error(err)
		return
	}
	t.Log(cmd.Val())
}
