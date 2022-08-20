package heuristic

import (
	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
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

// Inherited calculates the additional heuristic cost of inserting the input
// object into a tree at the current node. If an additional node is inserted as
// a child of the current node, then the chain of parents up to the root must
// have their bounds adjusted as well -- we will need to take that into
// consideration.
func Inherited(nodes allocation.C[*node.N], n *node.N, bound hyperrectangle.R) H {
	if n == nil {
		return 0
	}

	h := H(0.0)
	// This naturally means leaf nodes have no cost, as we are only
	// considering parent nodes. In the case we have multiple points in a
	// single leaf node, we may need to revise this.
	for m := node.Parent(nodes, n); m != nil; m = node.Parent(nodes, m) {
		h += Heuristic(bhr.Union(bound, m.Bound())) - Heuristic(bound)
	}
	return h
}

// Direct calculates the heuristic cost of creating a new node containing the
// input object at the current node. This new node will be the parent of the
// current node, and will therefore have a cost determined by the current node
// and the object leaf (of size bound).
func Direct(n *node.N, bound hyperrectangle.R) H {
	return Heuristic(bhr.Union(bound, n.Bound()))
}
