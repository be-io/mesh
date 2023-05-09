package graph

// BestPaths contains the list of best solutions
type BestPaths []BestPath

// ShortestAll calculates all of the shortest paths from src to dest
func (that *Graph) ShortestAll(src, dest int) (BestPaths, error) {
	return that.evaluateAll(src, dest, true)
}

// LongestAll calculates all the longest paths from src to dest
func (that *Graph) LongestAll(src, dest int) (BestPaths, error) {
	return that.evaluateAll(src, dest, false)
}

func (that *Graph) evaluateAll(src, dest int, shortest bool) (BestPaths, error) {
	if that.running {
		return BestPaths{}, ErrAlreadyCalculating
	}
	that.running = true
	defer func() { that.running = false }()
	//Setup graph
	that.setup(shortest, src, -1)
	return that.postSetupEvaluateAll(src, dest, shortest)
}

func (that *Graph) postSetupEvaluateAll(src, dest int, shortest bool) (BestPaths, error) {
	var current *Vertex
	oldCurrent := -1
	for that.visiting.Len() > 0 {
		//Visit the current lowest distanced Vertex
		current = that.visiting.PopOrdered()
		if oldCurrent == current.ID {
			continue
		}
		oldCurrent = current.ID
		//If the current distance is already worse than the best try another Vertex
		if shortest && current.distance > that.best {
			continue
		}
		for v, dist := range current.arcs {
			//If the arc has better access, than the current best, update the Vertex being touched
			if (shortest && current.distance+dist < that.Vertexes[v].distance) ||
				(!shortest && current.distance+dist > that.Vertexes[v].distance) ||
				(current.distance+dist == that.Vertexes[v].distance && !that.Vertexes[v].containsBest(current.ID)) {
				//if that.Vertexes[v].bestVertex == current.ID && that.Vertexes[v].ID != dest {
				if current.containsBest(v) && that.Vertexes[v].ID != dest {
					//also only do this if we aren't checkout out the best distance again
					//This seems familiar 8^)
					return BestPaths{}, newErrLoop(current.ID, v)
				}
				if current.distance+dist == that.Vertexes[v].distance {
					//At this point we know it's not in the list due to initial check
					that.Vertexes[v].bestVerticies = append(that.Vertexes[v].bestVerticies, current.ID)
				} else {
					that.Vertexes[v].distance = current.distance + dist
					that.Vertexes[v].bestVerticies = []int{current.ID}
				}
				if v == dest {
					that.visitedDest = true
					that.best = current.distance + dist
					continue
					//If this is the destination update best, so we can stop looking at
					// useless Vertexes
				}
				//Push this updated Vertex into the list to be evaluated, pushes in
				// sorted form
				that.visiting.PushOrdered(&that.Vertexes[v])
			}
		}
	}
	if !that.visitedDest {
		return BestPaths{}, ErrNoPath
	}
	return that.bestPaths(src, dest), nil
}

func (that *Graph) bestPaths(src, dest int) BestPaths {
	paths := that.visitPath(src, dest, dest)
	best := BestPaths{}
	for indexPaths := range paths {
		for i, j := 0, len(paths[indexPaths])-1; i < j; i, j = i+1, j-1 {
			paths[indexPaths][i], paths[indexPaths][j] = paths[indexPaths][j], paths[indexPaths][i]
		}
		best = append(best, BestPath{that.Vertexes[dest].distance, paths[indexPaths]})
	}

	return best
}

func (that *Graph) visitPath(src, dest, currentNode int) [][]int {
	if currentNode == src {
		return [][]int{
			{currentNode},
		}
	}
	paths := [][]int{}
	for _, vertex := range that.Vertexes[currentNode].bestVerticies {
		sps := that.visitPath(src, dest, vertex)
		for i := range sps {
			paths = append(paths, append([]int{currentNode}, sps[i]...))
		}
	}
	return paths
}
