/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package iset

import (
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/tidwall/redcon"
	"net"
)

var _ Conn = new(RConn)

func WithConn(conn redcon.Conn) Conn {
	return &RConn{conn: conn}
}

// Conn
// For Simple Strings, the first byte of the reply is "+"
// For Errors, the first byte of the reply is "-"
// For Integers, the first byte of the reply is ":"
// For Bulk Strings, the first byte of the reply is "$"
// For Arrays, the first byte of the reply is "*"
type Conn interface {
	redcon.Conn
	WriteErr(err error)
	WriteWriteBulkStrings(rs []string)
	WriteAnyArray(rs []interface{})
}

type RConn struct {
	conn redcon.Conn
	serv func(rc redcon.Conn, cmd redcon.Command)
}

func (that *RConn) WriteAnyArray(rs []interface{}) {
	if len(rs) < 1 {
		that.WriteArray(0)
		return
	}
	that.WriteArray(len(rs))
	for _, r := range rs {
		that.WriteAny(r)
	}
}

func (that *RConn) WriteWriteBulkStrings(rs []string) {
	if len(rs) < 1 {
		that.WriteArray(0)
		return
	}
	that.WriteArray(len(rs))
	for _, r := range rs {
		that.WriteBulkString(r)
	}
}

func (that *RConn) WriteErr(err error) {
	if IsNil(err) {
		that.conn.WriteNull()
	} else {
		that.conn.WriteError(err.Error())
	}
}

func (that *RConn) RemoteAddr() string {
	return that.conn.RemoteAddr()
}

func (that *RConn) Close() error {
	return that.conn.Close()
}

func (that *RConn) WriteError(msg string) {
	that.conn.WriteError(msg)
}

func (that *RConn) WriteString(str string) {
	that.conn.WriteString(str)
}

func (that *RConn) WriteBulk(bulk []byte) {
	that.conn.WriteBulk(bulk)
}

func (that *RConn) WriteBulkString(bulk string) {
	that.conn.WriteBulkString(bulk)
}

func (that *RConn) WriteInt(num int) {
	that.conn.WriteInt(num)
}

func (that *RConn) WriteInt64(num int64) {
	that.conn.WriteInt64(num)
}

func (that *RConn) WriteUint64(num uint64) {
	that.conn.WriteUint64(num)
}

func (that *RConn) WriteArray(count int) {
	that.conn.WriteArray(count)
}

func (that *RConn) WriteNull() {
	that.conn.WriteNull()
}

func (that *RConn) WriteRaw(data []byte) {
	that.conn.WriteRaw(data)
}

func (that *RConn) WriteAny(any interface{}) {
	that.conn.WriteAny(any)
}

func (that *RConn) Context() interface{} {
	return that.conn.Context()
}

func (that *RConn) SetContext(v interface{}) {
	that.conn.SetContext(v)
}

func (that *RConn) SetReadBuffer(bytes int) {
	that.conn.SetReadBuffer(bytes)
}

func (that *RConn) Detach() redcon.DetachedConn {
	return that.conn.Detach()
}

func (that *RConn) ReadPipeline() []redcon.Command {
	return that.conn.ReadPipeline()
}

func (that *RConn) PeekPipeline() []redcon.Command {
	return that.conn.PeekPipeline()
}

func (that *RConn) NetConn() net.Conn {
	return that.conn.NetConn()
}

var _ Conn = new(NoWriteConn)

type NoWriteConn struct {
	conn Conn
}

func (that *NoWriteConn) RemoteAddr() string {
	return that.conn.RemoteAddr()
}

func (that *NoWriteConn) Close() error {
	return that.conn.Close()
}

func (that *NoWriteConn) WriteError(msg string) {
	log.Warn0("Redis command exec, %s", msg)
}

func (that *NoWriteConn) WriteString(str string) {
}

func (that *NoWriteConn) WriteBulk(bulk []byte) {
}

func (that *NoWriteConn) WriteBulkString(bulk string) {
}

func (that *NoWriteConn) WriteInt(num int) {
}

func (that *NoWriteConn) WriteInt64(num int64) {
}

func (that *NoWriteConn) WriteUint64(num uint64) {
}

func (that *NoWriteConn) WriteArray(count int) {
}

func (that *NoWriteConn) WriteNull() {
}

func (that *NoWriteConn) WriteRaw(data []byte) {
}

func (that *NoWriteConn) WriteAny(any interface{}) {
}

func (that *NoWriteConn) Context() interface{} {
	return that.conn.Context()
}

func (that *NoWriteConn) SetContext(v interface{}) {
}

func (that *NoWriteConn) SetReadBuffer(bytes int) {
}

func (that *NoWriteConn) Detach() redcon.DetachedConn {
	return that.conn.Detach()
}

func (that *NoWriteConn) ReadPipeline() []redcon.Command {
	return that.conn.ReadPipeline()
}

func (that *NoWriteConn) PeekPipeline() []redcon.Command {
	return that.conn.PeekPipeline()
}

func (that *NoWriteConn) NetConn() net.Conn {
	return that.conn.NetConn()
}

func (that *NoWriteConn) WriteErr(err error) {
	if IsNil(err) {
		that.WriteNull()
		return
	}
	that.WriteError(err.Error())
}

func (that *NoWriteConn) WriteWriteBulkStrings(rs []string) {
}

func (that *NoWriteConn) WriteAnyArray(rs []interface{}) {
}
