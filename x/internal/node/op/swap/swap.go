package swap

import (
	"github.com/downflux/go-bvh/x/internal/node"
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

	if isAncestor(n, m) || isAncestor(m, n) {
		panic("cannot swap a node with its ancestor")
	}

	if !n.IsRoot() {
		p := n.Parent()
		if p.Left() == n {
			p.SetLeft(m)
		} else {
			p.SetRight(m)
		}
		p.InvalidateAABBCache()
	}
	if !m.IsRoot() {
		p := m.Parent()
		if p.Left() == m {
			p.SetLeft(n)
		} else {
			p.SetRight(n)
		}
		p.InvalidateAABBCache()
	}
	p := n.Parent()
	q := m.Parent()
	n.SetParent(q)
	m.SetParent(p)
}
