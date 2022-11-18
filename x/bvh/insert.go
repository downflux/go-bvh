package bvh

import (
	"math"
	"math/rand"

	"github.com/downflux/go-bvh/x/bvh/balance"
	"github.com/downflux/go-bvh/x/bvh/sibling"
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

// expand creates a new node with s as its sibling. This will re-link any
// existing parents or siblings of s and ensure that the generated cache is
// still valid.
//
// The input node s must not be nil.
//
// The input node s is a node within the cache.
//
//	  Q
//	 / \
//	S   T
//
// to
//
//	    Q
//	   / \
//	  P   T
//	 / \
//	S   N
func expand(c *cache.C, s node.N) node.N {
	if s == nil {
		panic("cannot expand a nil node")
	}

	var q node.N
	qid := cid.IDInvalid
	if !s.IsRoot() {
		q = s.Parent()
		qid = q.ID()
	}

	p := c.GetOrDie(c.Insert(qid, cid.IDInvalid, cid.IDInvalid, false))
	n := c.GetOrDie(c.Insert(p.ID(), cid.IDInvalid, cid.IDInvalid, false))

	if q != nil {
		q.SetChild(q.Branch(s.ID()), p.ID())
	}

	s.SetParent(p.ID())

	p.SetLeft(s.ID())
	p.SetRight(n.ID())

	return n
}

// partition splits a full node s by moving some objects into a new node t.
//
// We assume t is an empty leaf node and the size of s exceeds the cache leaf
// size (that is, it can afford to lose a single object without becoming empty).
func partition(s node.N, t node.N, axis vector.D, data map[id.ID]hyperrectangle.R) {
	// Find the upper and lower bounds of the tight-fitting AABB for the
	// data in s. We need to calculate this directly as the cached AABB may
	// include some tolerance factor.
	kmin := math.Inf(1)
	kmax := math.Inf(-1)
	for x := range s.Leaves() {
		if min := data[x].Min().X(axis); min < kmin {
			kmin = min
		}
		if max := data[x].Max().X(axis); max > kmax {
			kmax = max
		}
	}
	kmid := kmin + (kmax-kmin)/2.0

	// Because it is possible that no object in s will exist only in the
	// right-handed domain of the pivot plane, we will keep an index of the
	// "right-most" object that will be forcibly moved into t if needed.
	var x id.ID
	xmin := math.Inf(-1)

	// Because of the way kmin is constructed, we guarantee that there is at
	// least one object still left in s.
	for y := range s.Leaves() {
		ymin := data[y].Min().X(axis)
		if ymin > xmin {
			x = y
			xmin = ymin
		}

		if ymin > kmid {
			delete(s.Leaves(), y)
			t.Leaves()[y] = struct{}{}
		}
	}

	// Ensure t has some object in it.
	if len(t.Leaves()) == 0 {
		delete(s.Leaves(), x)
		t.Leaves()[x] = struct{}{}
	}
}

// raw is the base insert into the tree. This method only guarantees the
// newly-created nodes have valid links back up to (potentially a new) root.
//
// This function returns an (s, t) node tuple, where s is the existing insertion
// sibling candidate, and t is a newly-created node (in the case that s is
// full).
func raw(c *cache.C, root cid.ID, data map[id.ID]hyperrectangle.R, x id.ID) (node.N, node.N) {
	if _, ok := c.Get(root); !ok {
		t := c.GetOrDie(c.Insert(
			cid.IDInvalid,
			cid.IDInvalid,
			cid.IDInvalid,
			/* validate = */ false,
		))
		t.Leaves()[x] = struct{}{}
		return nil, t
	}

	// t is the new node into which we insert the AABB.
	var t node.N
	aabb := data[x]

	s := sibling.Find(c, root, aabb)
	if s == nil {
		panic("cannot find valid insertion sibling candidate")
	}

	if s.IsLeaf() {
		// If the leaf is full, we need repartition the leaf and split
		// its children into a new node.
		if s.IsFull() {
			s.Leaves()[x] = struct{}{}

			t = expand(c, s)
			partition(s, t, vector.D(rand.Intn(int(c.K()))), data)
		} else {
			s.Leaves()[x] = struct{}{}
		}
	} else {
		t = expand(c, s)
		t.Leaves()[x] = struct{}{}
	}

	return s, t
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
