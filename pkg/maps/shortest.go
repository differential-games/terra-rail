package maps

import "container/heap"

type dist struct {
	idx      int
	previous int
	dist     int
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
	head := (*d)[0]
	*d = (*d)[1:]
	return head
}

var _ heap.Interface = &distQueue{}

func Shortest(m *Map, from, to int) {
	visited := make([]bool, m.Width*m.Height)
	paths := make([]dist, m.Width*m.Height)

	toVisit := &distQueue{}
	heap.Push(toVisit, dist{
		idx: from,
		previous: -1,
		dist: 0.0,
	})

	for next := heap.Pop(toVisit).(dist); toVisit.Len() > 0; next = heap.Pop(toVisit).(dist) {
		visited[next.idx] = true

		heap.Push(toVisit, dist{
			idx:
		})
	}

}
