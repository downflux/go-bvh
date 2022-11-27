// Package balance ensures the local node subtree is properly balanced and has a
// minimal surface area heuristic (SAH).
package balance

import (
	"github.com/downflux/go-bvh/x/internal/cache/node"
)

// B balances a given input node by
//
// 1. ensuring the node itself is balanced after the transformation, and
// 2. minimizing the SAH of the node.
//
// Here, the avl() function is a flat tree height constraint (i.e. after calling
// the function,
//
//	| L.H() - R.H() | <= 1
//
// and rotate() minimizes the SAH.
//
// # This function is based on a mix of
//
// 1. the Catto 2019 slides,
// 2. the briannoyama BVH implementation, and
// 3. Kopta et al. 2012
//
// See the individual function docstrings for more information.
//
// The input node is valid (i.e. its AABB and height are correctly
// pre-calculated), but its children may be imbalanced as a result of an insert
// or delete operation. That is, the height differs at most by one.
//
// If H(n) <= 1 (where leaf nodes have a height of H(n) = 0), this function is a
// no-op.
//
// The returned node has its AABB and height updated.
func BrianNoyama(n node.N) node.N { return Rotate(AVL(n)) }
