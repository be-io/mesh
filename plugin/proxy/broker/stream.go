/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package broker

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type UnaryServerTransportStream struct {
	// Name is the full method name in "/service/method" format.
	Name string

	mu       sync.Mutex
	hdrs     metadata.MD
	hdrsSent bool
	tlrs     metadata.MD
	tlrsSent bool
}

func (that *UnaryServerTransportStream) Method() string {
	return that.Name
}

func (that *UnaryServerTransportStream) Finish() {
	that.mu.Lock()
	defer that.mu.Unlock()
	that.hdrsSent = true
	that.tlrsSent = true
}

func (that *UnaryServerTransportStream) SetHeader(md metadata.MD) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	return that.setHeaderLocked(md)
}

func (that *UnaryServerTransportStream) SendHeader(md metadata.MD) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if err := that.setHeaderLocked(md); nil != err {
		return err
	}
	that.hdrsSent = true
	return nil
}

func (that *UnaryServerTransportStream) setHeaderLocked(md metadata.MD) error {
	if that.hdrsSent {
		return fmt.Errorf("headers already sent")
	}
	if that.hdrs == nil {
		that.hdrs = metadata.MD{}
	}
	for k, v := range md {
		that.hdrs[k] = append(that.hdrs[k], v...)
	}
	return nil
}

func (that *UnaryServerTransportStream) GetHeaders() metadata.MD {
	that.mu.Lock()
	defer that.mu.Unlock()
	return that.hdrs
}

func (that *UnaryServerTransportStream) SetTrailer(md metadata.MD) error {
	that.mu.Lock()
	defer that.mu.Unlock()
	if that.tlrsSent {
		return fmt.Errorf("trailers already sent")
	}
	if that.tlrs == nil {
		that.tlrs = metadata.MD{}
	}
	for k, v := range md {
		that.tlrs[k] = append(that.tlrs[k], v...)
	}
	return nil
}

func (that *UnaryServerTransportStream) GetTrailers() metadata.MD {
	that.mu.Lock()
	defer that.mu.Unlock()
	return that.tlrs
}

type ServerTransportStream struct {
	// Name is the full method name in "/service/method" format.
	Name string
	// Stream is the underlying stream to which header and trailer calls are
	// delegated.
	Stream grpc.ServerStream
}

func (that *ServerTransportStream) Method() string {
	return that.Name
}

func (that *ServerTransportStream) SetHeader(md metadata.MD) error {
	return that.Stream.SetHeader(md)
}

func (that *ServerTransportStream) SendHeader(md metadata.MD) error {
	return that.Stream.SendHeader(md)
}

func (that *ServerTransportStream) SetTrailer(md metadata.MD) error {
	type trailerWithErrors interface {
		TrySetTrailer(md metadata.MD) error
	}
	if t, ok := that.Stream.(trailerWithErrors); ok {
		return t.TrySetTrailer(md)
	}
	that.Stream.SetTrailer(md)
	return nil
}
