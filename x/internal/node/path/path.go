package path

import (
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/node"
)

// P is a path iterator struct representing the current node.
type P struct {
	// P is the parent of the current node.
	P *node.N

	// B is the branch of the parent which navigates to the current node.
	B node.Branch
}

// N returns the actual node tracked by the current iterator.
func (p P) N(c *cache.C[*node.N]) *node.N { return p.P.Child(c, p.B) }

// Next returns the next iterator.
func (p P) Next(c *cache.C[*node.N]) P {
	n := p.N(c)
	if n.IsRoot() {
		panic("iterating past end of sequence")
	}

	q := n.Parent(c)
	return P{
		P: q,
		B: q.Branch(n.ID()),
	}
}
