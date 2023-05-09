/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package broker

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/be-io/mesh/client/golang/cause"
	"github.com/be-io/mesh/client/golang/grpc"
	"io"
	"math"
	"net/http"
	"strings"

	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
)

const (
	maxMessageSize = grpc.MaxSize // Not limited with grpc
)

func writeSizePreface(w io.Writer, sz int32) error {
	return binary.Write(w, binary.BigEndian, sz)
}

func writeMessage(w io.Writer, codec encoding.Codec, m interface{}, end bool) error {
	b, err := codec.Marshal(m)
	if nil != err {
		return cause.Error(err)
	}

	sz := len(b)
	if sz > math.MaxInt32 {
		return cause.Errorf("message too large to send: %d bytes", sz)
	}
	if end {
		// trailer message is indicated w/ negative size
		sz = -sz
	}
	err = writeSizePreface(w, int32(sz))
	if nil != err {
		return cause.Error(err)
	}

	_, err = w.Write(b)
	if err == nil {
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
	return cause.Error(err)
}

func readSizePreface(in io.Reader) (int32, error) {
	var sz int32
	err := binary.Read(in, binary.BigEndian, &sz)
	return sz, cause.Error(err)
}

func readMessage(in io.Reader, codec encoding.Codec, sz int32, m interface{}) error {
	if sz < 0 {
		return cause.Errorf("bad size preface: size cannot be negative: %d", sz)
	} else if sz > maxMessageSize {
		return cause.Errorf("bad size preface: indicated size is too large: %d", sz)
	}
	msg := make([]byte, sz)
	_, err := io.ReadAtLeast(in, msg, int(sz))
	if nil != err {
		return cause.Error(err)
	}
	return cause.Error(codec.Unmarshal(msg, m))
}

func asMetadata(header http.Header) (metadata.MD, error) {
	// metadata has same shape as http.Header,
	md := metadata.MD{}
	for k, vs := range header {
		k = strings.ToLower(k)
		for _, v := range vs {
			if strings.HasSuffix(k, "-bin") {
				vv, err := base64.URLEncoding.DecodeString(v)
				if nil != err {
					return nil, err
				}
				v = string(vv)
			}
			md[k] = append(md[k], v)
		}
	}
	return md, nil
}

var reservedHeaders = map[string]struct{}{
	"accept-encoding":   {},
	"connection":        {},
	"content-type":      {},
	"content-length":    {},
	"keep-alive":        {},
	"te":                {},
	"trailer":           {},
	"transfer-encoding": {},
	"upgrade":           {},
}

func toHeaders(md metadata.MD, h http.Header, prefix string) {
	// binary headers must be base-64-encoded
	for k, vs := range md {
		lowerK := strings.ToLower(k)
		if _, ok := reservedHeaders[lowerK]; ok {
			// ignore reserved header keys
			continue
		}
		isBin := strings.HasSuffix(lowerK, "-bin")
		for _, v := range vs {
			if isBin {
				v = base64.URLEncoding.EncodeToString([]byte(v))
			}
			h.Add(prefix+k, v)
		}
	}
}

type strAddr string

func (that strAddr) Network() string {
	if that != "" {
		// Per the documentation on net/http.Request.RemoteAddr, if this is
		// set, it's set to the IP:port of the peer (hence, TCP):
		// https://golang.org/pkg/net/http/#Request
		//
		// If we want to support Unix sockets later, we can
		// add our own grpc-specific convention within the
		// grpc codebase to set RemoteAddr to a different
		// format, or probably better: we can attach it to the
		// context and use that from serverHandlerTransport.RemoteAddr.
		return "tcp"
	}
	return ""
}

func (that strAddr) String() string { return string(that) }
