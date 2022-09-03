// Package defines an operation which adds a new subtree into an existing tree.
package insert

import (
	"github.com/downflux/go-bvh/x/internal/node"

	nid "github.com/downflux/go-bvh/x/internal/node/id"
)

// Execute inserts the node m as a sibling of n. The returned node is the new
// root of the tree.
func Execute(n *node.N, m *node.N) *node.N {
	if n == nil || m == nil {
		panic("cannot insert an empty node into a (possibly) empty tree")
	}

	if !m.IsRoot() {
		panic("cannot insert an internal node into an existing tree")
	}

	if n.Cache() != m.Cache() {
		panic("cannot insert nodes with mismatching lookup tables")
	}

	p := n.Parent()
	var aid nid.ID
	if !n.IsRoot() {
		aid = p.ID()
	}

	q := node.New(node.O{
		Nodes: n.Cache(),

		Left:   m.ID(),
		Right:  n.ID(),
		Parent: aid,
	})

	if !n.IsRoot() {
		if p.Left() == n {
			p.SetLeft(q)
		} else {
			p.SetRight(q)
		}
	}

	m.SetParent(q)
	n.SetParent(q)

	return n.Root()
}
