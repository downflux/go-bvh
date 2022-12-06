package node

import (
	"fmt"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

func New(k vector.D, n int, tolerance float64, x ID) *N {
	if tolerance < 1 {
		panic(fmt.Sprintf("cannot set expansion tolerance to be less than resultant AABB"))
	}

	return &N{
		id: x,

		k:         k,
		n:         n,
		tolerance: tolerance,

		aabbCache: hyperrectangle.New(
			vector.V(make([]float64, k)),
			vector.V(make([]float64, k)),
		).M(),
		children: make(map[ID]*N, n),
	}
}

func Free(n *N) {
	// AABB will set by subsequent update calls. Skip zeroing the bounding
	// box.

	for x := range n.children {
		delete(n.children, x)
	}

	n.isLeaf = false
}
