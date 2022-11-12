package bvh

import (
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

// avl will look a node x and conditionally optimize the local subtree for
// balanced node heights. That is, given the local subtree
//
//	  x
//	 / \
//	a   z
//	   / \
//	  b   y
//
// where H(y) > H(b), this function will swap a and y so that the resultant tree
// looks like
//
//	  x
//	 / \
//	y   z
//	   / \
//	  b   a
//
// Note that we are treating the BVH as an AVL tree here; that is, we assume the
// local tree was previously balanced (before the subtree rooted at x was
// mutated through an insert or delete operation). Because of this assumption,
// the maximum possible imbalance between the children of x (i.e. a and z) is
// two.
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
// The returned node is the balanced input node.
//
// See the briannoyama implementation for more details.
func avl(x node.N, data map[id.ID]hyperrectangle.R, epsilon float64) node.N {
	if x.Height() < 2 {
		return x
	}

	var a, y node.N
	// WLOG if sz := x.Right() is a leaf node, then its height is 0, and
	// therefore cannot be taller than sa := x.Left(). Therefore, we know sz
	// is an internal node.
	if sa, sz := x.Left(), x.Right(); sa.Height() < sz.Height() {
		if sb, sy := sz.Left(), sz.Right(); sb.Height() < sy.Height() {
			a, y = sa, sy
		} else if sy.Height() < sb.Height() {
			// sb is the shallower node by construction.
			sb, sy = sy, sb
			a, y = sa, sy
		}
	} else if sz.Height() < sa.Height() {
		// sa is the shallower node by construction.
		sa, sz = sz, sa
		if sb, sy := sz.Left(), sz.Right(); sb.Height() < sy.Height() {
			a, y = sa, sy
		} else if sy.Height() < sb.Height() {
			sb, sy = sy, sb
			a, y = sa, sy
		}
	}
	if a != nil && y != nil {
		swap(a, y)

		// By construction, the y node is always the deeper node, and a
		// the shallower node before swapping -- now that they are
		// swapped, the relative depth order is reversed. See the
		// function docstring for the expected subtree layout after
		// swap.
		node.SetAABB(a.Parent(), data, epsilon)
		node.SetHeight(a.Parent())

		node.SetAABB(x, data, epsilon)
		node.SetHeight(x)
	}
	return x
}
