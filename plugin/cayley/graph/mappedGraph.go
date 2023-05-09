package graph

import (
	"errors"
	"fmt"
)

// GetMapped gets the key assosciated with the mapped int
func (that *Graph) GetMapped(a int) (string, error) {
	if !that.usingMap || that.mapping == nil {
		return "", ErrNoMap
	}
	for k, v := range that.mapping {
		if v == a {
			return k, nil
		}
	}
	return "", errors.New(fmt.Sprint(a, " not found in mapping"))
}

// GetMapping gets the index associated with the specified key
func (that *Graph) GetMapping(a string) (int, error) {
	if !that.usingMap || that.mapping == nil {
		return -1, ErrNoMap
	}
	if b, ok := that.mapping[a]; ok {
		return b, nil
	}
	return -1, errors.New(fmt.Sprint(a, " not found in mapping"))
}

// AddMappedVertex adds a new Vertex with a mapped ID (or returns the index if
// ID already exists).
func (that *Graph) AddMappedVertex(ID string) int {
	if !that.usingMap || that.mapping == nil {
		that.usingMap = true
		that.mapping = map[string]int{}
		that.highestMapIndex = 0
	}
	if i, ok := that.mapping[ID]; ok {
		return i
	}
	i := that.highestMapIndex
	that.highestMapIndex++
	that.mapping[ID] = i
	return that.AddVertex(i).ID
}

// AddMappedArc adds a new Arc from Source to Destination, for when verticies are
// referenced by strings.
func (that *Graph) AddMappedArc(Source, Destination string, Distance int64) error {
	return that.AddArc(that.AddMappedVertex(Source), that.AddMappedVertex(Destination), Distance)
}

// AddArc is the default method for adding an arc from a Source Vertex to a
// Destination Vertex
func (that *Graph) AddArc(Source, Destination int, Distance int64) error {
	if len(that.Vertexes) <= Source || len(that.Vertexes) <= Destination {
		return ErrNodeNotFound
	}
	that.Vertexes[Source].AddArc(Destination, Distance)
	return nil
}

// RemoveArc removes and arc from the Source vertex to the Destination vertex
// fails if either vertex doesn't exist, but will succeed if destination is
// not an arc of Source (as a nop)
func (that *Graph) RemoveArc(Source, Destination int) error {
	if len(that.Vertexes) <= Source || len(that.Vertexes) <= Destination {
		return ErrNodeNotFound
	}
	that.Vertexes[Source].RemoveArc(Destination)
	return nil
}
