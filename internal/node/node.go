package node

import (
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type N[T point.P] struct {
	data   int
	left   int
	right  int
	bounds hyperrectangle.R
}

func (n *N[T]) B() hyperrectangle.R { return n.bounds }
func (n *N[T]) Leaf() bool          { return n.left < 0 && n.right < 0 }

func L[T point.P](n *N[T], nodes []*N[T]) *N[T] {
	if n.left < 0 {
		return nil
	}
	return nodes[n.left]
}

func R[T point.P](n *N[T], nodes []*N[T]) *N[T] {
	if n.right < 0 {
		return nil
	}
	return nodes[n.right]
}

func Data[T point.P](n *N[T], data []T) T {
	if n.Leaf() {
		var blank T
		return blank
	}
	return data[n.data]
}
