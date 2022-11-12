package bvh

import (
	"math"
	"math/rand"

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
	// include some expansion factor.
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

type Update struct {
	ID   id.ID
	Node cid.ID
}

// insert adds a new AABB into a tree, and returns the new root, along with any
// object node updates.
//
// The input data cache is a read-only map within the insert function.
func insert(c *cache.C, root cid.ID, data map[id.ID]hyperrectangle.R, nodes map[id.ID]cid.ID, x id.ID, expansion float64) (cid.ID, []Update) {
	if root == cid.IDInvalid {
		s := c.GetOrDie(c.Insert(
			cid.IDInvalid,
			cid.IDInvalid,
			cid.IDInvalid,
			/* validate = */ false,
		))
		s.Leaves()[x] = struct{}{}
		return s.ID(), []Update{
			{
				ID:   x,
				Node: s.ID(),
			},
		}
	}

	// t is the new node into which we insert the AABB.
	var t node.N
	aabb := data[x]

	s := sibling(c, root, aabb)
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

		node.SetAABB(s, data, expansion)
		node.SetHeight(s)
		node.SetAABB(t, data, expansion)
		node.SetHeight(t)
	} else {
		t = expand(c, s)
		t.Leaves()[x] = struct{}{}

		node.SetAABB(t, data, expansion)
		node.SetHeight(t)
	}

	// At this point in execution, nodes s and t have updated caches and
	// correct heights. As we traverse up to the root, we will incrementally
	// rebalance the trees.
	var n node.N
	for n = t; n != nil; n = n.Parent() {
		if !n.IsLeaf() {
			node.SetAABB(n, data, expansion)
			node.SetHeight(n)

			n = avl(n, data, expansion)
			// TODO(minkezhang): Optimize for AABB, then rebalance and set
			// height.
		}
	}

	// If we created a new node, we need to broadcast any node changes to
	// the caller.
	updates := make([]Update, 0, c.LeafSize())
	if t != s {
		for n := range t.Leaves() {
			updates = append(updates, Update{
				ID:   n,
				Node: t.ID(),
			})
		}
	}

	// It is possible during repartitioning for the new object to be
	// inserted into the old node. We need to broadcast this change
	// as well.
	if _, ok := t.Leaves()[x]; !ok {
		updates = append(updates, Update{
			ID:   x,
			Node: s.ID(),
		})
	}

	return n.ID(), updates
}
