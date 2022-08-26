package hyperrectangle

import (
	"math"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

// Contains checks if the input rectangle r fully encloses s.
func Contains(r hyperrectangle.R, s hyperrectangle.R) bool {
	if r.Min().Dimension() != s.Min().Dimension() {
		panic("mismatching vector dimensions")
	}

	for i := vector.D(0); i < r.Min().Dimension(); i++ {
		if s.Min().X(i) <= r.Min().X(i) || s.Max().X(i) >= r.Max().X(i) {
			return false
		}
	}
	return true
}

func Collision(r hyperrectangle.R, s hyperrectangle.R) bool {
	if r.Min().Dimension() != s.Min().Dimension() {
		panic("mismatching vector dimensions")
	}

	for i := vector.D(0); i < r.Min().Dimension(); i++ {
		l := math.Max(r.Min().X(i), s.Min().X(i))
		u := math.Min(r.Max().X(i), s.Max().X(i))
		if l < u {
			return false
		}
	}
	return true
}

func Union(r hyperrectangle.R, s hyperrectangle.R) hyperrectangle.R {
	if r.Min().Dimension() != s.Min().Dimension() {
		panic("mismatching vector dimensions")
	}

	min := make([]float64, r.Min().Dimension())
	max := make([]float64, r.Min().Dimension())

	for i := vector.D(0); i < r.Min().Dimension(); i++ {
		min[i] = math.Min(r.Min().X(i), s.Min().X(i))
		max[i] = math.Max(r.Max().X(i), s.Max().X(i))
	}

	return *hyperrectangle.New(min, max)
}
