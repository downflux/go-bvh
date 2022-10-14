package op

import (
	"fmt"

	"github.com/downflux/go-bvh/x/internal/cache"
)

func IsAncestor(c *cache.C, n cache.ID, m cache.ID) bool {
	for x := m; x.IsValid(); x = c.GetOrDie(x).Parent() {
		if x == n {
			return true
		}
	}
	return false
}

// Swap moves two nodes in the same tree. This function does not support
// swapping ancestors, e.g.
//
//	  n
//	 / \
//	A   m
//	   / \
//	  B   C
//
// we unconditionally set this to
//
//	  m
//	 / \
//	B   n
//	   / \
//	  A   C
//
// The caller should be aware of this and make further optimizations if e.g. the
// optimal configuration would have been
//
//	  m
//	 / \
//	C   n
//	   / \
//	  A   B
func Swap(c *cache.C, from cache.ID, to cache.ID, validate bool) {
	// We will call validate only in debugging situations, as this is an
	// O(log N) check.
	if validate && (IsAncestor(c, from, to) || IsAncestor(c, to, from)) {
		panic(fmt.Sprintf("cannot swap ancestor nodes %v, %v", from, to))
	}

	n, m := c.GetOrDie(from), c.GetOrDie(to)
	p, _ := c.Get(n.Parent())
	q, _ := c.Get(m.Parent())

	// Update parent links to the children.
	if p != nil {
		p.SetChild(p.Branch(n.ID()), m.ID())
	}
	if q != nil {
		q.SetChild(q.Branch(m.ID()), n.ID())
	}

	// Update child links to the parent.
	n.SetParent(q.ID())
	m.SetParent(p.ID())
}
