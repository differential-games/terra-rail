package maps

import (
	"github.com/differential-games/hyper-terrain/pkg/noise"
	"math"
)

type Map struct {
	Width, Height int

	InvMaxScale float64

	// Elevation is the distance above sea level.
	Elevation []float64
}

func NewMap(width, height int) Map {
	return Map{
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

func (m *Map) Neighbors(idx int) []int {
	x := idx / m.Height
	y := idx % m.Height

	var result []int
	if y > 0 {
		result = append(result, idx-1)
	}
	if y < (m.Height - 1) {
		result = append(result, idx+1)
	}
	if x > 0 {
		result = append(result, idx-m.Height)
	}
	if x < (m.Width - 1) {
		result = append(result, idx+m.Height)
	}
	return result
}
