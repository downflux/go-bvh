package point

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type P interface {
	Bound() hyperrectangle.R
	Heuristic() float64
}
