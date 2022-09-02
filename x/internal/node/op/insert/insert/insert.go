// Package defines an operation which adds a new subtree into an existing tree.
package insert

// Execute inserts the node m as a sibling of n. The returned node is the new
// root of the tree.
func Execute(n *node.N, m *node.N) *node.N {
	if n == nil || m == nil {
		panic("cannot insert an empty node into a (possibly) empty tree")
	}

	if n.Cache() != m.Cache() {
		panic("cannot insert nodes with mismatching lookup tables")
	}

	var aid nid.ID
	if !n.IsRoot() {
		aid = n.Parent().ID()
	}

	p := node.New(node.O{
		Cache: n.Cache(),

		Left: n.ID(),
		Right: m.ID(),
		Parent: aid,
	})

	m.SetParent(p)
	
}
