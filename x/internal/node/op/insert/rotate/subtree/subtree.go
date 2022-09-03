package subtree

import (
	"github.com/downflux/go-bvh/x/internal/node"
)

type T struct {
	A, B, C, D, E, F, G *node.N
}

func New(n *node.N) *T {
	t := &T{
		A: n,
	}

	if !n.IsLeaf() {
		t.B = n.Left()
		t.C = n.Right()

		if !n.Left().IsLeaf() {
			t.D = n.Left().Left()
			t.E = n.Left().Right()
		}
		if !n.Right().IsLeaf() {
			t.F = n.Right().Left()
			t.G = n.Right().Right()
		}
	}
	return t
}
