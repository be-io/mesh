/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package redis

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/be-io/mesh/client/golang/tool"
	"github.com/redis/go-redis/v9"
	"github.com/tjfoc/gmsm/sm4"
	"net/url"
	"runtime"
	"strings"
	"time"
)

func getEncDecode(ctx context.Context, input string) string {
	if strings.HasPrefix(input, "ENC(") && strings.HasSuffix(input, ")") {
		key := "C+0CF9Uj7oLNKmq69k0tAA=="
		keyAsBytes, _ := base64.StdEncoding.DecodeString(key)
		encAsBytes, _ := hex.DecodeString(input[4 : len(input)-1])
		plainBytes, err := sm4.Sm4Ecb(keyAsBytes, encAsBytes, false)
		if err != nil {
			log.Error(ctx, "getEncDecode error : %v", err)
			return input
		}
		return string(plainBytes)
	}
	return input
}

func (that *redisAccessLayer) NewClient(ctx context.Context, servers string) (redis.UniversalClient, error) {
	//servers = "redis://username:ENC(835a7ba1495475bd403667a1e699ec5c)@127.0.0.1:6379"
	uri, err := url.Parse(servers)
	if nil != err {
		return nil, cause.Errorf("Redis proxy dont startup because servers %s is invalid, %s. ", servers, err.Error())
	}
	userinfo := tool.Anyone(uri.User, &url.Userinfo{})
	password, _ := userinfo.Password()
	password = getEncDecode(ctx, password)
	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:                 strings.Split(uri.Host, ","),
		DB:                    0,
		Username:              userinfo.Username(),
		Password:              password,
		SentinelUsername:      "",
		SentinelPassword:      "",
		MaxRetries:            3,
		MinRetryBackoff:       time.Millisecond * 100,
		MaxRetryBackoff:       time.Millisecond * 300,
		DialTimeout:           time.Second * 10,
		ReadTimeout:           time.Second * 60,
		WriteTimeout:          time.Second * 60,
		ContextTimeoutEnabled: false,
		PoolFIFO:              false,
		PoolSize:              runtime.NumCPU(),
		PoolTimeout:           time.Second * 12,
		MinIdleConns:          1,
		MaxIdleConns:          3,
		ConnMaxIdleTime:       time.Minute * 2,
		ConnMaxLifetime:       time.Minute * 25,
		MaxRedirects:          3,
		ReadOnly:              false,
		RouteByLatency:        true,
		RouteRandomly:         false,
		MasterName:            "",
	}), nil
}
