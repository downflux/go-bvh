package point

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type P interface {
	B() hyperrectangle.R
	Heuristic() float64
}
