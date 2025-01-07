/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package tool

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/opendatav/mesh/client/golang/cause"
	"github.com/opendatav/mesh/client/golang/log"
	"github.com/opendatav/mesh/client/golang/types"
	"hash/fnv"
	"io"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"
)

func Anyone[T any](vs ...T) T {
	if len(vs) < 1 {
		var zero T
		return zero
	}
	for _, v := range vs {
		if reflect.ValueOf(&v).Elem().IsZero() {
			continue
		}
		return v
	}
	return vs[0]
}

func Ternary[T any](expr bool, left T, right T) T {
	if expr {
		return left
	}
	return right
}

func FistUpper(v string) string {
	if len(v) < 1 {
		return v
	}
	return fmt.Sprintf("%c%s", unicode.ToUpper(rune(v[0])), v[1:])
}

func FistLower(v string) string {
	if len(v) < 1 {
		return v
	}
	return fmt.Sprintf("%c%s", unicode.ToLower(rune(v[0])), v[1:])
}

func IsLocal(nodeId string) bool {
	return strings.ToUpper(nodeId) == types.LocalNodeId || strings.ToUpper(nodeId) == types.LocalInstId
}

func IsLocalEnv(environ *types.Environ, nodeIds ...string) bool {
	for _, nodeId := range nodeIds {
		if IsLocal(nodeId) {
			return true
		}
		lowerId := strings.ToLower(nodeId)
		if strings.Index(lowerId, "jg") == 0 {
			instId, err := types.FromInstID(lowerId)
			if nil != err {
				log.Warn0("%s is not standard identity", nodeId)
				continue
			}
			if instId.MatchNode(environ.NodeId) {
				return true
			}
		}
		if lowerId == strings.ToLower(environ.NodeId) {
			return true
		}
	}
	return false
}

func NewTraceId() string {
	return fmt.Sprintf("%s%s", IPHex.Get(), NextID())
}

func NewSpanId(spanId string, calls int) string {
	if "" == spanId {
		return "0"
	}
	if len(spanId) > 255 {
		return "0"
	}
	return fmt.Sprintf("%s.%d", spanId, calls)
}

func Contains[T comparable](vs []T, v T) bool {
	for _, x := range vs {
		if x == v {
			return true
		}
	}
	return false
}

// ContainsAny if source contain any candidate
func ContainsAny(source string, candidates ...string) bool {
	if len(candidates) == 0 {
		return false
	}
	for _, candidate := range candidates {
		if strings.Contains(source, candidate) {
			return true
		}
	}
	return false
}

func Must(b bool) {
	if b {
		return
	}
	log.Panic(cause.Errorf("assertion failed"))
}

func MustNoError(err error) {
	if err == nil {
		return
	}
	log.Panic(err)
}

func CheckAvailable(ctx context.Context, address string) bool {
	if connection, err := net.DialTimeout("tcp", address, 1*time.Second); nil != err {
		log.Warn(ctx, "Test address is reachable, %s", err.Error())
		return false
	} else {
		log.Catch(connection.Close())
		return true
	}
}

func IPInHex(ip string) string {
	traceId := bytes.Buffer{}
	for _, feg := range strings.Split(ip, ".") {
		digit, _ := strconv.Atoi(feg)
		traceId.WriteString(fmt.Sprintf("%x", digit))
	}
	return traceId.String()
}

func MakeDir(path string) error {
	_, err := os.Stat(path)
	if nil == err || !os.IsNotExist(err) {
		return nil
	}
	return cause.Error(os.MkdirAll(path, 0755))
}

func MakeFile(path string) error {
	_, err := os.Stat(path)
	if nil == err || !os.IsNotExist(err) {
		return nil
	}
	fd, err := os.Create(path)
	if nil != err {
		return cause.Error(err)
	}
	return cause.Error(fd.Close())
}

func Padding(v string, length int, x string) string {
	if len(v) == length {
		return v
	}
	if len(v) < length {
		return strings.Repeat(x, length-len(v)) + v
	}
	return v[0:length]
}

func Addressable(ctx context.Context, addr string) bool {
	if "" == addr {
		return false
	}
	uri, err := types.FormatURL(addr)
	if nil != err {
		log.Debug(ctx, err.Error())
		return false
	}
	if "" != uri.Query().Get("p") || "1" == uri.Query().Get("k") {
		return true
	}
	port := Ternary("" != uri.Port(), uri.Port(), Ternary(uri.Scheme == "https", "443", "80"))
	return CheckAvailable(ctx, fmt.Sprintf("%s:%s", uri.Hostname(), port))
}

func FieldBy[T any](any interface{}, name string) interface{} {
	value := reflect.ValueOf(any).Elem().FieldByName(name)
	readable := reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr())).Elem()
	return readable.Interface()
}

func MapBy[K comparable, V any](m map[K]V, key K) (v V) {
	if nil == m {
		return
	}
	return m[key]
}

func ReplaceName(rawName string) string {
	if strings.Contains(rawName, pandora) {
		return pandora
	}
	return rawName
}

func Copy(src string, dst string) error {
	sf, err := os.Open(src)
	if nil != err {
		return cause.Error(err)
	}
	defer func() { log.Catch(sf.Close()) }()
	df, err := os.Create(dst)
	if nil != err {
		return cause.Error(err)
	}
	defer func() { log.Catch(df.Close()) }()
	_, err = io.Copy(df, sf)
	return cause.Error(err)
}

func Merge[K comparable, V any](dicts ...map[K]V) map[K]V {
	dict := map[K]V{}
	for _, kvs := range dicts {
		for k, v := range kvs {
			dict[k] = v
		}
	}
	return dict
}

func Distinct[T comparable](vs []T) []T {
	var nvs []T
	for _, v := range vs {
		if !Contains(nvs, v) {
			nvs = append(nvs, v)
		}
	}
	return nvs
}

func DistinctCall[T any, K comparable](vs []T, deeper func(v T) K) []T {
	var nvs []T
	keys := map[K]bool{}
	for _, v := range vs {
		k := deeper(v)
		if ok := keys[k]; ok {
			continue
		}
		keys[k] = true
		nvs = append(nvs, v)
	}
	return nvs
}

func Hash(s string) string {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(s))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func Reverse[T any](vs []T) []T {
	for x, y := 0, len(vs)-1; x < y; x, y = x+1, y-1 {
		vs[x], vs[y] = vs[y], vs[x]
	}
	return vs
}

func Values[K comparable, V any](dicts ...map[K]V) []V {
	var list []V
	for _, kvs := range dicts {
		for _, v := range kvs {
			list = append(list, v)
		}
	}
	return list
}

func MapValues[K comparable, V any, T any](m func(v V) T, dicts ...map[K]V) []T {
	var list []T
	for _, kvs := range dicts {
		for _, v := range kvs {
			list = append(list, m(v))
		}
	}
	return list
}

func WithoutZero[K comparable](lists ...[]K) []K {
	var list []K
	for _, vs := range lists {
		for _, v := range vs {
			if !reflect.ValueOf(&v).Elem().IsZero() {
				list = append(list, v)
			}
		}
	}
	return list
}

func EncodeB64(cipher []byte) ([]byte, error) {
	dst := make([]byte, hex.EncodedLen(len(cipher)))
	hex.Encode(dst, cipher)
	return dst, nil
}

func DecodeB64(cipher []byte) ([]byte, error) {
	n, err := hex.Decode(cipher, cipher)
	return cipher[:n], cause.Error(err)
}

func EscapedPath(r *url.URL) string {
	path := r.RawPath
	if path == "" {
		path = r.EscapedPath()
	}
	return path
}

func Timestamp(ctx context.Context, v string) int64 {
	if "" == v {
		return time.Now().UnixMilli()
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if nil != err {
		log.Error(ctx, err.Error())
		return time.Now().UnixMilli()
	}
	return n
}
