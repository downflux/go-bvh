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

// L is a left rotate operation on the right child of x.
//
// N.B.: In normal binary search trees, the order of the left and right children
// is important, i.e.
//
//   A
//  / \
// B   C
//
// Is considered different from
//
//   A
//  / \
// C   B
//
// A BVH does not have this constraint, so we are free to do a direct swap in
// these rotations.
//
// The returned node is the root of the subtree. In our implementation, this is
// always x.
//
// We can imagine the rotation as
//
//   x
//  / \
// a   z
//    / \
//   b   c
//
// to
//
//   x
//  / \
// c   z
//    / \
//   b   a
//
// See https://en.wikipedia.org/wiki/AVL_tree for more information.
func L(x *node.N) *node.N {
	if x.Height() < 2 {
		panic(fmt.Sprintf("cannot rotate %v with height %v", x.ID(), x.Height()))
	}
	swap.Execute(x.Left(), x.Right().Right())
	return x
}

// R is a right rotate operation on the left child of x.
func R(x *node.N) *node.N {
	if x.Height() < 2 {
		panic(fmt.Sprintf("cannot rotate %v with height %v", x.ID(), x.Height()))
	}
	swap.Execute(x.Right(), x.Left().Left())
	return x
}

// RL is a composite right-left rotate operation where z is the right
// child of x and z is left-heavy.
//
// As the name implies, this rotation breaks into two separate rotations --
//
// z' := R(x.Left())
// L(z'.Parent())
//
// See https://en.wikipedia.org/wiki/AVL_tree for more information.
func RL(x *node.N) *node.N {
	if x.Height() < 3 {
		panic(fmt.Sprintf("cannot rotate %v with height %v", x.ID(), x.Height()))
	}

	if x.Right().Height() != x.Left().Height()+2 {
		panic(fmt.Sprintf("cannot rotate %v to imbalance the tree", x.ID()))
	}

	return L(R(x.Right()).Parent())
}

// LR is a composite left-right rotate operation where z is the left
// child of x and z is right-heavy.
func LR(x *node.N) *node.N {
	if x.Height() < 3 {
		panic(fmt.Sprintf("cannot rotate %v with height %v", x.ID(), x.Height()))
	}

	if x.Left().Height() != x.Right().Height()+2 {
		panic(fmt.Sprintf("cannot rotate %v to imbalance the tree", x.ID()))
	}

	return R(L(x.Left()).Parent())
}
