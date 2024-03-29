// Package BVH implements an AABB-backed BVH tree.
package bvh

import (
	"fmt"

	"github.com/downflux/go-bvh/bvh/op/insert"
	"github.com/downflux/go-bvh/bvh/op/query"
	"github.com/downflux/go-bvh/bvh/op/remove"
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/cache"
	"github.com/downflux/go-bvh/internal/cache/node/util/metrics"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/ray"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/internal/cache/id"
)

type T struct {
	c    *cache.C
	root cid.ID

	nodes map[id.ID]cid.ID
	data  map[id.ID]hyperrectangle.R

	tolerance float64

	insert insert.O
	remove remove.O
}

type O struct {
	K        vector.D
	LeafSize int

	// Tolerance specifies the bounding buffer width around leaf nodes as a
	// percentage of the volume of the AABB. This value must be greater than
	// one (as the resultant AABB must encapsulate the leaf).
	Tolerance float64
}

func New(o O) *T {
	if o.Tolerance < 1 {
		panic(fmt.Sprintf("cannot set tolerance factor %v < 1", o.Tolerance))
	}

	return &T{
		c: cache.New(cache.O{
			K:        o.K,
			LeafSize: o.LeafSize,
		}),

		root: cid.IDInvalid,

		nodes:     make(map[id.ID]cid.ID, 1024),
		data:      make(map[id.ID]hyperrectangle.R, 1024),
		tolerance: o.Tolerance,

		insert: insert.Default,
		remove: remove.Default,
	}
}

// SAH returns the surface area heuristic as defined by MacDonald and Booth
// 1990. The heuristic constants are set to the values specified in Aila et al.
func (t *T) SAH() float64 {
	n, ok := t.c.Get(t.root)
	if !ok {
		return 0
	}

	return metrics.SAH(n)
}

func (t *T) IDs() []id.ID {
	ids := make([]id.ID, 0, len(t.data))
	for x := range t.data {
		ids = append(ids, x)
	}
	return ids
}

// Insert adds a new AABB into the BVH. The specific data structure which tracks
// this AABB is managed by the user (external to this library). After the AABB
// is mutated (e.g. during a simulation tick), the user must call Update to
// ensure the tree remains valid.
//
// Because the AABB must remain static inside the BVH, we will create a new copy
// of the input.
func (t *T) Insert(x id.ID, aabb hyperrectangle.R) error {
	if _, ok := t.data[x]; ok {
		return fmt.Errorf("cannot insert a duplicate node %v", x)
	}

	buf := hyperrectangle.New(
		vector.V(make([]float64, aabb.Min().Dimension())),
		vector.V(make([]float64, aabb.Min().Dimension())),
	).M()
	buf.Copy(aabb)

	t.data[x] = buf.R()

	root, mutations := t.insert.Insert(
		t.c, t.root, t.data, x, t.tolerance,
	)
	t.root = root.ID()
	for _, n := range mutations {
		for x := range n.Leaves() {
			t.nodes[x] = n.ID()
		}
	}

	return nil
}

// Update will move a corresponding object. Depending on the BVH tolerance and
// how fast an object is moving, we would expect this function to filter out a
// large number of Delete and subsequent Insert calls.
func (t *T) Update(x id.ID, aabb hyperrectangle.R) error {
	if _, ok := t.data[x]; !ok {
		return fmt.Errorf("cannot update a non-existent node %v", x)
	}

	n := t.c.GetOrDie(t.nodes[x])
	if !hyperrectangle.Contains(n.AABB().R(), aabb) {
		if err := t.Remove(x); err != nil {
			return fmt.Errorf("cannot update node %v: %v", x, err)
		}
		if err := t.Insert(x, aabb); err != nil {
			return fmt.Errorf("cannot update node %v: %v", x, err)
		}
	} else {
		// As with the Insert call, we do not have any guarantees the
		// AABB will not be mutated after the update call, so we must
		// ensure this is true by making a copy of the input.
		t.data[x].M().Copy(aabb)
	}

	return nil
}

func (t *T) Remove(x id.ID) error {
	if _, ok := t.data[x]; !ok {
		return fmt.Errorf("cannot remove a non-existent node %v", x)
	}

	root := t.remove.Remove(
		t.c, t.data, t.nodes[x], x, t.tolerance,
	)
	if root != nil {
		t.root = root.ID()
	} else {
		t.root = cid.IDInvalid
	}

	delete(t.nodes, x)
	delete(t.data, x)

	return nil
}

// BroadPhase finds all objects which intersect with the given input AABB.
func (t *T) BroadPhase(q hyperrectangle.R) []id.ID {
	return query.BroadPhase(t.c, t.root, t.data, q)
}

// Raycast finds all objcets which intersects the given ray.
func (t *T) Raycast(q ray.R) []id.ID {
	return query.Raycast(t.c, t.root, t.data, q)
}

// Query finds all objects which passes the input filtering function. BroadPhase
// and Raycast are special cases of the Query function. The input filter will be
// recursively applied; that is, the child of an internal BVH node will be
// searched only if the parent AABB also passes the filter.
func (t *T) Query(f func(r hyperrectangle.R) bool) []id.ID {
	return query.Query(t.c, t.root, t.data, f)
}
