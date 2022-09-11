package bvh

import (
	"fmt"
	"log"
	"sync"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-bvh/internal/node/op/insert"
	"github.com/downflux/go-bvh/internal/node/op/remove"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

type O struct {
	Size   uint
	Logger *log.Logger
}

type M struct {
	Height       uint
	MaxImbalance uint
	Cost         float64
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

	l sync.RWMutex
}

func New(o O) *BVH {
	return &BVH{
		lookup: map[id.ID]*node.N{},
		size:   o.Size,
		logger: o.Logger,
	}
}

// updateNodeCache ensures that the underling nodes have an updated cache of
// dynamically calculated properties. These cached properties may have been
// invalidated upon a tree mutation. Explicitly generating the cache ensures we
// can call read operations in parallel.
//
// We assume the write mutex is already held in this function.
func (bvh *BVH) updateNodeCache() {
	if bvh.root != nil {
		bvh.root.AABB()
		bvh.root.Height()
	}
}

func (bvh *BVH) IDs() []id.ID {
	bvh.l.RLock()
	defer bvh.l.RUnlock()

	ids := make([]id.ID, 0, len(bvh.lookup))
	for x := range bvh.lookup {
		ids = append(ids, x)
	}
	return ids
}

func (bvh *BVH) Report() M {
	bvh.l.RLock()
	defer bvh.l.RUnlock()

	return M{
		Height:       bvh.root.Height(),
		MaxImbalance: util.MaxImbalance(bvh.root),
		Cost:         util.Cost(bvh.root),
	}
}

// Insert adds a new AABB bounding box into the BVH tree. The input AABB may be
// larger than the actual object if e.g. the object is not a rectangle, or to
// account for movement updates.
func (bvh *BVH) Insert(x id.ID, aabb hyperrectangle.R) error {
	bvh.l.Lock()
	defer bvh.l.Unlock()

	return bvh.insert(x, aabb)
}

func (bvh *BVH) insert(x id.ID, aabb hyperrectangle.R) error {
	if bvh.lookup[x] != nil {
		return fmt.Errorf("cannot insert a node with duplicate ID %v", x)
	}

	n := insert.Execute(bvh.root, bvh.size, x, aabb)

	// We may have split the leaf node, in which case some data may have
	// shifted.
	for k := range n.Data() {
		bvh.lookup[k] = n
	}
	bvh.root = n.Root()
	bvh.updateNodeCache()

	return nil
}

// Remove will delete the BVH node which encapsulates the given object.
func (bvh *BVH) Remove(x id.ID) error {
	bvh.l.Lock()
	defer bvh.l.Unlock()

	return bvh.remove(x)
}

func (bvh *BVH) remove(x id.ID) error {
	if bvh.lookup[x] == nil {
		return fmt.Errorf("cannot remove a non-existent object with ID %v", x)
	}

	bvh.root = remove.Execute(bvh.lookup[x], x)
	bvh.updateNodeCache()
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
	bvh.l.Lock()
	defer bvh.l.Unlock()

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
		if err := bvh.remove(x); err != nil {
			return fmt.Errorf("cannot update object: %v", err)
		}
		if err := bvh.insert(x, aabb); err != nil {
			return fmt.Errorf("cannot update object: %v", err)
		}
	}

	return nil
}

// BroadPhase returns all objects which may collide with the input query.  The
// list of objects returned may not actually collide, e.g. if the actual hitbox
// in user-space is smaller than the query AABB. The user is responsible for
// further collision refinement.
func (bvh *BVH) BroadPhase(q hyperrectangle.R) []id.ID {
	bvh.l.RLock()
	defer bvh.l.RUnlock()

	if bvh.root == nil {
		return nil
	}
	return bvh.root.BroadPhase(q)
}
