package heuristic

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

func H(r hyperrectangle.R) float64 { return hyperrectangle.SA(r) }
