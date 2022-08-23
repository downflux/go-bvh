package point

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type ID string

type RO interface {
	Bound() hyperrectangle.R
	ID() ID
}
