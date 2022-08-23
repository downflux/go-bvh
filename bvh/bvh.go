package bvh

import (
	"fmt"

	"github.com/downflux/go-bvh/filter"
	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/allocation/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/insert"
	"github.com/downflux/go-bvh/internal/node/move"
	"github.com/downflux/go-bvh/internal/node/remove"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type BVH[T point.RO] struct {
	data   map[point.ID]T
	lookup map[point.ID]id.ID

	allocation allocation.C[*node.N]
	root       id.ID
}

func New[T point.RO](data []T) *BVH[T] {
	bvh := &BVH[T]{
		data:       map[point.ID]T{},
		lookup:     map[point.ID]id.ID{},
		allocation: *allocation.New[*node.N](),
	}
	for _, p := range data {
		bvh.Insert(p)
	}
	return bvh
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

// Move moves a new node into a new position. N.B.: The caller must manually set
// the bound on the point external to this struct.
func (bvh *BVH[T]) Move(i point.ID, r hyperrectangle.R) {
	if _, ok := bvh.data[i]; !ok {
		panic(fmt.Sprintf("attempting to move a point which does not exist in the BVH: %v", i))
	}

	nid := bvh.lookup[i]
	nid, bvh.root = move.Execute(bvh.allocation, nid, r)

	bvh.lookup[i] = nid
}

func (bvh *BVH[T]) Remove(i point.ID) {
	if _, ok := bvh.data[i]; !ok {
		panic(fmt.Sprintf("attempting to delete a point which does not exist in the BVH: %v", i))
	}
	bvh.root = remove.Execute(bvh.allocation, bvh.lookup[i])

	delete(bvh.data, i)
	delete(bvh.lookup, i)
}

func RangeSearch[T point.RO](bvh BVH[T], q hyperrectangle.R, f filter.F[T]) []T {
	panic("unimplemented")
}
