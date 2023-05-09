package graph

import (
	"errors"
	"fmt"
)

// Graph contains all the graph details
type Graph struct {
	Home            string
	best            int64
	visitedDest     bool
	Vertexes        []Vertex //slice of all Vertexes available
	visiting        dijkstraList
	mapping         map[string]int
	usingMap        bool
	highestMapIndex int
	running         bool
}

// NewGraph creates a new empty graph
func NewGraph(home string) *Graph {
	return &Graph{
		Home:    home,
		mapping: map[string]int{},
	}
}

// AddNewVertex adds a new vertex at the next available index
func (that *Graph) AddNewVertex() *Vertex {
	for i, v := range that.Vertexes {
		if i != v.ID {
			that.Vertexes[i] = Vertex{ID: i}
			return &that.Vertexes[i]
		}
	}
	return that.AddVertex(len(that.Vertexes))
}

// AddVertex adds a single vertex
func (that *Graph) AddVertex(ID int) *Vertex {
	that.AddVerticies(Vertex{ID: ID})
	return &that.Vertexes[ID]
}

// GetVertex gets the reference of the specified vertex. An error is thrown if
// there is no vertex with that index/ID.
func (that *Graph) GetVertex(ID int) (*Vertex, error) {
	if ID >= len(that.Vertexes) {
		return nil, errors.New("Vertex not found")
	}
	return &that.Vertexes[ID], nil
}

func (that Graph) validate() error {
	for _, v := range that.Vertexes {
		for a := range v.arcs {
			if a >= len(that.Vertexes) || (that.Vertexes[a].ID == 0 && a != 0) {
				return errors.New(fmt.Sprint("Graph validation error;", "Vertex ", a, " referenced in arcs by Vertex ", v.ID))
			}
		}
	}
	return nil
}

// SetDefaults sets the distance and best node to that specified
func (that *Graph) setDefaults(Distance int64, BestNode int) {
	for i := range that.Vertexes {
		that.Vertexes[i].bestVerticies = []int{BestNode}
		that.Vertexes[i].distance = Distance
	}
}
