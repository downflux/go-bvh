// Package rotate defines the operations on a node subtree. Here, x is the root
// of the subtree, and z is one of its children.
//
// In-order tree traversal is preserved under rotation.
//
// See https://en.wikipedia.org/wiki/Tree_rotation and
// https://en.wikipedia.org/wiki/AVL_tree for more information.
package rotate

import (
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/rotate/swap"
)

// Rebalance will look a node y and its ancestors z and x, and conditionally swap an
// with the a node from the opposite subtree. That is, given
//
//	  x
//	 / \
//	a   z
//	   / \
//	  b   y
//
// Rebalance may do a swap of a and y, depending on the subtree heights of a and
// y.
//
// Note that we are treating the BVH here as an AVL tree; that is, we assume the
// tree was previously balanced, and we only need to do the toation here due to
// a single insert or delete operation, and therefore, the maximum imbalance
// possible between sibling nodes (e.g. a and z) is two.
//
// Further note that the rebalance operation is simplified because a BVH tree
// is invariant under sibling swaps --
//
//	  x
//	 / \
//	a   z
//
// and
//
//	  x
//	 / \
//	z   a
//
// Behave exactly the same, since the only property we care about in the query
// operation is on AABB intersections, and does not rely on the order between z
// and a, as is the case in e.g. a binary search tree.
//
// In an AVL tree where in-order traversal behavior needs to be preserved and
// WLOG z is the right child of x, we will need to apply the standard L or RL on
// x, depending on if the c or b is heavier, respectively.
//
// The returned node is the parent node of the input.
//
// See the briannoyama implementation for more details.
func Rebalance(y *node.N) *node.N {
	if y == nil {
		panic("cannot rebalance an empty node")
	}

	if y.IsRoot() {
		return nil
	}
	if y.Parent().IsRoot() {
		return y.Parent()
	}
	if y.IsLeaf() {
		return y.Parent()
	}

	z := y.Parent()
	x := z.Parent()
	a := map[bool]*node.N{
		true:  x.Right(),
		false: x.Left(),
	}[z == x.Left()]

	if y.Height() > a.Height() {
		swap.Execute(a, y)
	}
	return z
}
