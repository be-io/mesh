/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package iset

import (
	"context"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/redis/go-redis/v9"
	"github.com/tidwall/redcon"
	"net"
	"strings"
)

func init() {
	var _ redcon.Conn = new(PipConn)
	Register(new(Multi))
}

type Multi struct {
}

func (that *Multi) Length() int {
	return 1
}

func (that *Multi) Name() string {
	return "MULTI"
}

// Serve
// Marks the start of a transaction block. Subsequent commands will be queued for atomic execution using EXEC.
//
// Return
// Simple string reply: always OK.
func (that *Multi) Serve(ctx context.Context, conn Conn, cmd redcon.Command, client redis.Cmdable) {
	rc, ok := client.(redis.UniversalClient)
	if !ok {
		conn.WriteError("Unknown connection state.")
		return
	}
	pipeline := client.TxPipeline()
	pip := &Pip{conn: conn.Detach(), client: rc, pipe: pipeline}
	conn.WriteString("OK")
	go pip.Begin(ctx)
}

type Pip struct {
	conn   redcon.DetachedConn
	client redis.UniversalClient
	pipe   redis.Pipeliner
}

func (that *Pip) End(ctx context.Context) {
	if err := that.conn.Close(); nil != err {
		log.Error(ctx, err.Error())
	}
}

func (that *Pip) Begin(ctx context.Context) {
	if err := log.PError(ctx, func() error {
		pc := &PipConn{conn: that.conn, ctx: ctx}
		for {
			command, err := that.conn.ReadCommand()
			if nil != err {
				that.conn.WriteError(err.Error())
				return err
			}
			switch strings.ToUpper(string(command.Args[0])) {
			case "EXEC":
				Exec(ctx, that.conn, command, that.pipe)
				if err = that.conn.Flush(); nil != err {
					log.Error(ctx, "Redis proxy flush conn, %s", err.Error())
					return nil
				}
			case "DISCARD":
				that.conn.WriteString("OK")
				return nil
			case "WATCH":
				that.conn.WriteString("OK")
				return nil
			default:
				Exec(ctx, pc, command, that.pipe)
				that.conn.WriteString("QUEUED")
			}
		}
	}); nil != err {
		log.Error(ctx, err.Error())
	}
}

type PipConn struct {
	conn redcon.Conn
	ctx  context.Context
}

func (that *PipConn) RemoteAddr() string {
	return that.conn.RemoteAddr()
}

func (that *PipConn) Close() error {
	return that.conn.Close()
}

func (that *PipConn) WriteError(msg string) {
	log.Error(that.ctx, msg)
}

func (that *PipConn) WriteString(str string) {

}

func (that *PipConn) WriteBulk(bulk []byte) {

}

func (that *PipConn) WriteBulkString(bulk string) {

}

func (that *PipConn) WriteInt(num int) {

}

func (that *PipConn) WriteInt64(num int64) {

}

func (that *PipConn) WriteUint64(num uint64) {

}

func (that *PipConn) WriteArray(count int) {

}

func (that *PipConn) WriteNull() {

}

func (that *PipConn) WriteRaw(data []byte) {

}

func (that *PipConn) WriteAny(any interface{}) {

}

func (that *PipConn) Context() interface{} {
	return that.conn.Context()
}

func (that *PipConn) SetContext(v interface{}) {
	that.conn.SetContext(v)
}

func (that *PipConn) SetReadBuffer(bytes int) {
	that.conn.SetReadBuffer(bytes)
}

func (that *PipConn) Detach() redcon.DetachedConn {
	return that.conn.Detach()
}

func (that *PipConn) ReadPipeline() []redcon.Command {
	return that.conn.ReadPipeline()
}

func (that *PipConn) PeekPipeline() []redcon.Command {
	return that.conn.PeekPipeline()
}

func (that *PipConn) NetConn() net.Conn {
	return that.conn.NetConn()
}
