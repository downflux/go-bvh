package node

import (
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type N[T point.P] struct {
	data   T
	left   *N[T]
	right  *N[T]
	bounds hyperrectangle.R
}

func (n *N[T]) B() hyperrectangle.R { return n.bounds }
func (n *N[T]) Leaf() bool          { return n.left == nil && n.right == nil }

func (n *N[T]) Data() T  { return n.data }
func (n *N[T]) L() *N[T] { return n.left }
func (n *N[T]) R() *N[T] { return n.right }
