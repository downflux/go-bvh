package bvh

import (
	"fmt"

	"github.com/downflux/go-bvh/filter"
	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/allocation/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/insert"
	"github.com/downflux/go-bvh/internal/node/remove"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type BVH[T point.P] struct {
	data   map[point.ID]T
	lookup map[point.ID]id.ID

	allocation allocation.C[*node.N]
	root       id.ID
}

func New[T point.P](data []T) *BVH[T] {
	return &BVH[T]{
		data:       map[point.ID]T{},
		lookup:     map[point.ID]id.ID{},
		allocation: *allocation.New[*node.N](),
	}
}

func (bvh *BVH[T]) Insert(p T) {
	// Cannot re-insert a point which already exists in the BVH.
	if _, ok := bvh.data[p.ID()]; ok {
		panic(fmt.Sprintf("cannot insert a point which already exists in the BVH: %v", p.ID()))
	}

	bvh.data[p.ID()] = p
	// TODO(minkezhang): Expand bound by some percentage.
	var i id.ID
	i, bvh.root = insert.Execute(bvh.allocation, bvh.root, p.ID(), p.Bound())
	bvh.lookup[p.ID()] = i
}

func (bvh *BVH[T]) Move(id point.ID, dp vector.V) bool {
	panic("unimplemented")
}

func (bvh *BVH[T]) Remove(i point.ID) {
	if _, ok := bvh.data[i]; !ok {
		panic(fmt.Sprintf("attempting to delete a point which does not exist in the BVH: %v", i))
	}
	bvh.root = remove.Execute(bvh.allocation, bvh.lookup[i])

	delete(bvh.data, i)
	delete(bvh.lookup, i)
}

func RangeSearch[T point.P](bvh BVH[T], q hyperrectangle.R, f filter.F[T]) []T {
	panic("unimplemented")
}
