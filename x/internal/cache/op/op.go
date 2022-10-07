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
//	C   n
//	   / \
//	  B   C
//
// N.B.: the order of B and b here do not matter, since our BVH tree is
// invariant under child swaps.
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
//n
//

func Swap(c *cache.C, from cache.ID, to cache.ID) {
	n := c.GetOrDie(from)
	m := c.GetOrDie(to)

}
