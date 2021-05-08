package maps

import (
	"github.com/differential-games/hyper-terrain/pkg/noise"
)

type Map struct {
	Width, Height int

	InvMaxScale float64

	// Elevation is the distance above sea level.
	Elevation []float64
}

func NewMap(width, height int) Map {
	return Map{
		Width: width,
		Height: height,

		InvMaxScale: 1/400,
		Elevation: make([]float64, width*height),
	}
}

func (m *Map) Fill(n *noise.Fractal) {
	for x := 0; x < m.Width; x++ {
		px := float64(x)*m.InvMaxScale
		for y := 0; y < m.Height; y++ {
			py := float64(y)*m.InvMaxScale
			m.Elevation[x*m.Height+y] = n.Cubic(px, py)
		}
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
