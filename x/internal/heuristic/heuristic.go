package heuristic

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/x/hyperrectangle"
)

func H(r hyperrectangle.R) float64 { return bhr.V(r) }
