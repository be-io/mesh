/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package iset

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/tidwall/redcon"
)

func init() {
	Register(new(Auth))
}

type Auth struct {
}

func (that *Auth) Length() int {
	return 1
}

func (that *Auth) Name() string {
	return "AUTH"
}

// Serve
// Simple string reply or an error if the password, or username/password pair, is invalid.
func (that *Auth) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	conn.WriteString("OK")
}
