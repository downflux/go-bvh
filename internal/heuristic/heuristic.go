package heuristic

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type H float64

func Heuristic(r hyperrectangle.R) H {
	h := 1.0
	d := r.D()
	for i := vector.D(0); i < d.Dimension(); i++ {
		h *= d.X(i)
	}
	return H(h)
}
