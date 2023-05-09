package graph

import (
	"math"
)

// BestPath contains the solution of the most optimal path
type BestPath struct {
	Distance int64
	Path     []int
}

// Shortest calculates the shortest path from src to dest
func (that *Graph) Shortest(src, dest int) (BestPath, error) {
	return that.evaluate(src, dest, true)
}

// Longest calculates the longest path from src to dest
func (that *Graph) Longest(src, dest int) (BestPath, error) {
	return that.evaluate(src, dest, false)
}

func (that *Graph) setup(shortest bool, src int, list int) {
	//-1 auto list
	//Get a new list regardless
	if list >= 0 {
		that.forceList(list)
	} else if shortest {
		that.forceList(-1)
	} else {
		that.forceList(-2)
	}
	//Reset state
	that.visitedDest = false
	//Reset the best current value (worst so it gets overwritten)
	// and set the defaults *almost* as bad
	// set all best verticies to -1 (unused)
	if shortest {
		that.setDefaults(int64(math.MaxInt64)-2, -1)
		that.best = int64(math.MaxInt64)
	} else {
		that.setDefaults(int64(math.MinInt64)+2, -1)
		that.best = int64(math.MinInt64)
	}
	//Set the distance of initial vertex 0
	that.Vertexes[src].distance = 0
	//Add the source vertex to the list
	that.visiting.PushOrdered(&that.Vertexes[src])
}

func (that *Graph) forceList(i int) {
	//-2 long auto
	//-1 short auto
	//0 short pq
	//1 long pq
	//2 short ll
	//3 long ll
	switch i {
	case -2:
		if len(that.Vertexes) < 800 {
			that.forceList(2)
		} else {
			that.forceList(0)
		}
	case -1:
		if len(that.Vertexes) < 800 {
			that.forceList(3)
		} else {
			that.forceList(1)
		}
	case 0:
		that.visiting = priorityQueueNewShort()
	case 1:
		that.visiting = priorityQueueNewLong()
	case 2:
		that.visiting = linkedListNewShort()
	case 3:
		that.visiting = linkedListNewLong()
	default:
		panic(i)
	}
}

func (that *Graph) bestPath(src, dest int) BestPath {
	var path []int
	for c := that.Vertexes[dest]; c.ID != src; c = that.Vertexes[c.bestVerticies[0]] {
		path = append(path, c.ID)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return BestPath{that.Vertexes[dest].distance, path}
}

func (that *Graph) evaluate(src, dest int, shortest bool) (BestPath, error) {
	if that.running {
		return BestPath{}, ErrAlreadyCalculating
	}
	that.running = true
	defer func() { that.running = false }()
	//Setup graph
	that.setup(shortest, src, -1)
	return that.postSetupEvaluate(src, dest, shortest)
}

func (that *Graph) postSetupEvaluate(src, dest int, shortest bool) (BestPath, error) {
	var current *Vertex
	oldCurrent := -1
	for that.visiting.Len() > 0 {
		//Visit the current lowest distanced Vertex
		//TODO WTF
		current = that.visiting.PopOrdered()
		if oldCurrent == current.ID {
			continue
		}
		oldCurrent = current.ID
		//If the current distance is already worse than the best try another Vertex
		if shortest && current.distance >= that.best {
			continue
		}
		for v, dist := range current.arcs {
			//If the arc has better access, than the current best, update the Vertex being touched
			if (shortest && current.distance+dist < that.Vertexes[v].distance) ||
				(!shortest && current.distance+dist > that.Vertexes[v].distance) {
				if current.bestVerticies[0] == v && that.Vertexes[v].ID != dest {
					//also only do this if we aren't checkout out the best distance again
					//This seems familiar 8^)
					return BestPath{}, newErrLoop(current.ID, v)
				}
				that.Vertexes[v].distance = current.distance + dist
				that.Vertexes[v].bestVerticies[0] = current.ID
				if v == dest {
					//If this is the destination update best, so we can stop looking at
					// useless Vertexes
					that.best = current.distance + dist
					that.visitedDest = true
					continue // Do not push if dest
				}
				//Push this updated Vertex into the list to be evaluated, pushes in
				// sorted form
				that.visiting.PushOrdered(&that.Vertexes[v])
			}
		}
	}
	return that.finally(src, dest)
}

func (that *Graph) finally(src, dest int) (BestPath, error) {
	if !that.visitedDest {
		return BestPath{}, ErrNoPath
	}
	return that.bestPath(src, dest), nil
}
