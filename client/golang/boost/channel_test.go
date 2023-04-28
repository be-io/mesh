/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package boost

import (
	"context"
	"github.com/be-io/mesh/client/golang/mpc"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"testing"
	"time"
)

var (
	InitialCap = 5
	MaxIdleCap = 10
	MaximumCap = 100
	network    = "tcp"
	address    = "127.0.0.1:7777"
	//factory    = func() (interface{}, error) { return net.Dial(network, address) }
	factory = func(ctx context.Context) (Connection, error) {
		return rpc.DialHTTP("tcp", address)
	}
	closeFac = func(v Connection) error {
		nc := v.(*rpc.Client)
		return nc.Close()
	}
)

func init() {
	// used for factory function
	go rpcServer()
	time.Sleep(time.Millisecond * 300) // wait until tcp server has been settled

	rand.Seed(time.Now().UTC().UnixNano())
}

func TestNew(t *testing.T) {
	ctx := mpc.Context()
	p, err := newChannelPool()
	defer p.Release(ctx)
	if nil != err {
		t.Errorf("New error: %s", err)
	}
}
func TestPool_Get_Impl(t *testing.T) {
	ctx := mpc.Context()
	p, _ := newChannelPool()
	defer p.Release(ctx)

	conn, err := p.Borrow(ctx)
	if nil != err {
		t.Errorf("Get error: %s", err)
	}
	_, ok := conn.(*rpc.Client)
	if !ok {
		t.Errorf("Conn is not of type poolConn")
	}
	p.Return(ctx, conn)
}

func TestPool_Get(t *testing.T) {
	ctx := mpc.Context()
	p, _ := newChannelPool()
	defer p.Release(ctx)

	_, err := p.Borrow(ctx)
	if nil != err {
		t.Errorf("Get error: %s", err)
	}

	// after one get, current capacity should be lowered by one.
	if p.Len() != (InitialCap - 1) {
		t.Errorf("Get error. Expecting %d, got %d",
			(InitialCap - 1), p.Len())
	}

	// get them all
	var wg sync.WaitGroup
	for i := 0; i < (MaximumCap - 1); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.Borrow(mpc.Context())
			if nil != err {
				t.Errorf("Get error: %s", err)
			}
		}()
	}
	wg.Wait()

	if p.Len() != 0 {
		t.Errorf("Get error. Expecting %d, got %d",
			(InitialCap - 1), p.Len())
	}

	_, err = p.Borrow(ctx)
	if err != ErrPoolMaxActiveConnReached {
		t.Errorf("Get error: %s", err)
	}

}

func TestPool_Put(t *testing.T) {
	ctx := mpc.Context()
	pconf := PoolConfig{InitialCap: InitialCap, MaxCap: MaximumCap, Factory: factory, IdleTimeout: time.Second * 20,
		MaxIdle: MaxIdleCap}
	p, err := NewChannelPool(ctx, &pconf)
	if nil != err {
		t.Fatal(err)
	}
	defer p.Release(ctx)

	// get/create from the pool
	conns := make([]Connection, MaximumCap)
	for i := 0; i < MaximumCap; i++ {
		conn, _ := p.Borrow(ctx)
		conns[i] = conn
	}

	// now put them all back
	for _, conn := range conns {
		p.Return(ctx, conn)
	}

	if p.Len() != MaxIdleCap {
		t.Errorf("Put error len. Expecting %d, got %d",
			1, p.Len())
	}

	p.Release(ctx) // close pool

}

func TestPool_UsedCapacity(t *testing.T) {
	ctx := mpc.Context()
	p, _ := newChannelPool()
	defer p.Release(ctx)

	if p.Len() != InitialCap {
		t.Errorf("InitialCap error. Expecting %d, got %d",
			InitialCap, p.Len())
	}
}

func TestPool_Close(t *testing.T) {
	ctx := mpc.Context()
	p, _ := newChannelPool()

	// now close it and test all cases we are expecting.
	p.Release(ctx)

	c := p.(*channelPool)

	if c.conns != nil {
		t.Errorf("Close error, conns channel should be nil")
	}

	if c.factory != nil {
		t.Errorf("Close error, factory should be nil")
	}

	_, err := p.Borrow(ctx)
	if err == nil {
		t.Errorf("Close error, get conn should return an error")
	}

	if p.Len() != 0 {
		t.Errorf("Close error used capacity. Expecting 0, got %d", p.Len())
	}
}

func TestPoolConcurrent(t *testing.T) {
	ctx := mpc.Context()
	p, _ := newChannelPool()
	pipe := make(chan Connection, 0)

	go func() {
		p.Release(ctx)
	}()

	for i := 0; i < MaximumCap; i++ {
		go func() {
			conn, _ := p.Borrow(ctx)

			pipe <- conn
		}()

		go func() {
			conn := <-pipe
			if conn == nil {
				return
			}
			p.Return(ctx, conn)
		}()
	}
}

func TestPoolWriteRead(t *testing.T) {
	ctx := mpc.Context()
	//p, _ := NewChannelPool(0, 30, factory)
	p, _ := newChannelPool()
	conn, _ := p.Borrow(ctx)
	cli := conn.(*rpc.Client)
	var resp int
	err := cli.Call("Arith.Multiply", Args{1, 2}, &resp)
	if nil != err {
		t.Error(err)
	}
	if resp != 2 {
		t.Error("rpc.err")
	}
}

func TestPoolConcurrent2(t *testing.T) {
	ctx := mpc.Context()
	//p, _ := NewChannelPool(0, 30, factory)
	p, _ := newChannelPool()

	var wg sync.WaitGroup

	go func() {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				conn, _ := p.Borrow(ctx)
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				p.Close(ctx, conn)
				wg.Done()
			}(i)
		}
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			conn, _ := p.Borrow(ctx)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			p.Close(ctx, conn)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

//
//func TestPoolConcurrent3(t *testing.T) {
//	p, _ := NewChannelPool(0, 1, factory)
//
//	var wg sync.WaitGroup
//
//	wg.Add(1)
//	go func() {
//		p.Close()
//		wg.Done()
//	}()
//
//	if conn, err := p.Get(); err == nil {
//		conn.Close()
//	}
//
//	wg.Wait()
//}

func newChannelPool() (Pool, error) {
	pconf := PoolConfig{InitialCap: InitialCap, MaxCap: MaximumCap, Factory: factory, IdleTimeout: time.Second * 20,
		MaxIdle: MaxIdleCap}
	return NewChannelPool(mpc.Context(), &pconf)
}

func rpcServer() {
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", address)
	if e != nil {
		panic(e)
	}
	go http.Serve(l, nil)
}

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}
