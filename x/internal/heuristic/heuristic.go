package heuristic

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

var (
	// H is intended to be the SAH as defined in canonical literature.
	H = hyperrectangle.SA
)

// The briannoyama implementation uses the total edge-weight of a hyperrectangle
// as a ray collision intersection heuristic.
func BrianNoyama(r hyperrectangle.R) float64 {
	rmin, rmax := r.Min(), r.Max()

	h := 0.0
	for i := vector.D(0); i < rmin.Dimension(); i++ {
		h += rmax[i] - rmin[i]
	}
	return h
}
