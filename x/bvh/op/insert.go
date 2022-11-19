package op

import (
	"math"
	"math/rand"

	"github.com/downflux/go-bvh/x/bvh/balance"
	"github.com/downflux/go-bvh/x/bvh/sibling"
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/op/unsafe"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

// partition splits a full node s by moving some objects into a new node t.
//
// We assume t is an empty leaf node and the size of s exceeds the cache leaf
// size (that is, it can afford to lose a single object without becoming empty).
func split(c *cache.C, data map[id.ID]hyperrectangle.R, from node.N, to node.N) {
	if c.LeafSize() == 2 {
		for x := range from.Leaves() {
			to.Leaves()[x] = struct{}{}
			delete(from.Leaves, x)
			return
		}
	}

	high := make([]id.ID, c.K())
	low := make([]id.ID, c.K())

	for i := vector.D(0); i < c.K(); i++ {
		separation[i] = math.Inf(1)
	}

	node.SetAABB(from, data, 1)
	buf := hyperrectangle.New(
		vector.V(make([]float64, c.K())),
		vector.V(make([]float64, c.K())),
	).M()
	buf.Copy(node.AABB().R())

	for i := vector.D(0); i < c.K(); i++ {
		min := math.Inf(1)
		max := math.Inf(-1)
		separation := 0

		for x := range from.Leaves() {
			aabb := data[x]

			if max - min

			if 	
			high.SetX(i, math.Max(high.X(i), aabb.Min().X(i)))
			low.SetX(i, math.Min(low.X(i), aabb.Max().X(i)))
			buf.Union(aabb)
		}
	}

	seperation := math.Inf(-1)
	var k vector.D
	for i := vector.D(0); i < c.K(); i++ {
		normal := buf.Max().X(i) - buf.Min().X(i)
		high.SetX(i, high.X(i)/normal)
		low.SetX(i, low.X(i)/normal)

		if s := high.X(i) - low.X(i); s > separation {
			separation = s
			k = i
		}
	}

}

func Insert(c *cache.C, root cid.ID, data map[id.ID]hyperrectangle.R, nodes map[id.ID]cid.ID, x id.ID, tolerance float64) (cid.ID, map[id.ID]cid.ID) {
	updates := make(map[id.ID]cid.ID, c.LeafSize())

	var n node.N
	var ok bool
	if n, ok = c.Get(root); !ok {
		n = c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, false))
	} else {
		n = candidate.BrianNoyama(c, n)
	}

	if n.IsFull() {
		n.Leaves()[x] = struct{}{}
		m := unsafe.Expand(c, n)
		split(c, data, n, m)

		for _, y := range m.Leaves() {
			updates[y] = m.ID()
		}
		if _, ok := updates[x]; !ok {
			updates[x] = n.ID()
		}
	} else {
		n.Leaves()[x] = struct{}{}

		updates[x] = n.ID()
	}

	var r node.N
	for m := n; m != nil; m = m.Parent() {
		if !m.IsLeaf() {
			node.SetAABB(m, data, tolerance)
			node.SetHeight(m)

			m = balance.B(m, data, tolerance)
		}

		if m.Parent() == nil {
			r = m
		}
	}

	return r.ID(), updates
}

// insert adds a new AABB into a tree, and returns the new root, along with any
// object node updates.
//
// The input data cache is a read-only map within the insert function.
func insert(c *cache.C, root cid.ID, data map[id.ID]hyperrectangle.R, nodes map[id.ID]cid.ID, x id.ID, tolerance float64) (cid.ID, map[id.ID]cid.ID) {
	s, t := raw(c, root, data, x)

	updates := make(map[id.ID]cid.ID, c.LeafSize())
	if s != nil {
		node.SetAABB(s, data, tolerance)
		node.SetHeight(s)

		// x may be inserted into either s or t.
		if _, ok := s.Leaves()[x]; ok {
			updates[x] = s.ID()
		}
	}

	if t != nil {
		node.SetAABB(t, data, tolerance)
		node.SetHeight(t)

		// In the course of creating the new node t, we will need to
		// update any migrated nodes.
		for x := range t.Leaves() {
			updates[x] = t.ID()
		}
	}

	// Walk back up the tree, while at the same time getting the root. Since
	// nodes s and t are already balanced by construction, the balancing op
	// will only start at the shared parent.
	n := s
	if n == nil {
		n = t
	}

	// At this point in execution, nodes s and t have updated caches and
	// correct heights. As we traverse up to the root, we will incrementally
	// rebalance the trees.
	var r node.N
	for m := n; m != nil; m = m.Parent() {
		if !m.IsLeaf() {
			node.SetAABB(m, data, tolerance)
			node.SetHeight(m)

			m = balance.B(m, data, tolerance)
		}

		if m.Parent() == nil {
			r = m
		}
	}

	return r.ID(), updates
}
