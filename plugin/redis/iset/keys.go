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
	Register(new(Keys))
}

type Keys struct {
}

func (that *Keys) Length() int {
	return 1
}

func (that *Keys) Name() string {
	return "KEYS"
}

// Serve
// Returns all keys matching pattern.
//
// While the time complexity for this operation is O(N), the constant times are fairly low. For example, Redis running on an entry level laptop can scan a 1 million key database in 40 milliseconds.
//
// Warning: consider KEYS as a command that should only be used in production environments with extreme care. It may ruin performance when it is executed against large databases. This command is intended for debugging and special operations, such as changing your keyspace layout. Don't use KEYS in your regular application code. If you're looking for a way to find keys in a subset of your keyspace, consider using SCAN or sets.
//
// Supported glob-style patterns:
//
// h?llo matches hello, hallo and hxllo
// h*llo matches hllo and heeeello
// h[ae]llo matches hello and hallo, but not hillo
// h[^e]llo matches hallo, hbllo, ... but not hello
// h[a-b]llo matches hallo and hbllo
// Use \ to escape special characters if you want to match them verbatim.
// Array reply: list of keys matching pattern.
func (that *Keys) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	switch ref := client.(type) {
	case *redis.ClusterClient:
		that.ClusterKeys(ctx, conn, cmd, ref)
	default:
		that.MasterKeys(ctx, conn, cmd, client)
	}
}

func (that *Keys) ClusterKeys(ctx context.Context, conn Conn, cmd redcon.Command, client *redis.ClusterClient) {
	var keys []string
	err := client.ForEachMaster(ctx, func(ctx context.Context, ref *redis.Client) error {
		iter := client.Scan(ctx, 0, that.KeysPattern(cmd), 10000).Iterator()
		for iter.Next(ctx) {
			keys = append(keys, iter.Val())
		}
		if err := iter.Err(); nil != err {
			return err
		}
		return nil
	})
	if nil != err {
		conn.WriteErr(err)
		return
	}
	conn.WriteWriteBulkStrings(keys)
}

func (that *Keys) MasterKeys(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	var keys []string
	iter := client.Scan(ctx, 0, that.KeysPattern(cmd), 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); nil != err {
		conn.WriteErr(err)
		return
	}
	conn.WriteWriteBulkStrings(keys)
}

func (that *Keys) KeysPattern(cmd redcon.Command) string {
	if len(cmd.Args) > 1 {
		return string(cmd.Args[1])
	}
	return ""
}
