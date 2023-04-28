/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package boost

import (
	"context"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/log"
	"io"
	"sync"
	"time"
	//"reflect"
)

var (
	// ErrPoolClosed is the error resulting if the pool is closed via pool.Close().
	ErrPoolClosed = cause.Errorf("Pool is closed ")
	//ErrPoolMaxActiveConnReached connection pool limited
	ErrPoolMaxActiveConnReached = cause.Errorf("MaxActiveConnReached ")
	// ErrPoolTimeout is the error when the client pool timed out
	ErrPoolTimeout = cause.Errorf("Pool timed out ")
)

// Connection is the duck of pooled connection interface.
type Connection interface {
	io.Closer
}

// Pool interface describes a pool implementation. A pool should have maximum
// capacity. An ideal pool is threadsafe and easy to use.
type Pool interface {
	// Borrow returns a new connection from the pool. Closing the connections puts
	// it back to the Pool. Closing it when the pool is destroyed or full will
	// be counted as an error.
	Borrow(ctx context.Context) (Connection, error)

	// Return the connection to the pool
	Return(ctx context.Context, conn Connection) error

	// Close closes the pool and all its connections. After Close() the pool is
	// no longer usable.
	Close(ctx context.Context, conn Connection) error

	// Release the connections
	Release(ctx context.Context)

	// Len returns the current number of connections of the pool.
	Len() int
}

// PoolConfig is the connection pool configuration
type PoolConfig struct {
	//min connections
	InitialCap int
	//max connections
	MaxCap int
	//max idle connections
	MaxIdle int
	//connection factory
	Factory func(ctx context.Context) (Connection, error)
	//connection health check hook
	Ping func(ctx context.Context, conn Connection) error
	//max timeout of idle connection
	IdleTimeout time.Duration
}

type connReq struct {
	idleConn *idleConn
}

// channelPool connection pool
type channelPool struct {
	mu                       sync.RWMutex
	conns                    chan *idleConn
	factory                  func(ctx context.Context) (Connection, error)
	ping                     func(ctx context.Context, conn Connection) error
	idleTimeout, waitTimeOut time.Duration
	maxActive                int
	openingConns             int
	connReqs                 []chan connReq
}

type idleConn struct {
	conn Connection
	t    time.Time
}

// NewChannelPool will create a connection channel pool
func NewChannelPool(ctx context.Context, poolConfig *PoolConfig) (Pool, error) {
	if !(poolConfig.InitialCap <= poolConfig.MaxIdle && poolConfig.MaxCap >= poolConfig.MaxIdle && poolConfig.InitialCap >= 0) {
		return nil, cause.Errorf("Invalid capacity settings")
	}
	if nil == poolConfig.Factory {
		return nil, cause.Errorf("Invalid factory func settings")
	}

	c := &channelPool{
		conns:        make(chan *idleConn, poolConfig.MaxIdle),
		factory:      poolConfig.Factory,
		idleTimeout:  poolConfig.IdleTimeout,
		maxActive:    poolConfig.MaxCap,
		openingConns: poolConfig.InitialCap,
	}

	if nil != poolConfig.Ping {
		c.ping = poolConfig.Ping
	}

	for i := 0; i < poolConfig.InitialCap; i++ {
		conn, err := c.factory(ctx)
		if nil != err {
			c.Release(ctx)
			return nil, cause.Errorf("Factory is not able to fill the pool: %s ", err.Error())
		}
		c.conns <- &idleConn{conn: conn, t: time.Now()}
	}

	return c, nil
}

// getConns will get all connections
func (c *channelPool) getConns() chan *idleConn {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()
	return conns
}

// Borrow a connection from pool
func (c *channelPool) Borrow(ctx context.Context) (Connection, error) {
	conns := c.getConns()
	if nil == conns {
		return nil, ErrPoolClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			if nil == wrapConn {
				return nil, ErrPoolClosed
			}
			//discard if idle timeout
			if timeout := c.idleTimeout; timeout > 0 {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					log.Catch(c.Close(ctx, wrapConn.conn))
					continue
				}
			}
			//discard if ping failure
			if nil != c.ping {
				if err := c.Ping(ctx, wrapConn.conn); nil != err {
					log.Catch(c.Close(ctx, wrapConn.conn))
					continue
				}
			}
			return wrapConn.conn, nil
		case <-ctx.Done():
			return nil, ErrPoolTimeout
		default:
			c.mu.Lock()
			log.Debug(ctx, "OpenConn %v %v", c.openingConns, c.maxActive)
			if c.openingConns >= c.maxActive {
				req := make(chan connReq, 1)
				c.connReqs = append(c.connReqs, req)
				c.mu.Unlock()
				ret, ok := <-req
				if !ok {
					return nil, ErrPoolMaxActiveConnReached
				}
				if timeout := c.idleTimeout; timeout > 0 {
					if ret.idleConn.t.Add(timeout).Before(time.Now()) {
						log.Catch(c.Close(ctx, ret.idleConn.conn))
						continue
					}
				}
				return ret.idleConn.conn, nil
			}
			if nil == c.factory {
				c.mu.Unlock()
				return nil, ErrPoolClosed
			}
			conn, err := c.factory(ctx)
			if nil != err {
				c.mu.Unlock()
				return nil, err
			}
			c.openingConns++
			c.mu.Unlock()
			return conn, nil
		}
	}
}

// Return will return the connection to the pool
func (c *channelPool) Return(ctx context.Context, conn Connection) error {
	if nil == conn {
		return cause.Errorf("Connection is nil. rejecting")
	}

	c.mu.Lock()

	if nil == c.conns {
		c.mu.Unlock()
		return c.Close(ctx, conn)
	}

	if l := len(c.connReqs); l > 0 {
		req := c.connReqs[0]
		copy(c.connReqs, c.connReqs[1:])
		c.connReqs = c.connReqs[:l-1]
		req <- connReq{
			idleConn: &idleConn{conn: conn, t: time.Now()},
		}
		c.mu.Unlock()
		return nil
	} else {
		select {
		case c.conns <- &idleConn{conn: conn, t: time.Now()}:
			c.mu.Unlock()
			return nil
		default:
			c.mu.Unlock()
			// connection is fill, close it
			return c.Close(ctx, conn)
		}
	}
}

// Close the connection
func (c *channelPool) Close(ctx context.Context, conn Connection) error {
	if nil == conn {
		return cause.Errorf("Connection is nil. rejecting")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.openingConns--
	return conn.Close()
}

// Ping check the connection
func (c *channelPool) Ping(ctx context.Context, conn Connection) error {
	if nil == conn {
		return cause.Errorf("Connection is nil. rejecting")
	}
	return c.ping(ctx, conn)
}

// Release will release all connections
func (c *channelPool) Release(ctx context.Context) {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.ping = nil
	c.mu.Unlock()

	if nil == conns {
		return
	}

	close(conns)
	for wrapConn := range conns {
		log.Catch(wrapConn.conn.Close())
	}
}

// Len return the size of the connections in pool
func (c *channelPool) Len() int {
	return len(c.getConns())
}
