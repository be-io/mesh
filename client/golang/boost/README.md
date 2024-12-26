```go
package main

import (
	"fmt"
	"net"
	"time"
	"github.com/opendatav/mesh/client/golang/pool"
)

func main() {
	//factory Specify the method to create the connection
	factory := func() (interface{}, error) { return net.Dial("tcp", "127.0.0.1:4000") }

	//close Specify the method to close the connection
	close := func(v interface{}) error { return v.(net.Conn).Close() }

	//ping Specify the method to detect whether the connection is invalid
	//ping := func(v interface{}) error { return nil }

	//Create a connection pool: Initialize the number of connections to 5, the maximum idle connection is 20, and the maximum concurrent connection is 30
	poolConfig := &pool.Config{
		InitialCap: 5,
		MaxIdle:    20,
		MaxCap:     30,
		Factory:    factory,
		Close:      close,
		//Ping:       ping,
		//The maximum idle time of the connection, the connection exceeding this time will be closed, which can avoid the problem of automatic failure when connecting to EOF when idle
		IdleTimeout: 15 * time.Second,
	}
	p, err := pool.NewChannelPool(poolConfig)
	if nil != err {
		fmt.Println("err=", err)
	}

	//Get a connection from the connection pool
	v, err := p.Get()

	//do something
	//conn=v.(net.Conn)

	//Put the connection back into the connection pool, when the connection is no longer in use
	p.Put(v)

	//Release all connections in the connection pool, when resources need to be destroyed
	p.Release()

	//View the number of connections in the current connection pool
	current := p.Len()
}
```