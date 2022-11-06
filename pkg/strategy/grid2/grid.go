package grid2

import (
	"math"

	"github.com/c9s/bbgo/pkg/fixedpoint"
)

type Grid struct {
	UpperPrice fixedpoint.Value `json:"upperPrice"`
	LowerPrice fixedpoint.Value `json:"lowerPrice"`

	// Size is the number of total grids
	Size fixedpoint.Value `json:"size"`

	// Pins are the pinned grid prices, from low to high
	Pins []Pin `json:"pins"`

	pinsCache map[Pin]struct{} `json:"-"`
}

type Pin fixedpoint.Value

func calculateArithmeticPins(lower, upper, size, tickSize fixedpoint.Value) []Pin {
	var height = upper.Sub(lower)
	var spread = height.Div(size)

	var pins []Pin
	for p := lower; p.Compare(upper) <= 0; p = p.Add(spread) {
		// tickSize here = 0.01
		pp := math.Trunc(p.Float64()/tickSize.Float64()) * tickSize.Float64()
		pins = append(pins, Pin(fixedpoint.NewFromFloat(pp)))
	}

	return pins
}

func buildPinCache(pins []Pin) map[Pin]struct{} {
	cache := make(map[Pin]struct{}, len(pins))
	for _, pin := range pins {
		cache[pin] = struct{}{}
	}

	return cache
}

func NewGrid(lower, upper, size, tickSize fixedpoint.Value) *Grid {
	var pins = calculateArithmeticPins(lower, upper, size, tickSize)

	grid := &Grid{
		UpperPrice: upper,
		LowerPrice: lower,
		Size:       size,
		Pins:       pins,
		pinsCache:  buildPinCache(pins),
	}

	return grid
}

func (g *Grid) Height() fixedpoint.Value {
	return g.UpperPrice.Sub(g.LowerPrice)
}

// Spread returns the spread of each grid
func (g *Grid) Spread() fixedpoint.Value {
	return g.Height().Div(g.Size)
}

func (g *Grid) Above(price fixedpoint.Value) bool {
	return price.Compare(g.UpperPrice) > 0
}

func (g *Grid) Below(price fixedpoint.Value) bool {
	return price.Compare(g.LowerPrice) < 0
}

func (g *Grid) OutOfRange(price fixedpoint.Value) bool {
	return price.Compare(g.LowerPrice) < 0 || price.Compare(g.UpperPrice) > 0
}

func (g *Grid) updatePinsCache() {
	g.pinsCache = buildPinCache(g.Pins)
}

func (g *Grid) HasPin(pin Pin) (ok bool) {
	_, ok = g.pinsCache[pin]
	return ok
}

func (g *Grid) ExtendUpperPrice(upper fixedpoint.Value) (newPins []Pin) {
	g.UpperPrice = upper

	// since the grid is extended, the size should be updated as well
	spread := g.Spread()
	g.Size = (g.UpperPrice - g.LowerPrice).Div(spread).Floor()

	lastPinPrice := fixedpoint.Value(g.Pins[len(g.Pins)-1])
	for p := lastPinPrice.Add(spread); p <= g.UpperPrice; p = p.Add(spread) {
		newPins = append(newPins, Pin(p))
	}

	g.Pins = append(g.Pins, newPins...)
	g.updatePinsCache()
	return newPins
}

func (g *Grid) ExtendLowerPrice(lower fixedpoint.Value) (newPins []Pin) {
	spread := g.Spread()

	g.LowerPrice = lower

	height := g.Height()

	// since the grid is extended, the size should be updated as well
	g.Size = height.Div(spread).Floor()

	firstPinPrice := fixedpoint.Value(g.Pins[0])
	numToAdd := (firstPinPrice.Sub(g.LowerPrice)).Div(spread).Floor()
	if numToAdd == 0 {
		return newPins
	}

	for p := firstPinPrice.Sub(spread.Mul(numToAdd)); p.Compare(firstPinPrice) < 0; p = p.Add(spread) {
		newPins = append(newPins, Pin(p))
	}

	g.Pins = append(newPins, g.Pins...)
	g.updatePinsCache()
	return newPins
}
