package remove

import (
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/remove/remove/subtree"
)

// Execute will remove an entire subtree. The returned node is the sibling node
// for the input node.
func Execute(n *node.N) *node.N {
	if n == nil {
		panic("cannot remove a nil node")
	}

	t := subtree.New(n)

	if t.A != nil {
		if t.A.Left() == t.C {
			t.A.SetLeft(t.F)
		} else {
			t.A.SetRight(t.F)
		}

		t.A.InvalidateAABBCache()
	}
	if t.C != nil {
		n.Cache().Delete(t.C.ID())
	}
	if t.F != nil {
		t.F.SetParent(t.A)
	}
	// We are deleting internal references to the local root -- children of
	// this node will evaluate t.C = nil.
	n.Cache().Delete(n.ID())

	if !n.IsLeaf() {
		// N.B.: We cannot execute this in parallel because the cache is
		// being mutated. We can add a lock, but in reality we are going
		// to remove only a leaf node, so there will be a performance
		// penalty.
		Execute(n.Left())
		Execute(n.Right())
	}

	// We now expect the local subtree to be of the format
	//
	//   A
	//  / \
	// B   F
	//
	// Where F may be the new root.
	return t.F
}
