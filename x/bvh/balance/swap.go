package balance

import (
	"github.com/downflux/go-bvh/x/internal/cache/node"
)

// swap will switch two nodes. We are assuming x and y are not ancestors of one
// another, e.g. neither x nor y is the root node.
//
// N.B.: This will leave nodes x and y (and their parents) in an invalid state.
// The caller is responsible for maintaining the consistency of the tree by
// manually updating the height and AABB of all affected nodes.
func swap(x node.N, y node.N) {
	p, q := x.Parent(), y.Parent()
	b, c := p.Branch(x.ID()), q.Branch(y.ID())

	p.SetChild(b, y.ID())
	q.SetChild(c, x.ID())

	x.SetParent(q.ID())
	y.SetParent(p.ID())
}
