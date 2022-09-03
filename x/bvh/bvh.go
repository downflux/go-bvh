package bvh

import (
	"fmt"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/op/insert"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

// BVH is an AABB-backed bounded volume hierarchy. This struct does not store
// any data associated with the AABBs aside from the object ID. The caller is
// responsible for storing any data associated with the AABB (e.g. in a separate
// map).
type BVH struct {
	lookup map[id.ID]*node.N
	root   *node.N
}

func New() *BVH {
	return &BVH{
		lookup: map[id.ID]*node.N{},
	}
}

// Insert adds a new AABB bounding box into the BVH tree. The input AABB should
// be larger than the actual object to account for movement updates.
func (bvh *BVH) Insert(id id.ID, aabb hyperrectangle.R) {
	n := insert.Execute(bvh.root, id, aabb)
	bvh.lookup[id] = n
	bvh.root = n.Root()
}

// Remove will delete the BVH node which encapsulates the given object.
func (bvh *BVH) Remove(id id.ID) error {
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
func (bvh *BVH) Update(id id.ID, q hyperrectangle.R, aabb hyperrectangle.R) error {
	return fmt.Errorf("unimplemented")
}

func BroadPhase(bvh *BVH, aabb hyperrectangle.R) []id.ID { return bvh.root.BroadPhase(aabb) }
