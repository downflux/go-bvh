package subtree

import (
	"github.com/downflux/go-bvh/internal/node"
)

// T is a local subtree representation around an input "local" root A.
//
//	 A
//	/ \
//
// B   C
//
//	 / \
//	F   G
//
// Where the following conditions are met --
//
//  1. C.Height() > B.Height() and
//  2. F.Height() >= G.Height()
//
// We will use this struct to transform the local tree into
//
//	 C
//	/ \
//
// F   A
//
//	 / \
//	B   G
//
// This is from the reference Box2D implementation (github.com/erincatto/box2d).
// Note that this is a different rotation strategy than the Catto 2019 slides.
type T struct {
	A *node.N
	B *node.N
	C *node.N
	F *node.N
	G *node.N
}

func New(n *node.N) *T {
	t := &T{
		A: n,
	}

	if n.IsLeaf() || n.Height() < 2 {
		return nil
	}

	if n.Right().Height() == n.Left().Height() {
		return nil
	} else if n.Right().Height() > n.Left().Height() {
		t.B = n.Left()
		t.C = n.Right()
	} else {
		t.B = n.Right()
		t.C = n.Left()
	}

	if t.C.Left().Height() > t.C.Right().Height() {
		t.F = t.C.Left()
		t.G = t.C.Right()
	} else {
		t.F = t.C.Right()
		t.G = t.C.Left()
	}
	return t
}
