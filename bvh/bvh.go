package bvh

import (
	"fmt"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/insert"
	"github.com/downflux/go-bvh/internal/node/op/remove"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
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

// Insert adds a new AABB bounding box into the BVH tree. The input AABB may be
// larger than the actual object if e.g. the object is not a rectangle, or to
// account for movement updates.
func (bvh *BVH) Insert(id id.ID, aabb hyperrectangle.R) error {
	if bvh.lookup[id] != nil {
		return fmt.Errorf("cannot insert a node with duplicate ID %v", id)
	}

	n := insert.Execute(bvh.root, id, aabb)

	// We may have split the leaf node, in which case some data may have
	// shifted.
	for _, x := range n.Data() {
		bvh.lookup[x] = n
	}
	bvh.root = n.Root()

	return nil
}

// Remove will delete the BVH node which encapsulates the given object.
func (bvh *BVH) Remove(id id.ID) error {
	if bvh.lookup[id] == nil {
		return fmt.Errorf("cannot remove a non-existent object with ID %v", id)
	}

	bvh.root = remove.Execute(bvh.lookup[id], id)
	delete(bvh.lookup, id)

	return nil
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
	if bvh.lookup[id] == nil {
		return fmt.Errorf("cannot update a non-existent object with ID %v", id)
	}

	r, ok := bvh.lookup[id].Get(id)
	// If the tracked leaf node cannot find the given object, then something
	// has gone wrong and the tree cannot recover, as the lookup table is no
	// longer trustworthy.
	if !ok {
		panic(fmt.Sprintf("object %v has vanished", id))
	}

	// Update the BVH tree.
	if !bhr.Contains(r, q) {
		if err := bvh.Remove(id); err != nil {
			return fmt.Errorf("cannot update object: %v", err)
		}
		if err := bvh.Insert(id, aabb); err != nil {
			return fmt.Errorf("cannot update object: %v", err)
		}
	}

	return nil
}

// BroadPhase returns all objects which may collide with the input query.  The
// list of objects returned may not actually collide, e.g. if the actual hitbox
// in user-space is smaller than the query AABB. The user is responsible for
// further collision refinement.
func BroadPhase(bvh *BVH, q hyperrectangle.R) []id.ID { return bvh.root.BroadPhase(q) }
