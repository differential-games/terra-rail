package maps

import (
	"container/heap"
	"fmt"
	"math"
)

type dist struct {
	idx      int
	previous int
	dist     float64
}

type distQueue []dist

func (d *distQueue) Len() int {
	return len(*d)
}

func (d *distQueue) Less(i, j int) bool {
	return (*d)[i].dist < (*d)[j].dist
}

func (d *distQueue) Swap(i, j int) {
	(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
}

func (d *distQueue) Push(x interface{}) {
	*d = append(*d, x.(dist))
}

func (d *distQueue) Pop() interface{} {
	old := *d
	n := len(old)
	x := old[n-1]
	*d = old[0 : n-1]
	return x
}

var _ heap.Interface = &distQueue{}

func Shortest(m *Map, from, to int) []int {
	visited := make([]bool, m.Width*m.Height)
	paths := make([]dist, m.Width*m.Height)

	toVisit := &distQueue{}
	heap.Init(toVisit)
	heap.Push(toVisit, dist{
		idx:      from,
		previous: -1,
		dist:     0.0,
	})

	i := 0
	for toVisit.Len() > 0 {
		next := heap.Pop(toVisit).(dist)
		if visited[next.idx] {
			continue
		}

		visited[next.idx] = true
		paths[next.idx] = next
		if next.idx == to {
			break
		}

		curElevation := m.Elevation[next.idx]
		neighbors := m.Neighbors(next.idx)
		for _, n := range neighbors {
			dElevation := math.Abs(curElevation - m.Elevation[n])
			toN := 1 + 100000000*dElevation*dElevation

			if visited[n] {
				continue
			}

			heap.Push(toVisit, dist{
				idx:      n,
				previous: next.idx,
				dist:     next.dist + toN,
			})
		}
		//fmt.Println(toVisit.Len())
		i++
		if i > 10000000 {
			fmt.Println(next.dist)
			panic("WHAT")
		}
	}

	var result []int
	for cur := to; cur != from; {
		result = append(result, paths[cur].idx)
		cur = paths[cur].previous
	}
	result = append(result, paths[from].idx)
	return result
}
