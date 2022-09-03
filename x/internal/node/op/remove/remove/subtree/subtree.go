package subtree

import (
	"github.com/downflux/go-bvh/x/internal/node"
)

// T is a local subtree representation around an input "local" leaf G.
//
//	  A
//	 / \
//	B   C
//	   / \
//	  F   G
type T struct {
	A, B, C, F, G *node.N
}

func New(n *node.N) *T {
	t := &T{
		G: n,
	}

	if !n.IsRoot() {
		t.C = n.Parent()

		if t.C.Left() == n {
			t.F = t.C.Right()
		} else {
			t.F = t.C.Left()
		}

		if !t.C.IsRoot() {
			t.A = t.C.Parent()

			if t.A.Left() == t.C {
				t.B = t.A.Right()
			} else {
				t.B = t.A.Left()
			}
		}
	}
	return t
}
