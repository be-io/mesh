/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package iset

import (
	"context"
	"github.com/be-io/mesh/client/golang/log"
	"github.com/redis/go-redis/v9"
	"github.com/tidwall/redcon"
)

func init() {
	Register(new(Quit))
}

type Quit struct {
}

func (that *Quit) Length() int {
	return 1
}

func (that *Quit) Name() string {
	return "QUIT"
}

func (that *Quit) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	conn.WriteString("OK")
	log.Catch(conn.Close())
}
