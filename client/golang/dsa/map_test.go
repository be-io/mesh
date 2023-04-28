/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package dsa

import (
	"hash/fnv"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"testing"
)

type Animal struct {
	name string
}

func TestMapCreation(t *testing.T) {
	m := NewStringMap[string]()
	if nil != m.Values() {
		t.Error("map is not nil.")
	}

	if m.Size() != 0 {
		t.Error("new map should be empty.")
	}
}

func TestInsert(t *testing.T) {
	m := NewStringMap[Animal]()
	elephant := Animal{"elephant"}
	monkey := Animal{"monkey"}

	m.Put("elephant", elephant)
	m.Put("monkey", monkey)

	if m.Size() != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

func TestInsertAbsent(t *testing.T) {
	m := NewStringMap[Animal]()
	elephant := Animal{"elephant"}
	monkey := Animal{"monkey"}

	m.PutIfa("elephant", elephant)
	m.PutIfa("elephant", monkey)
	if v, ok := m.Get("elephant"); !ok || v != elephant {
		t.Error("map set a new value even the entry is already present")
	}
}

func TestGet(t *testing.T) {
	m := NewStringMap[Animal]()

	// Get a missing element.
	val, ok := m.Get("Money")

	if ok == true {
		t.Error("ok should be false when item is missing from map.")
	}

	if (val != Animal{}) {
		t.Error("Missing values should return as null.")
	}

	elephant := Animal{"elephant"}
	m.Put("elephant", elephant)

	// Retrieve inserted element.
	elephant, ok = m.Get("elephant")
	if ok == false {
		t.Error("ok should be true for item stored within the map.")
	}

	if elephant.name != "elephant" {
		t.Error("item was modified.")
	}
}

func TestHas(t *testing.T) {
	m := NewStringMap[Animal]()

	// Get a missing element.
	if m.Exist("Money") == true {
		t.Error("element shouldn't exists")
	}

	elephant := Animal{"elephant"}
	m.Put("elephant", elephant)

	if m.Exist("elephant") == false {
		t.Error("element exists, expecting Has to return True.")
	}
}

func TestRemove(t *testing.T) {
	m := NewStringMap[Animal]()

	monkey := Animal{"monkey"}
	m.Put("monkey", monkey)

	m.Remove("monkey")

	if m.Size() != 0 {
		t.Error("Expecting count to be zero once item was removed.")
	}

	temp, ok := m.Get("monkey")

	if ok != false {
		t.Error("Expecting ok to be false for missing items.")
	}

	if (temp != Animal{}) {
		t.Error("Expecting item to be nil after its removal.")
	}

	// Remove a none existing element.
	m.Remove("noone")
}

func TestRemoveCb(t *testing.T) {
	m := NewStringMap[Animal]()

	monkey := Animal{"monkey"}
	m.Put("monkey", monkey)
	elephant := Animal{"elephant"}
	m.Put("elephant", elephant)

	var (
		mapKey   string
		mapVal   Animal
		wasFound bool
	)
	cb := func(key string, val Animal, exists bool) bool {
		mapKey = key
		mapVal = val
		wasFound = exists

		return val.name == "monkey"
	}

	// Monkey should be removed
	result := m.RemoveIfy("monkey", cb)
	if !result {
		t.Errorf("Result was not true")
	}

	if mapKey != "monkey" {
		t.Error("Wrong key was provided to the callback")
	}

	if mapVal != monkey {
		t.Errorf("Wrong value was provided to the value")
	}

	if !wasFound {
		t.Errorf("Key was not found")
	}

	if m.Exist("monkey") {
		t.Errorf("Key was not removed")
	}

	// Elephant should not be removed
	result = m.RemoveIfy("elephant", cb)
	if result {
		t.Errorf("Result was true")
	}

	if mapKey != "elephant" {
		t.Error("Wrong key was provided to the callback")
	}

	if mapVal != elephant {
		t.Errorf("Wrong value was provided to the value")
	}

	if !wasFound {
		t.Errorf("Key was not found")
	}

	if !m.Exist("elephant") {
		t.Errorf("Key was removed")
	}

	// Unset key should remain unset
	result = m.RemoveIfy("horse", cb)
	if result {
		t.Errorf("Result was true")
	}

	if mapKey != "horse" {
		t.Error("Wrong key was provided to the callback")
	}

	if (mapVal != Animal{}) {
		t.Errorf("Wrong value was provided to the value")
	}

	if wasFound {
		t.Errorf("Key was found")
	}

	if m.Exist("horse") {
		t.Errorf("Key was created")
	}
}

func TestPop(t *testing.T) {
	m := NewStringMap[Animal]()

	monkey := Animal{"monkey"}
	m.Put("monkey", monkey)

	v, exists := m.Pop("monkey")

	if !exists || v != monkey {
		t.Error("Pop didn't find a monkey.")
	}

	v2, exists2 := m.Pop("monkey")

	if exists2 || v2 == monkey {
		t.Error("Pop keeps finding monkey")
	}

	if m.Size() != 0 {
		t.Error("Expecting count to be zero once item was Pop'ed.")
	}

	temp, ok := m.Get("monkey")

	if ok != false {
		t.Error("Expecting ok to be false for missing items.")
	}

	if (temp != Animal{}) {
		t.Error("Expecting item to be nil after its removal.")
	}
}

func TestCount(t *testing.T) {
	m := NewStringMap[Animal]()
	for i := 0; i < 100; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}

	if m.Size() != 100 {
		t.Error("Expecting 100 element within map.")
	}
}

func TestIsEmpty(t *testing.T) {
	m := NewStringMap[Animal]()

	if m.IsEmpty() == false {
		t.Error("new map should be empty")
	}

	m.Put("elephant", Animal{"elephant"})

	if m.IsEmpty() != false {
		t.Error("map shouldn't be empty.")
	}
}

func TestIterator(t *testing.T) {
	m := NewStringMap[Animal]()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}

	counter := 0
	// Iterate over elements.
	for _, val := range m.Entries() {
		if (val.Value == Animal{}) {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestBufferedIterator(t *testing.T) {
	m := NewStringMap[Animal]()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}

	counter := 0
	// Iterate over elements.
	for _, item := range m.Entries() {
		val := item.Value

		if (val == Animal{}) {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestClear(t *testing.T) {
	m := NewStringMap[Animal]()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}

	m.Clear()

	if m.Size() != 0 {
		t.Error("We should have 0 elements.")
	}
}

func TestIterCb(t *testing.T) {
	m := NewStringMap[Animal]()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}

	counter := 0
	// Iterate over elements.
	m.ForEach(func(key string, v Animal) {
		counter++
	})
	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestItems(t *testing.T) {
	m := NewStringMap[Animal]()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}

	items := m.Entries()

	if len(items) != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestConcurrent(t *testing.T) {
	m := NewStringMap[int]()
	ch := make(chan int)
	const iterations = 1000
	var a [iterations]int

	// Using go routines insert 1000 ints into our map.
	go func() {
		for i := 0; i < iterations/2; i++ {
			// Add item to map.
			m.Put(strconv.Itoa(i), i)

			// Retrieve item from map.
			val, _ := m.Get(strconv.Itoa(i))

			// Write to channel inserted value.
			ch <- val
		} // Call go routine with current index.
	}()

	go func() {
		for i := iterations / 2; i < iterations; i++ {
			// Add item to map.
			m.Put(strconv.Itoa(i), i)

			// Retrieve item from map.
			val, _ := m.Get(strconv.Itoa(i))

			// Write to channel inserted value.
			ch <- val
		} // Call go routine with current index.
	}()

	// Wait for all go routines to finish.
	counter := 0
	for elem := range ch {
		a[counter] = elem
		counter++
		if counter == iterations {
			break
		}
	}

	// Sorts array, will make is simpler to verify all inserted values we're returned.
	sort.Ints(a[0:iterations])

	// Make sure map contains 1000 elements.
	if m.Size() != iterations {
		t.Error("Expecting 1000 elements.")
	}

	// Make sure all inserted values we're fetched from map.
	for i := 0; i < iterations; i++ {
		if i != a[i] {
			t.Error("missing value", i)
		}
	}
}

func TestKeys(t *testing.T) {
	m := NewStringMap[Animal]()

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}

	keys := m.Keys()
	if len(keys) != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestMInsert(t *testing.T) {
	animals := map[string]Animal{
		"elephant": {"elephant"},
		"monkey":   {"monkey"},
	}
	m := NewStringMap[Animal]()
	for k, v := range animals {
		m.Put(k, v)
	}

	if m.Size() != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

func TestFnv32(t *testing.T) {
	key := []byte("ABC")

	hasher := fnv.New32()
	_, err := hasher.Write(key)
	if err != nil {
		t.Errorf(err.Error())
	}
	if FNV32Hash(string(key)) != hasher.Sum32() {
		t.Errorf("Bundled fnv32 produced %d, expected result from hash/fnv32 is %d", FNV32Hash(string(key)), hasher.Sum32())
	}

}

func TestUpsert(t *testing.T) {
	dolphin := Animal{"dolphin"}
	whale := Animal{"whale"}
	tiger := Animal{"tiger"}
	lion := Animal{"lion"}

	m := NewStringMap[Animal]()
	m.Put("marine", dolphin)
	m.Update("marine", func(k string, v Animal) (Animal, error) {
		return whale, nil
	})
	m.Update("predator", func(k string, v Animal) (Animal, error) {
		return tiger, nil
	})
	m.Update("predator", func(k string, v Animal) (Animal, error) {
		return lion, nil
	})

	if m.Size() != 2 {
		t.Error("map should contain exactly two elements.")
	}

	marineAnimals, ok := m.Get("marine")
	if marineAnimals.name != "dolphinwhale" || !ok {
		t.Error("Set, then Upsert failed")
	}

	predators, ok := m.Get("predator")
	if !ok || predators.name != "tigerlion" {
		t.Error("Upsert, then Upsert failed")
	}
}

func TestKeysWhenRemoving(t *testing.T) {
	m := NewStringMap[Animal]()

	// Insert 100 elements.
	Total := 100
	for i := 0; i < Total; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}

	// Remove 10 elements concurrently.
	Num := 10
	for i := 0; i < Num; i++ {
		go func(c Map[string, Animal], n int) {
			c.Remove(strconv.Itoa(n))
		}(m, i)
	}
	keys := m.Keys()
	for _, k := range keys {
		if k == "" {
			t.Error("Empty keys returned")
		}
	}
}

func TestUnDrainedIter(t *testing.T) {
	m := NewStringMap[Animal]()
	// Insert 100 elements.
	Total := 100
	for i := 0; i < Total; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	counter := 0
	// Iterate over elements.
	ch := m.Entries()
	for _, item := range ch {
		val := item.Value

		if (val == Animal{}) {
			t.Error("Expecting an object.")
		}
		counter++
		if counter == 42 {
			break
		}
	}
	for i := Total; i < 2*Total; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	for _, item := range ch {
		val := item.Value

		if (val == Animal{}) {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 142 {
		t.Error("We should have been right where we stopped", counter, len(m.Keys()))
	}

	counter = 0
	for _, item := range m.Entries() {
		val := item.Value

		if (val == Animal{}) {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 200 {
		t.Error("We should have counted 200 elements.")
	}
}

func TestUnDrainedIterBuffered(t *testing.T) {
	m := NewStringMap[Animal]()
	// Insert 100 elements.
	Total := 100
	for i := 0; i < Total; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	counter := 0
	// Iterate over elements.
	ch := m.Entries()
	for _, item := range ch {
		val := item.Value

		if (val == Animal{}) {
			t.Error("Expecting an object.")
		}
		counter++
		if counter == 42 {
			break
		}
	}
	for i := Total; i < 2*Total; i++ {
		m.Put(strconv.Itoa(i), Animal{strconv.Itoa(i)})
	}
	for _, item := range ch {
		val := item.Value

		if (val == Animal{}) {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 142 {
		t.Error("We should have been right where we stopped")
	}

	counter = 0
	for _, item := range m.Entries() {
		val := item.Value

		if (val == Animal{}) {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 200 {
		t.Error("We should have counted 200 elements.")
	}
}

func TestPutGet(t *testing.T) {
	wg := sync.WaitGroup{}
	m := NewAnyMap[any, any]()
	for index := 0; index < 20; index++ {
		t.Log("Routine", index)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := 0; idx < 99999; idx++ {
				k := rand.Int63()
				v := rand.Float32()
				m.Put(k, v)
				nv, ok := m.Get(k)
				if !ok {
					t.Error("Unexpected concurrency read write ")
					return
				}
				if nv != v {
					t.Error("Unexpected concurrency read write compare")
					return
				}
			}
		}()
	}
	wg.Wait()
	t.Log("Test success")
}
