package bvh

import (
	"fmt"

	"github.com/downflux/go-bvh/filter"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type BVH[T point.P] struct {
	root *node.N

	data map[point.ID]T
}

func New[T point.P](data []T) *BVH[T] {
	panic("unimpemented")
}

func (bvh *BVH[T]) Insert(p T) bool {
	// Cannot re-insert a point which already exists in the BVH.
	if _, ok := bvh.data[p.ID()]; ok {
		panic(fmt.Sprintf("cannot insert a point which already exists in the BVH: %v", p.ID()))
	}

	if bvh.root == nil {
		bvh.root = node.New(node.O{
			ID:    p.ID(),
			Bound: p.Bound(),
		})
	}

	panic("unimplemented")

	return true
}

func (bvh *BVH[T]) Move(id point.ID, offset vector.V) bool {
	panic("unimplemented")
}

func (bvh *BVH[T]) Remove(id point.ID) bool {
	panic("unimplemented")
}

func Collisions[T point.P](bvh BVH[T], q hyperrectangle.R, f filter.F[T]) []T {
	panic("unimplemented")
}
