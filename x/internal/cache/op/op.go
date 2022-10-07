package op

import (
	"github.com/downflux/go-bvh/x/internal/cache"
)

// Swap moves two nodes in the same tree.
//
// Case: parent / child
//
//	  n
//	 / \
//	A   m
//	   / \
//	  B   C
//
// Swap(n, m) should generate the following tree
//
//	  m
//	 / \
//	B   n
//	   / \
//	  A   C
//
// N.B.: the order of A and C here do not matter, since our BVH tree is
// invariant under child swaps.
//
// Note too that this may make the overall tree quality worse -- additonal
// checks will be necessary to swap B
//
// Case: child / child
//
//    A
//   / \
//  n   m
//
// will be invariant (since it doesn't matter).
//
// Case: ancestor / child
//
//    n
//   / \
//  A   B
//     / \
//    C   m
//       / \
//      D   E
//
// to
//
//    m
//   / \
//  D   n
//     / \
//    A   B
//       / \
//      C   E

func Swap(c *cache.C, from cache.ID, to cache.ID) {
	n := c.GetOrDie(from)
	m := c.GetOrDie(to)

	// Handle child / child case.
	if n.Parent() == m.Parent() {
		return
	}

}
