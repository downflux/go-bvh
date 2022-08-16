package node

import (
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type N[T point.P] struct {
	data T

	parent *N[T]
	left   *N[T]
	right  *N[T]
	bounds hyperrectangle.R
}

func (n *N[T]) B() hyperrectangle.R { return n.bounds }
func (n *N[T]) Leaf() bool          { return n.left == nil && n.right == nil }

func (n *N[T]) Data() T       { return n.data }
func (n *N[T]) L() *N[T]      { return n.left }
func (n *N[T]) R() *N[T]      { return n.right }
func (n *N[T]) Parent() *N[T] { return n.parent }

func (n *N[T]) Move(id point.ID, offset vector.V) bool { return false }
func (n *N[T]) Remove(id point.ID) bool                { return false }
func (n *N[T]) Insert(p T)                             {}
