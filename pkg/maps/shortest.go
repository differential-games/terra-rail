package maps

import (
	"container/heap"
	"fmt"
	"math"
)

type Origin int8

var (
	None  Origin = 0
	Start Origin = 1
	End   Origin = 2
	Both         = Start | End
)

// node represents a visited location on a Map.
type node struct {
	// idx is the id of tis node in the Map.
	idx int
	// previous is the index of the node which gets to this node fastest.
	// -1 indicates this is the starting node.
	previous int
	// origin is whether this node was first reached from the Start or End.
	origin Origin
	// dist is the cost to get to node from a starting point.
	dist float64
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
func Shortest(m *Map, from, to int) ([]int, []int) {
	visited := make([]Origin, m.Width*m.Height)
	paths := make([]node, m.Width*m.Height)

	// Initialize the node priority queue.
	toVisit := &nodeQueue{}
	heap.Init(toVisit)
	heap.Push(toVisit, node{
		idx:      from,
		origin:   Start,
		previous: -1,
		dist:     0.0,
	})
	heap.Push(toVisit, node{
		idx:      to,
		origin:   End,
		previous: -1,
		dist:     0.0,
	})

	end1 := 0
	end2 := 0
	for toVisit.Len() > 0 {
		// get the next node.
		cur := heap.Pop(toVisit).(node)

		prevOrigin := visited[cur.idx]
		if prevOrigin == cur.origin {
			// We've already visited this from this Origin in a faster way.
			continue
		}

		if cur.origin|prevOrigin == Both {
			// We've finally met.
			end1 = cur.previous
			end2 = cur.idx
			break
		}
		paths[cur.idx] = cur

		// Mark this node as visited and add it to known shortest paths.
		visited[cur.idx] = cur.origin

		curElevation := m.Elevation[cur.idx]
		// Get the nodes neighboring this one.
		neighbors := m.Neighbors(cur.idx)
		for _, n := range neighbors {
			if visited[n.idx] == cur.origin {
				// We've already visited this neighbor from this Origin, so
				// nothing to do.
				continue
			}

			dElevation := math.Abs(curElevation-m.Elevation[n.idx]) / n.dist
			// The distance metric. Heavily weights against large changes in
			// elevation.
			toN := n.dist + 1000000*dElevation*dElevation

			//slope := math.Abs(curElevation - m.Elevation[n.idx])/n.dist
			// The distance metric. Heavily weights against large changes in
			// elevation.
			//n.dist += 00000*slope*slope

			// Add this neighbor to the priority queue.
			heap.Push(toVisit, node{
				idx:      n.idx,
				origin:   cur.origin,
				previous: cur.idx,
				dist:     cur.dist + toN,
			})
		}
	}

	// Recreate the path from "from" to "to" by starting at "to" and
	// backtracking.
	var result1 []int
	for cur := end1; cur != -1; cur = paths[cur].previous {
		result1 = append(result1, paths[cur].idx)
	}
	var result2 []int
	for cur := end2; cur != -1; cur = paths[cur].previous {
		result2 = append(result2, paths[cur].idx)
	}
	fmt.Println(end1, end2)
	return result1, result2
}
