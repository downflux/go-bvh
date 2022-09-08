package heuristic

import (
	"math"

	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

func H(r hyperrectangle.R) float64 { return bhr.V(r) }

// G is an experimental heuristic describing the rough surface area heuristic of
// a given hyperrectangle.
func G(r hyperrectangle.R) float64 { return math.Pow(bhr.V(r), 1.0/float64(r.Min().Dimension())) }
