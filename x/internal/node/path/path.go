package path

import (
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
func (p P) N() *node.N { return p.P.Child(p.B) }

// Next returns the next iterator.
func (p P) Next() P {
	n := p.N()
	if n.IsRoot() {
		panic("iterating past end of sequence")
	}

	q := n.Parent()
	return P{
		P: q,
		B: q.Branch(n.ID()),
	}
}
