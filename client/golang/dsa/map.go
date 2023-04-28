/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package dsa

import (
	"reflect"
	"sync"
	"unsafe"
)

var _ Map[any, any] = new(ConcurrentMap[any, any])

type Map[K any, V any] interface {
	IsEmpty() bool
	Size() int
	Exist(k K) bool
	Get(k K) (V, bool)
	Put(k K, v V)
	Remove(k K)
	// RemoveIfy deadlock warning
	RemoveIfy(key K, fn func(key K, v V, exist bool) bool) bool
	Clear()
	Keys() []K
	Values() []V
	Entries() []Entry[K, V]
	// ForEach deadlock warning
	ForEach(fn func(key K, v V))
	// FindAny deadlock warning
	FindAny(fn func(key K, v V) bool) V
	// Update deadlock warning
	Update(k K, fn func(k K, v V) (V, error)) (V, error)
	Pop(key K) (v V, exists bool)
	PutIfa(k K, v V) V
	// PutIfy deadlock warning
	PutIfy(k K, fn func(k K) V) V
	PutIfe(k K, fn func(k K) (V, error)) (V, error)
}

var ShardCount = 32

func FNV32Hash(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

// AnyHash returns hash taken from any object
func AnyHash[T any](i T) uint32 {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		if !v.CanAddr() {
			return FNV32Hash("")
		}

		v = v.Addr()
	}
	if v.IsNil() {
		return FNV32Hash(v.String())
	}

	size := unsafe.Sizeof(v.Interface())
	b := (*[1 << 10]uint8)(unsafe.Pointer(v.Pointer()))[:size:size]

	return FNV32Hash(string(b))
}

func NewAnyMap[K any, V any]() Map[K, V] {
	return NewMap[K, V](AnyHash[K])
}

func NewStringMap[V any]() Map[string, V] {
	return NewMap[string, V](FNV32Hash)
}

func NewMap[K any, V any](sharding func(key K) uint32) Map[K, V] {
	m := &ConcurrentMap[K, V]{
		sharding: sharding,
		shards:   make([]*MapEntries[K, V], ShardCount),
	}
	for idx := 0; idx < ShardCount; idx++ {
		m.shards[idx] = &MapEntries[K, V]{items: make(map[interface{}]V)}
	}
	return m
}

// ConcurrentMap
// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (ShardCount) map shards.
type ConcurrentMap[K any, V any] struct {
	shards   []*MapEntries[K, V]
	sharding func(key K) uint32
}

// IsEmpty
// IsEmpty returns the number of elements within the map.
func (that *ConcurrentMap[K, V]) IsEmpty() bool {
	return that.Size() < 1
}

// Size
// Size returns the number of elements within the map.
func (that *ConcurrentMap[K, V]) Size() int {
	count := 0
	for idx := 0; idx < ShardCount; idx++ {
		func() {
			shard := that.shards[idx]
			shard.RLock()
			defer shard.RUnlock()
			count += len(shard.items)
		}()
	}
	return count
}

// Exist
// Exist up an item under specified key
func (that *ConcurrentMap[K, V]) Exist(key K) bool {
	// Get shard
	shard := that.GetEntry(key)
	shard.RLock()
	defer shard.RUnlock()
	// See if element is within shard.
	_, ok := shard.items[key]
	return ok
}

// Get
// Get retrieves an element from map under given key.
func (that *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	// Get shard
	shard := that.GetEntry(key)
	shard.RLock()
	defer shard.RUnlock()
	// Get item from shard.
	val, ok := shard.items[key]
	return val, ok
}

// Put
// Put the given value under the specified key.
func (that *ConcurrentMap[K, V]) Put(key K, value V) {
	// Get map shard.
	shard := that.GetEntry(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

// Remove
// Remove removes an element from the map.
func (that *ConcurrentMap[K, V]) Remove(key K) {
	// Try to get shard.
	shard := that.GetEntry(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

// RemoveIfy locks the shard containing the key, retrieves its current value and calls the callback with those params
// If callback returns true and element exists, it will remove it from the map
// Returns the value returned by the callback (even if element was not present in the map)
// is a callback executed in a map.RemoveCb() call, while Lock is held
// If returns true, the element will be removed from the map
func (that *ConcurrentMap[K, V]) RemoveIfy(key K, fn func(key K, v V, exist bool) bool) bool {
	// Try to get shard.
	shard := that.GetEntry(key)
	shard.Lock()
	defer shard.Unlock()
	v, ok := shard.items[key]
	remove := fn(key, v, ok)
	if remove && ok {
		delete(shard.items, key)
	}
	return remove
}

// Clear removes all items from map.
func (that *ConcurrentMap[K, V]) Clear() {
	for item := range that.IterBuffered() {
		that.Remove(item.Key)
	}
}

// Keys returns all keys as []string
func (that *ConcurrentMap[K, V]) Keys() []K {
	count := that.Size()
	ch := make(chan K, count)
	go func() {
		// Foreach shard.
		wg := sync.WaitGroup{}
		wg.Add(ShardCount)
		for _, shard := range that.shards {
			go func(shard *MapEntries[K, V]) {
				// Foreach key, value pair.
				shard.RLock()
				defer func() {
					shard.RUnlock()
					wg.Done()
				}()
				for key := range shard.items {
					ch <- key.(K)
				}
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	// Generate keys
	keys := make([]K, 0, count)
	for k := range ch {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all items as map[string]V
func (that *ConcurrentMap[K, V]) Values() []V {
	var vs []V
	// Insert items to temporary map.
	for item := range that.IterBuffered() {
		vs = append(vs, item.Value)
	}
	return vs
}

// Entries returns all items as map[string]V
func (that *ConcurrentMap[K, V]) Entries() []Entry[K, V] {
	var entries []Entry[K, V]
	// Insert items to temporary map.
	for item := range that.IterBuffered() {
		entries = append(entries, item)
	}
	return entries
}

// PutIfa
// Put the given value under the specified key if no value was associated with it.
func (that *ConcurrentMap[K, V]) PutIfa(k K, v V) V {
	return that.PutIfy(k, func(key K) V { return v })
}

// PutIfy
// Compute the given value under the specified key if no value was associated with it.
func (that *ConcurrentMap[K, V]) PutIfy(k K, fn func(k K) V) V {
	// Get map shard.
	shard := that.GetEntry(k)
	shard.Lock()
	defer shard.Unlock()
	ov, ok := shard.items[k]
	if ok {
		return ov
	}
	nv := fn(k)
	shard.items[k] = nv
	return nv
}

func (that *ConcurrentMap[K, V]) PutIfe(k K, fn func(k K) (V, error)) (V, error) {
	// Get map shard.
	shard := that.GetEntry(k)
	shard.Lock()
	defer shard.Unlock()
	ov, ok := shard.items[k]
	if ok {
		return ov, nil
	}
	nv, err := fn(k)
	if nil != err {
		return nv, err
	}
	shard.items[k] = nv
	return nv, nil
}

// IterBuffered returns a buffered iterator which could be used in a for range loop.
func (that *ConcurrentMap[K, V]) IterBuffered() <-chan Entry[K, V] {
	cs := that.Snapshot()
	total := 0
	for _, c := range cs {
		total += cap(c)
	}
	ch := make(chan Entry[K, V], total)
	go that.Pip(cs, ch)
	return ch
}

// ForEach
// Callback based iterator, cheapest way to read
// all elements in a map.
// Iterator callbacalled for every key,value found in
// maps. RLock is held for all calls for a given shard
// therefore callback sess consistent view of a shard,
// but not across the shards
func (that *ConcurrentMap[K, V]) ForEach(fn func(key K, v V)) {
	for idx := range that.shards {
		shard := (that.shards)[idx]
		func() {
			shard.RLock()
			defer shard.RUnlock()
			for key, value := range shard.items {
				fn(key.(K), value)
			}
		}()
	}
}

// FindAny find any match return one
func (that *ConcurrentMap[K, V]) FindAny(fn func(key K, v V) bool) V {
	var mv V
	for idx := range that.shards {
		shard := (that.shards)[idx]
		if v, ok := func() (V, bool) {
			shard.RLock()
			defer shard.RUnlock()
			for key, value := range shard.items {
				if fn(key.(K), value) {
					return value, true
				}
			}
			return mv, false
		}(); ok {
			return v
		}
	}
	return mv
}

// Update
// Insert or Update - updates existing element or inserts a new one using UpsertCb
// Callback to return new element to be inserted into the map
// It is called while lock is held, therefore it MUST NOT
// try to access other keys in same map, as it can lead to deadlock since
// Go sync.RWLock is not reentrant
func (that *ConcurrentMap[K, V]) Update(k K, fn func(k K, v V) (V, error)) (V, error) {
	shard := that.GetEntry(k)
	shard.Lock()
	defer shard.Unlock()
	ov, _ := shard.items[k]
	r, err := fn(k, ov)
	if nil != err {
		return r, err
	}
	shard.items[k] = r
	return r, nil
}

// Pop removes an element from the map and returns it
func (that *ConcurrentMap[K, V]) Pop(key K) (V, bool) {
	// Try to get shard.
	shard := that.GetEntry(key)
	shard.Lock()
	defer shard.Unlock()
	v, ok := shard.items[key]
	delete(shard.items, key)
	return v, ok
}

// Pip reads elements from channels `cs` into channel `out`
func (that *ConcurrentMap[K, V]) Pip(cs []chan Entry[K, V], out chan Entry[K, V]) {
	wg := sync.WaitGroup{}
	wg.Add(len(cs))
	for _, ch := range cs {
		go func(ch chan Entry[K, V]) {
			defer wg.Done()
			for t := range ch {
				out <- t
			}
		}(ch)
	}
	wg.Wait()
	close(out)
}

// GetEntry returns shard under given key
func (that *ConcurrentMap[K, V]) GetEntry(key K) *MapEntries[K, V] {
	return that.shards[uint(that.sharding(key))%uint(ShardCount)]
}

// Snapshot
// Returns a array of channels that contains elements in each shard,
// which likely takes a snapshot of `m`.
// It returns once the size of each buffered channel is determined,
// before all the channels are populated using goroutines.
func (that *ConcurrentMap[K, V]) Snapshot() (cs []chan Entry[K, V]) {
	//When you access map items before initializing.
	if len(that.shards) == 0 {
		panic(`cmap.ConcurrentMap is not initialized. Should run New() before usage.`)
	}
	cs = make([]chan Entry[K, V], ShardCount)
	wg := sync.WaitGroup{}
	wg.Add(ShardCount)
	// Foreach shard.
	for index, shard := range that.shards {
		go func(index int, shard *MapEntries[K, V]) {
			// Foreach key, value pair.
			shard.RLock()
			defer func() {
				shard.RUnlock()
				wg.Done()
			}()
			cs[index] = make(chan Entry[K, V], len(shard.items))
			for key, val := range shard.items {
				cs[index] <- Entry[K, V]{key.(K), val}
			}
			close(cs[index])
		}(index, shard)
	}
	wg.Wait()
	return cs
}

// MapEntries
// A "thread" safe string to anything map.
type MapEntries[K any, V any] struct {
	items        map[any]V
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// Entry
// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type Entry[K any, V any] struct {
	Key   K
	Value V
}
