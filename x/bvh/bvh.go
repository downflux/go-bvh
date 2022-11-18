package bvh

import (
	"fmt"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

type T struct {
	c    *cache.C
	root cid.ID

	nodes map[id.ID]cid.ID
	data  map[id.ID]hyperrectangle.R

	tolerance float64
}

type O struct {
	K         vector.D
	LeafSize  int
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
	}
}

func (t *T) SAH() float64 {
	n, ok := t.c.Get(t.root)
	if !ok {
		return 0
	}

	return util.SAH(n)
}

func (t *T) H() int {
	n, ok := t.c.Get(t.root)
	if !ok {
		return 0
	}
	return n.Height()
}

func (t *T) K() vector.D { return t.c.K() }

func (t *T) IDs() []id.ID {
	ids := make([]id.ID, 0, len(t.data))
	for x := range t.data {
		ids = append(ids, x)
	}
	return ids
}

func (t *T) Insert(x id.ID, aabb hyperrectangle.R) error {
	if _, ok := t.data[x]; ok {
		return fmt.Errorf("cannot insert a duplicate node %v", x)
	}

	t.data[x] = aabb

	var updates map[id.ID]cid.ID
	t.root, updates = insert(
		t.c, t.root, t.data, t.nodes, x, t.tolerance,
	)
	for x, n := range updates {
		t.nodes[x] = n
	}

	return nil
}

func (t *T) BroadPhase(q hyperrectangle.R) []id.ID { return broadphase(t.c, t.root, t.data, q) }

func (t *T) Update(x id.ID, aabb hyperrectangle.R) error { return nil }
func (t *T) Remove(x id.ID) error                        { return nil }