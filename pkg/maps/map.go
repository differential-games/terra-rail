package maps

import (
	"math"

	"github.com/differential-games/hyper-terrain/pkg/noise"
)

type Map struct {
	Width, Height int

	InvMaxScale float64

	// Elevation is the distance above sea level.
	Elevation []float64
}

func NewMap(width, height int) *Map {
	return &Map{
		Width:  width,
		Height: height,

		InvMaxScale: 1.0 / 400.0,
		Elevation:   make([]float64, width*height),
	}
}

func (m *Map) Fill(n *noise.Fractal) {
	min := math.MaxFloat64
	max := 0.0
	for x := 0; x < m.Width; x++ {
		px := float64(x) * m.InvMaxScale
		for y := 0; y < m.Height; y++ {
			py := float64(y) * m.InvMaxScale
			elevation := n.Cubic(px, py)
			m.Elevation[x*m.Height+y] = elevation
			min = math.Min(min, elevation)
			max = math.Max(max, elevation)
		}
	}

	for i, e := range m.Elevation {
		m.Elevation[i] = (e - min) / (max - min)
	}
}

const sqrt2 = 1.4142135623730950488016887242096980785696718753769480731766797379

func (m *Map) IndexOf(x, y int) int {
	return x*m.Height+y
}

func (m *Map) Neighbors2(idx int) []int {
	x := idx / m.Height
	y := idx % m.Height

	var result []int
	if x > 0 {
		result = append(result, idx-m.Height)
	}
	if x > 0 && y > 0 {
		result = append(result, idx-m.Height-1)
	}
	if y > 0 {
		result = append(result, idx-1)
	}
	if x < (m.Width - 1) && y > 0 {
		result = append(result, idx+m.Height-1)
	}
	if x < (m.Width - 1) {
		result = append(result, idx+m.Height)
	}
	if x < (m.Width - 1) && y < (m.Height - 1) {
		result = append(result, idx+m.Height+1)
	}
	if y < (m.Height - 1) {
		result = append(result, idx+1)
	}
	if x > 0 && y < (m.Height - 1) {
		result = append(result, idx-m.Height+1)
	}
	return result
}

func (m *Map) Neighbors(idx int) []node {
	x := idx / m.Height
	y := idx % m.Height

	var result []node
	if x > 0 {
		result = append(result, node{
			idx: idx-m.Height,
			previous: idx,
			dist: 1.0,
		})
	}
	if x > 0 && y > 0 {
		result = append(result, node{
			idx: idx-m.Height-1,
			previous: idx,
			dist: sqrt2,
		})
	}
	if y > 0 {
		result = append(result, node{
			idx: idx-1,
			previous: idx,
			dist: 1.0,
		})
	}
	if x < (m.Width - 1) && y > 0 {
		result = append(result, node{
			idx: idx+m.Height-1,
			previous: idx,
			dist: sqrt2,
		})
	}
	if x < (m.Width - 1) {
		result = append(result, node{
			idx: idx+m.Height,
			previous: idx,
			dist: 1.0,
		})
	}
	if x < (m.Width - 1) && y < (m.Height - 1) {
		result = append(result, node{
			idx: idx+m.Height+1,
			previous: idx,
			dist: sqrt2,
		})
	}
	if y < (m.Height - 1) {
		result = append(result, node{
			idx: idx+1,
			previous: idx,
			dist: 1.0,
		})
	}
	if x > 0 && y < (m.Height - 1) {
		result = append(result, node{
			idx: idx-m.Height+1,
			previous: idx,
			dist: sqrt2,
		})
	}
	return result
}
