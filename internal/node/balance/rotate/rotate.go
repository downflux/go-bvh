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

// Execute will look a node y and its ancestors z and x, and conditionally swap an
// with the a node from the opposite subtree. That is, given
//
//	  x
//	 / \
//	a   z
//	   / \
//	  b   y
//
// Execute may do a swap of a and y, depending on the subtree heights of a and
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
// The returned node is the input node.
//
// See the briannoyama implementation for more details.
func Execute(x *node.N) *node.N {
	if x == nil {
		panic("cannot rebalance an empty node")
	}

	if x.IsLeaf() {
		return x
	}

	var b, f *node.N
	// If x.Left() is leaf, then its height is minimal, and therefore we
	// implicitly know x.Left() is an internal node here.
	if x.Left().Height() > x.Right().Height() {
		if x.Left().Left().Height() > x.Left().Right().Height() {
			b, f = x.Left().Left(), x.Right()
		} else if x.Left().Right().Height() > x.Left().Left().Height() {
			b, f = x.Left().Right(), x.Right()
		}
	} else if x.Right().Height() > x.Left().Height() {
		if x.Right().Left().Height() > x.Right().Right().Height() {
			b, f = x.Right().Left(), x.Left()
		} else if x.Right().Right().Height() > x.Right().Left().Height() {
			b, f = x.Right().Right(), x.Left()
		}
	}
	if b != nil && f != nil {
		swap.Execute(b, f)
	}
	return x
}
