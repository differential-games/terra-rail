package maps

import (
	"container/heap"
	"math"
)

type Growth struct {
	Idx       int
	Elevation float64
	D         float64
}

type growthQueue []Growth

// Len implements heap.Interface.
func (d *growthQueue) Len() int {
	return len(*d)
}

// Less implements heap.Interface.
func (d *growthQueue) Less(i, j int) bool {
	return (*d)[i].D < (*d)[j].D
}

// Swap implements heap.Interface.
func (d *growthQueue) Swap(i, j int) {
	(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
}

// Push implements heap.Interface.
func (d *growthQueue) Push(x interface{}) {
	*d = append(*d, x.(Growth))
}

// Pop implements heap.Interface.
func (d *growthQueue) Pop() interface{} {
	old := *d
	n := len(old)
	x := old[n-1]
	*d = old[0 : n-1]
	return x
}

func GrowthAround(m *Map, idx int) chan Growth {
	result := make(chan Growth)
	visited := make([]bool, m.Width*m.Height)

	toVisit := &growthQueue{}
	heap.Init(toVisit)
	heap.Push(toVisit, Growth{
		Idx:       idx,
		Elevation: m.Elevation[idx],
		D:         0.0,
	})

	go func() {
		for toVisit.Len() > 0 {
			cur := heap.Pop(toVisit).(Growth)
			if visited[cur.Idx] {
				continue
			}
			visited[cur.Idx] = true

			for _, n := range m.Neighbors2(cur.Idx) {
				nElevation := m.Elevation[n]
				dElevation := math.Abs(cur.Elevation - nElevation)*10000

				base := 1.0
				if (cur.Idx+ n) % 2 == 1 {
					base = sqrt2
				}

				toN := base + dElevation*dElevation*dElevation*dElevation*dElevation

				heap.Push(toVisit, Growth{
					Idx:       n,
					Elevation: nElevation,
					D:         cur.D + toN,
				})
			}

			result <- cur
		}
		close(result)
	}()

	return result
}
