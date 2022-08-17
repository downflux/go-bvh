package util

import (
	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
)

type N struct {
	ID point.ID

	Left  *N
	Right *N

	Bound hyperrectangle.R
}

func Construct(c allocation.C[*node.N], n *node.N) *N {
	if n == nil {
		return nil
	}

	m := &N{
		ID:    n.ID(),
		Bound: n.Bound(),
	}

	m.Left = Construct(c, node.Left(c, n))
	m.Right = Construct(c, node.Right(c, n))

	return m
}

// Equal is a test-only function which determines the equality between two
// allocation objects. We consider allocations equal if the node relations are
// invariant under allocation IDs.
func Equal(
	a allocation.C[*node.N],
	r allocation.ID,
	b allocation.C[*node.N],
	s allocation.ID,
) bool {
	return cmp.Equal(
		Construct(a, a[r]),
		Construct(b, b[s]),
		cmp.AllowUnexported(
			hyperrectangle.R{},
		),
	)
}
