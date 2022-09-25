package candidate

import (
	"math"

	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

type I interface {
	Height() uint
	AABB() hyperrectangle.R
}

// P is a pseudo-node used for prospective changes.
type P struct {
	L *node.N
	R *node.N

	Buf hyperrectangle.R
}

func (p P) Height() uint { return uint(math.Max(float64(p.L.Height()), float64(p.R.Height()))) + 1 }
func (p P) AABB() hyperrectangle.R {
	bhr.UnionBuf(p.L.AABB(), p.R.AABB(), p.Buf)
	return p.Buf
}

// C is a candidate for swapping nodes in the subtree. The B and C nodes here
// may be pseudo-nodes -- that is, what a node will look like assuming the swap
// occurs. We use these pseudo-nodes to check for the optimality of the swap;
// candidates which unbalance the tree or decreases the quality of the SAH must
// not be used.
//
// Here, the local subtree root A is implicit, and the nodes B and C are the
// left and right children of A (after the swap).
type C struct {
	B I
	C I

	Src    *node.N
	Target *node.N
}

// Simulate the B -> F rotation. In this case, C's  height (and AABB) may
// change, as its new children will be B and G. We need to ensure A is still
// balanced to fulfill the AVL guarantees -- that is,
//
//	|Height(A.L) - Height(A.R)| < 2
func (c C) Balanced() bool { return math.Abs(float64(c.B.Height())-float64(c.C.Height())) < 2 }
