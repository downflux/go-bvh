package subtree

import (
	"github.com/downflux/go-bvh/internal/node"
)

// T is a local subtree given the node G.
//
//	  A
//	 / \
//	B   C
//	   / \
//	  F   G
type T struct {
	A, B, C, F, G *node.N
}

// New generates a subtree about the G-node. The returned value is nil if no
// subtree exists.
func New(n *node.N) *T {
	if n.IsRoot() || n.Parent().IsRoot() {
		return nil
	}

	return &T{
		A: n.Parent().Parent(),
		B: map[bool]*node.N{
			true:  n.Parent().Parent().Right(),
			false: n.Parent().Parent().Left(),
		}[n.Parent() == n.Parent().Parent().Left()],
		C: n.Parent(),
		F: map[bool]*node.N{
			true:  n.Parent().Right(),
			false: n.Parent().Left(),
		}[n == n.Parent().Left()],
		G: n,
	}
}
