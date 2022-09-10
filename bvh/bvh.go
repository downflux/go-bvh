package bvh

import (
	"fmt"
	"log"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/insert"
	"github.com/downflux/go-bvh/internal/node/op/remove"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

type O struct {
	Size   uint
	Logger *log.Logger
}

// BVH is an AABB-backed bounded volume hierarchy. This struct does not store
// any data associated with the AABBs aside from the object ID. The caller is
// responsible for storing any data associated with the AABB (e.g. in a separate
// map).
type BVH struct {
	lookup map[id.ID]*node.N
	root   *node.N
	size   uint
	logger *log.Logger
}

func New(o O) *BVH {
	return &BVH{
		lookup: map[id.ID]*node.N{},
		size:   o.Size,
		logger: o.Logger,
	}
}

// Insert adds a new AABB bounding box into the BVH tree. The input AABB may be
// larger than the actual object if e.g. the object is not a rectangle, or to
// account for movement updates.
func (bvh *BVH) Insert(x id.ID, aabb hyperrectangle.R) error {
	if bvh.lookup[x] != nil {
		return fmt.Errorf("cannot insert a node with duplicate ID %v", x)
	}

	if bvh.logger != nil {
		bvh.logger.Printf("inserting rectangle ID: %v, AABB: %v", x, aabb)
	}

	n := insert.Execute(bvh.root, bvh.size, x, aabb)

	if bvh.logger != nil {
		bvh.logger.Printf("inserted rectangle into node NID: %v", n.ID())
	}

	// We may have split the leaf node, in which case some data may have
	// shifted.
	for k := range n.Data() {
		bvh.lookup[k] = n
	}
	bvh.root = n.Root()
	if bvh.logger != nil {
		bvh.logger.Printf(
			"tree root NID: %v, Len: %v, H: %v, Imbalance: %v",
			bvh.root.ID(),
			len(bvh.lookup),
			bvh.root.Height(),
			util.MaxImbalance(bvh.root),
		)
		util.Log(bvh.logger, bvh.root)
	}

	return nil
}

// Remove will delete the BVH node which encapsulates the given object.
func (bvh *BVH) Remove(x id.ID) error {
	if bvh.lookup[x] == nil {
		return fmt.Errorf("cannot remove a non-existent object with ID %v", x)
	}

	bvh.root = remove.Execute(bvh.lookup[x], x)
	delete(bvh.lookup, x)

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
func (bvh *BVH) Update(x id.ID, q hyperrectangle.R, aabb hyperrectangle.R) error {
	if bvh.lookup[x] == nil {
		return fmt.Errorf("cannot update a non-existent object with ID %v", x)
	}

	r, ok := bvh.lookup[x].Get(x)
	// If the tracked leaf node cannot find the given object, then something
	// has gone wrong and the tree cannot recover, as the lookup table is no
	// longer trustworthy.
	if !ok {
		panic(fmt.Sprintf("object %v has vanished", x))
	}

	// Update the BVH tree.
	if !bhr.Contains(r, q) {
		if err := bvh.Remove(x); err != nil {
			return fmt.Errorf("cannot update object: %v", err)
		}
		if err := bvh.Insert(x, aabb); err != nil {
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
