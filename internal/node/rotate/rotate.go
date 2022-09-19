// Package rotate defines the operations on a node subtree. Here, x is the root
// of the subtree, and z is one of its children.
//
// In-order tree traversal is preserved under rotation.
//
// See https://en.wikipedia.org/wiki/Tree_rotation and
// https://en.wikipedia.org/wiki/AVL_tree for more information.
package rotate

import (
	"fmt"

	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/rotate/swap"
)

// Rebalance will look a node z and its parent x, and conditionally swap an
// unbalanced child of z with the sibling of z. That is, given
//
//	  x
//	 / \
//	a   z
//	   / \
//	  b   c
//
// Rebalance will do a swap(a, b) or swap(a, c), depending on the subtree
// heights.
//
// Note that we are treating the BVH here as an AVL tree; that is, we assume the
// tree was previously balanced, and we only need to do the toation here due to
// a single insert or delete operation, and therefore, the maximum imbalance
// possible between sibling nodes (e.g. a and z) is two.
//
// Further note that the rebalance operation is simplified because a BVH tree
// is invariant under sibling swaps --
//
//	  A
//	 / \
//	B   C
//
// and
//
//	  A
//	 / \
//	C   B
//
// Behave exactly the same, since the only property we care about in the query
// operation is on AABB intersections, and does not rely on the order between B
// and C, as is the case in e.g. a binary search tree.
//
// In an AVL tree where in-order traversal behavior needs to be preserved and
// WLOG z is the right child of x, we will need to apply the standard L or RL on
// x, depending on if the c or b is heavier, respectively.
//
// The returned node is the parent node.
func Rebalance(x *node.N, z *node.N) *node.N {
	if x.IsLeaf() {
		panic(fmt.Sprintf("parent node %v cannot be a leaf", x.ID()))
	}
	if z.IsRoot() {
		panic(fmt.Sprintf("child node %v cannot be the root node", z.ID()))
	}
	if z.Parent() != x || (x.Left() != z && x.Right() != z) {
		panic(fmt.Sprintf("%v is not a parent of %v", x.ID(), z.ID()))
	}
	if z.IsLeaf() {
		return x
	}

	var m *node.N
	if z == x.Left() {
		m = x.Right()
	} else {
		m = x.Left()
	}

	if m.Height() < z.Height()+2 {

	}
}

func RebalanceHeight(n *node.N) *node.N {
	if n.IsRoot() {
		return n
	}

	// Check the relative heights between n and its sibling -- if the height
	// is too imbalanced, rotate the tree accordingly such that the
	// p-subtree is better balanced. This is the (single-)node rebalance
	// step for a self-balancing binary tree.
	p := n.Parent()
	if n == p.Left() {
		// Handle the case of the n-subtree mutation being the result of
		// an insert operation.
		if p.Right().Height() > n.Height()+1 {
			return r(p, n)
		}
		// Handle the case of the n-subtree mutation being the result of
		// a remove operation.
		if p.Right().Height()+1 < n.Height() {
			return l(p, n)
		}
	}
	if n == p.Right() {
		if p.Left().Height() > n.Height()+1 {
			return l(p, n)
		}
		if p.Left().Height()+1 > n.Height() {
			return r(p, n)
		}
	}
	return n
}

func rotate(p *node.N, n *node.N) *node.N {
	if n.IsRoot() || n.Parent() != p || (p.Left() != n) && (p.Right() != n) {
		panic("invalid parent / child node relationship")
	}

	if p.Right() == n {
		return l(p)
	}
	return r(p)
}

// l is a left rotate operation on the right child of x.
//
// N.B.: In normal binary search trees, the order of the left and right children
// is important, i.e.
//
//	  A
//	 / \
//	B   C
//
// Is considered different from
//
//	  A
//	 / \
//	C   B
//
// A BVH does not have this constraint, so we are free to do a direct swap in
// these rotations.
//
// The returned node is the root of the subtree. In our implementation, this is
// always x.
//
// We can imagine the rotation as
//
//	  x
//	 / \
//	a   z
//	   / \
//	  b   c
//
// to
//
//	  x
//	 / \
//	c   z
//	   / \
//	  b   a
//
// See https://en.wikipedia.org/wiki/AVL_tree for more information.
func l(x *node.N) *node.N {
	if x.Height() < 2 {
		panic(fmt.Sprintf("cannot rotate %v with height %v", x.ID(), x.Height()))
	}
	if x.Right().IsLeaf() {
		panic(fmt.Sprintf("cannot rotate leaf %v", x.Right().ID()))
	}

	swap.Execute(x.Left(), x.Right().Right())
	return x
}

// r is a right rotate operation on the left child of x.
func r(x *node.N) *node.N {
	if x.Height() < 2 {
		panic(fmt.Sprintf("cannot rotate %v with height %v", x.ID(), x.Height()))
	}
	if x.Left().IsLeaf() {
		panic(fmt.Sprintf("cannot rotate leaf %v", x.Left().ID()))
	}

	swap.Execute(x.Right(), x.Left().Left())
	return x
}

// rl is a composite right-left rotate operation where z is the right
// child of x and z is left-heavy.
//
// As the name implies, this rotation breaks into two separate rotations --
//
// z' := R(x.Left())
// L(z'.Parent())
//
// See https://en.wikipedia.org/wiki/AVL_tree for more information.
func rl(x *node.N) *node.N {
	if x.Height() < 3 {
		panic(fmt.Sprintf("cannot rotate %v with height %v", x.ID(), x.Height()))
	}

	if x.Right().Height() != x.Left().Height()+2 {
		panic(fmt.Sprintf("cannot rotate %v to imbalance the tree", x.ID()))
	}

	return l(r(x.Right()).Parent())
}

// lr is a composite left-right rotate operation where z is the left
// child of x and z is right-heavy.
func lr(x *node.N) *node.N {
	if x.Height() < 3 {
		panic(fmt.Sprintf("cannot rotate %v with height %v", x.ID(), x.Height()))
	}

	if x.Left().Height() != x.Right().Height()+2 {
		panic(fmt.Sprintf("cannot rotate %v to imbalance the tree", x.ID()))
	}

	return r(l(x.Left()).Parent())
}
