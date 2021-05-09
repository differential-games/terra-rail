package maps

import (
	"container/heap"
	"math"
)

// node represents a visited location on a Map.
type node struct {
	// idx is the id of tis node in the Map.
	idx      int
	// previous is the index of the node which gets to this node fastest.
	// -1 indicates this is the starting node.
	previous int
	// dist is the cost to get to node from a starting point.
	dist     float64
}

// nodeQueue is a priority queue of nodes to visit.
type nodeQueue []node

// Len implements heap.Interface.
func (d *nodeQueue) Len() int {
	return len(*d)
}

// Less implements heap.Interface.
func (d *nodeQueue) Less(i, j int) bool {
	return (*d)[i].dist < (*d)[j].dist
}

// Swap implements heap.Interface.
func (d *nodeQueue) Swap(i, j int) {
	(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
}

// Push implements heap.Interface.
func (d *nodeQueue) Push(x interface{}) {
	*d = append(*d, x.(node))
}

// Pop implements heap.Interface.
func (d *nodeQueue) Pop() interface{} {
	old := *d
	n := len(old)
	x := old[n-1]
	*d = old[0 : n-1]
	return x
}

var _ heap.Interface = &nodeQueue{}

// Shortest calculates the "shortest" path between two points.
// An implementation of Dijkstra's algorithm.
//
// TODO: make the distance metric configurable.
// TODO: use A* to make more efficient.
func Shortest(m *Map, from, to int) []int {
	visited := make([]bool, m.Width*m.Height)
	paths := make([]node, m.Width*m.Height)

	// Initialize the node priority queue.
	toVisit := &nodeQueue{}
	heap.Init(toVisit)
	heap.Push(toVisit, node{
		idx:      from,
		previous: -1,
		dist:     0.0,
	})

	for toVisit.Len() > 0 {
		// get the next node.
		next := heap.Pop(toVisit).(node)
		if visited[next.idx] {
			// We got to this node by a shorter path, so discard.
			continue
		}

		// Mark this node as visited and add it to known shortest paths.
		visited[next.idx] = true
		paths[next.idx] = next
		if next.idx == to {
			// We've reached the target node.
			break
		}

		curElevation := m.Elevation[next.idx]
		// Get the nodes neighboring this one.
		neighbors := m.Neighbors(next.idx)
		for _, n := range neighbors {
			dElevation := math.Abs(curElevation - m.Elevation[n])
			// The distance metric. Heavily weights against large changes in
			// elevation.
			toN := 1 + 100000000*dElevation*dElevation

			if visited[n] {
				// We've already visited this neighbor, so nothing to do.
				continue
			}

			// Add this neighbor to the priority queue.
			heap.Push(toVisit, node{
				idx:      n,
				previous: next.idx,
				dist:     next.dist + toN,
			})
		}
	}

	// Recreate the path from "from" to "to" by starting at "to" and
	// backtracking.
	var result []int
	for cur := to; cur != from; {
		result = append(result, paths[cur].idx)
		cur = paths[cur].previous
	}
	result = append(result, paths[from].idx)
	return result
}
