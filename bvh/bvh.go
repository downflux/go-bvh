package bvh

import (
	"fmt"

	"github.com/downflux/go-bvh/filter"
	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/insert"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type BVH[T point.P] struct {
	data map[point.ID]T

	allocation allocation.C[*node.N]
	root       allocation.ID
}

func New[T point.P](data []T) *BVH[T] {
	return &BVH[T]{
		data:       map[point.ID]T{},
		allocation: *allocation.New[*node.N](),
	}
}

func (bvh *BVH[T]) Insert(p T) {
	// Cannot re-insert a point which already exists in the BVH.
	if _, ok := bvh.data[p.ID()]; ok {
		panic(fmt.Sprintf("cannot insert a point which already exists in the BVH: %v", p.ID()))
	}

	bvh.data[p.ID()] = p
	bvh.root = insert.New(bvh.allocation).Execute(bvh.root, p.ID(), p.Bound())
}

func (bvh *BVH[T]) Move(id point.ID, dp vector.V) bool {
	panic("unimplemented")
}

func (bvh *BVH[T]) Remove(id point.ID) bool {
	panic("unimplemented")
}

func RangeSearch[T point.P](bvh BVH[T], q hyperrectangle.R, f filter.F[T]) []T {
	panic("unimplemented")
}
