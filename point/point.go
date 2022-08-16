package point

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type ID string

type P interface {
	Bound() hyperrectangle.R
	Heuristic() float64
	ID() ID
}
