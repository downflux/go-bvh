package cache

import (
	"fmt"
)

func IsAncestor(c *C, n ID, m ID) bool {
	var x *N
	for x = c.GetOrDie(m); !x.IsRoot(); x = x.Parent() {
		if x.ID() == n {
			return true
		}
	}
	return x.ID() == n
}

// Swap moves two nodes in the same tree. This function does not support
// swapping ancestors, e.g.
//
//	  n
//	 / \
//	A   m
//	   / \
//	  B   C
func Swap(c *C, from ID, to ID, validate bool) {
	// We will call validate only in debugging situations, as this is an
	// O(log N) check.
	if validate && (IsAncestor(c, from, to) || IsAncestor(c, to, from)) {
		panic(fmt.Sprintf("cannot swap ancestor nodes %v, %v", from, to))
	}

	n, m := c.GetOrDie(from), c.GetOrDie(to)
	var p, q *N
	if !n.IsRoot() {
		p = n.Parent()
	}
	if !m.IsRoot() {
		q = m.Parent()
	}

	var b, d B
	if p != nil {
		b = p.Branch(n.ID())
	}
	if q != nil {
		d = q.Branch(m.ID())
	}

	// Update parent links to the children.
	if b.IsValid() {
		p.SetChild(b, m.ID())
	}
	if d.IsValid() {
		q.SetChild(d, n.ID())
	}

	// Update child links to the parent.
	n.SetParent(q.ID())
	m.SetParent(p.ID())
}
