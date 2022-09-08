package swap

import (
	"github.com/downflux/go-bvh/internal/node"
)

// isAncestor returns true if the current node n is an ancestor of m.
func isAncestor(n *node.N, m *node.N) bool {
	if n == nil || m == nil {
		return false
	}
	if n == m {
		return true
	}
	if m.IsRoot() {
		return false
	}
	return isAncestor(n, m.Parent())
}

// Swap will move the input nodes n and m such that they and their subtrees
// belong to the other's parents.
func Execute(n *node.N, m *node.N) {
	if n == nil || m == nil {
		panic("cannot swap an empty node")
	}
	if n.Root() != m.Root() {
		panic("cannot swap nodes of disjoint trees")
	}
	if isAncestor(n, m) || isAncestor(m, n) {
		panic("cannot swap a node with its ancestor")
	}

	// Since the root node is an ancestor of all nodes, we can assume n and
	// m are not the root.
	p := n.Parent()
	q := m.Parent()

	if p.Left() == n {
		p.SetLeft(m)
	} else {
		p.SetRight(m)
	}

	// If n and m are direct siblings, we want to ensure nodes are not
	// swapped back.
	if q.Left() == m && p != q {
		q.SetLeft(n)
	} else {
		q.SetRight(n)
	}

	n.SetParent(q)
	m.SetParent(p)
}
