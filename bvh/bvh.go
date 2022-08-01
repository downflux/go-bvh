package bvh

import (
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type BVH[T point.P] struct {
	nodes []*node.N[T]
	data  []T
}

func New[T point.P](data []T) *BVH[T] {
	panic("unimpemented")
}

func (bvh *BVH[T]) Insert(t T) {
	panic("unimplemented")
}

func (bvh *BVH[T]) Remove(r hyperrectangle.R, f func(p T) bool) (T, bool) {
	panic("unimplemented")
}

func Collisions[T point.P](bvh BVH[T], q hyperrectangle.R, f func(p T) bool) []T {
	panic("unimplemented")
}
