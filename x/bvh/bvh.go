package bvh

import (
	"fmt"

	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

// BVH is an AABB-backed bounded volume hierarchy. This struct does not store
// any data associated with the AABBs aside from the object ID. The caller is
// responsible for storing any data associated with the AABB (e.g. in a separate
// map).
type BVH[T comparable] struct {
	lookup map[T]*node.N[T]
	root   *node.N[T]
}

type D[T comparable] struct {
	ID   T
	AABB hyperrectangle.R
}

func New[T comparable](data []D[T]) *BVH[T] {
	bvh := &BVH[T]{
		lookup: map[T]*node.N[T]{},
	}

	for _, p := range data {
		bvh.Insert(p.ID, p.AABB)
	}
	return bvh
}

// Insert adds a new AABB bounding box into the BVH tree. The input AABB should
// be larger than the actual object to account for movement updates.
func (bvh *BVH[T]) Insert(id T, aabb hyperrectangle.R) error {
	return fmt.Errorf("unimplemented")
}

// Remove will delete the BVH node which encapsulates the given object.
func (bvh *BVH[T]) Remove(id T) error {
	return fmt.Errorf("unimplemented")
}

// Update will conditionally update the BVH tree if the new position of the
// bounding box is no longer completely contained by the associated BVH node.
//
// The input aabb is the new bounding box which will be used to encapsulate the
// object.
//
// N.B.: The BVH is not responsible for updating the object itself -- the caller
// will need to do that separately. This function is called only to update the
// state of the BVH.
func (bvh *BVH[T]) Update(id T, q hyperrectangle.R, aabb hyperrectangle.R) error {
	return fmt.Errorf("unimplemented")
}
