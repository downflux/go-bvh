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
	if n.Left() != nil {
		t.B = n.Left()
		t.D = n.Left().Left()
		t.E = n.Left().Right()
	}
	if n.Right() != nil {
		t.C = n.Right()
		t.F = n.Right().Left()
		t.G = n.Right().Right()
	}
	return t
}
