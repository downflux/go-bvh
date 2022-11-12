package bvh

import (
	"github.com/downflux/go-bvh/x/internal/cache/node"
)

// swap will switch two nodes. We are assuming x and y are not ancestors of one
// another, e.g. neither x nor y is the root node.
func swap(x node.N, y node.N) {
	p, q := x.Parent(), y.Parent()
	b, c := p.Branch(x.ID()), q.Branch(y.ID())

	p.SetChild(b, y.ID())
	q.SetChild(c, x.ID())

	x.SetParent(q.ID())
	y.SetParent(p.ID())

	// N.B.: We are not setting the AABB nor the height of the parent nodes
	// p and q here.
}
