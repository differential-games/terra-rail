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

func (m *Map) Fill(n *noise.Value) {
	for x := 0; x < m.Width; x++ {
		px := float64(x)*m.InvMaxScale
		for y := 0; y < m.Height; y++ {
			py := float64(y)*m.InvMaxScale
			m.Elevation[x*m.Width+y] = n.Cubic(px, py)
		}
	}
}
