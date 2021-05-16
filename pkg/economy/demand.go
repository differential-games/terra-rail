package economy

import (
	"github.com/differential-games/terra-rail/pkg/maps"
	"math"
)

type Demand struct {
	// Price is the market price of the good.
	Price []float64

	// Supply is the amount of the resource on the map.
	Supply []float64
	// Demand is the per-cycle consumption of the resource.
	Demand map[int]float64
	// MaxPrice is the maximum a buyer will ever pay for the resource.
	MaxPrice float64
}

func NewDemand(width, height, maxPrice float64) *Demand {
	return &Demand{
		Price: make([]float64, width*height),
		Supply: make([]float64, width*height),
		Demand: make(map[int]float64),
		MaxPrice: maxPrice,
	}
}

func (d *Demand) moveSupply(m *maps.Map) {
	// Up to 1 of Supply moves to an adjacent cell with a higher Price.
	for idx, s := range d.Supply {
		if s == 0 {
			continue
		}
		destination := idx
		// maxPrice is the highest market value for the resource in this or
		// adjacent cells.
		maxPrice := d.Price[idx]
		for _, n := range m.Neighbors2(idx) {
			nPrice := d.Price[n]
			if nPrice > maxPrice {
				destination = n
				maxPrice = nPrice
			}
		}
		if destination == idx {
			// The goods aren't moving anywhere.
			continue
		}

		dSupply := math.Min(1.0, d.Supply[idx])
		d.Supply[idx] -= dSupply
		d.Supply[destination] += dSupply
	}
}

func (d *Demand) consume() {
	// Assume the number of demand points is small compared to the size of the
	// map.
	for idx, want := range d.Demand {
		// Consume existing supply at the demand point.
		oldSupply := d.Supply[idx]
		// Negative newSupply indicates unmet demand.
		newSupply := oldSupply - want

		// Increase price by the unmet demand.
		// Decrease price by excess supply.
		d.Price[idx] = math.Max(d.MaxPrice, d.Price[idx]-newSupply)
		// Supply can't go lower than 0.
		d.Supply[idx] = math.Max(0.0, newSupply)
	}
}

func (d *Demand) updatePrice(m *maps.Map) {
	// We want each cell to point towards the neighbor which allows it to take
	// on the highest price. Effectively, it's finding the cheapest path to a
	// demander.
	newPrices := make([]float64, len(d.Price))
	for idx, elevation := range m.Elevation {
		// Default to targeting self if unable to offer a better price.
		maxPrice := 0.0
		if d.Demand[idx] > 0 {
			// This is a demander, so it has an intrinsic price.
			maxPrice = d.Price[idx]
		}
		for _, n := range m.Neighbors2(idx) {
			// deltaElevation is the change in elevation in the direction we
			// expect goods to flow.
			deltaElevation := m.Elevation[n] - elevation
			// moveCost is the cost to move from idx to n.
			moveCost := 1.0 + deltaElevation*deltaElevation

			toNPrice := d.Price[n] - moveCost
			// This neighbor allows us to charge a higher price.
			maxPrice = math.Max(maxPrice, toNPrice)
		}
		newPrices[idx] = maxPrice
	}
	d.Price = newPrices
}

func (d *Demand) Tick(m *maps.Map) {
	d.moveSupply(m)
	d.updatePrice(m)
	d.consume()
}
