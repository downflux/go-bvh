package remove

import (
	"fmt"

	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/allocation/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/refit"
	"github.com/downflux/go-bvh/internal/node/rotate"
)

// Execute removes a node from the tree and rebalances the resulting structure.
//
// We are guaranteed the node to be deleted nid exists in the allocation.
//
// Delete returns the new root.
func Execute(nodes allocation.C[*node.N], nid id.ID) id.ID {
	n := nodes[nid]
	if !node.Leaf(nodes, n) {
		panic(fmt.Sprintf("cannot directly remove an interior node: %v", nid))
	}

	p := node.Parent(nodes, n)

	// Handle case where n is the root.
	if p == nil {
		nodes.Remove(nid)
		return 0
	}

	var s *node.N
	if node.Left(nodes, p) == n {
		s = node.Right(nodes, p)
	} else {
		s = node.Left(nodes, p)
	}

	a := node.Parent(nodes, p)

	// Handle the case where the sibling is not root.
	if a != nil {
		if node.Left(nodes, a) == p {
			a.SetLeft(s.Index())
		} else {
			a.SetRight(s.Index())
		}

		s.SetParent(a.Index())
	} else {
		s.SetParent(0)
	}

	nodes.Remove(p.Index())
	nodes.Remove(n.Index())

	refit.Execute(nodes, s.Index())
	rotate.Execute(nodes, s.Index())

	r := node.Root(nodes, s)
	if r == nil {
		return 0
	}
	return r.Index()
}
