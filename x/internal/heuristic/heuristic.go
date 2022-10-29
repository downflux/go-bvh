package heuristic

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

var (
	// H is intended to be the SAH as defined in canonical literature.
	// Experimentally, for BroadPhase, we have found the volume produces a
	// BVH with faster Insert and BroadPhase operations.
	H = hyperrectangle.V
)
