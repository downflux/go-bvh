package tree

import (
	"math"
	"math/rand"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

// expand creates a new node with s as its sibling. This will re-link any
// existing parents or siblings of s and ensure that the generated cache is
// still valid.
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
func expand(c *cache.C, s *cache.N) *cache.N {
	q := s.Parent()

	p := c.GetOrDie(c.Insert(q.ID(), cache.IDInvalid, cache.IDInvalid, false))
	n := c.GetOrDie(c.Insert(p.ID(), cache.IDInvalid, cache.IDInvalid, false))

	q.SetChild(q.Branch(s.ID()), p.ID())
	s.SetParent(p.ID())

	p.SetLeft(s.ID())
	p.SetRight(n.ID())

	return n
}

// partition splits a full node s by moving some objects into a new node t.
//
// We assume t is an empty leaf node and the size of s exceeds the cache leaf
// size (that is, it can afford to lose a single object without becoming empty).
func partition(s *cache.N, t *cache.N, axis vector.D, data map[id.ID]hyperrectangle.R) {
	// Find the upper and lower bounds of the tight-fitting AABB for the
	// data in s. This helps circumvent the case where the cached AABB for s
	// includes some expansion factor.
	kmin := math.Inf(1)
	kmax := math.Inf(-1)
	for x := range s.Data() {
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
	for y := range s.Data() {
		ymin := data[y].Min().X(axis)
		if ymin > xmin {
			x = y
			xmin = ymin
		}

		if ymin > kmid {
			delete(s.Data(), y)
			t.Data()[y] = struct{}{}
		}
	}

	// Ensure t has some object in it.
	if len(t.Data()) == 0 {
		delete(s.Data(), x)
		t.Data()[x] = struct{}{}
	}
}

type Update struct {
	ID   id.ID
	Node cache.ID
}

// setAABB sets a leaf node's AABB with a given expansion factor. The input node
// must be a leaf node and contain at least one object.
func setAABB(data map[id.ID]hyperrectangle.R, n *cache.N, c float64) {
	var initialized bool
	var k vector.D
	for x := range n.Data() {
		if !initialized {
			n.AABB().Copy(data[x])
			k = data[x].Min().Dimension()
		} else {
			n.AABB().Union(data[x])
		}
	}
	n.AABB().Scale(math.Pow(c, 1/float64(k)))
}

// insert adds a new AABB into a tree, and returns the new root, along with any
// object node updates.
//
// The input data cache is a read-only map within the insert function.
func insert(c *cache.C, root cache.ID, data map[id.ID]hyperrectangle.R, nodes map[id.ID]cache.ID, x id.ID, expansion float64) (cache.ID, []Update) {
	if root == cache.IDInvalid {
		s := c.GetOrDie(c.Insert(
			cache.IDInvalid,
			cache.IDInvalid,
			cache.IDInvalid,
			/* validate = */ false,
		))
		s.Data()[x] = struct{}{}
		return s.ID(), []Update{
			{
				ID:   x,
				Node: s.ID(),
			},
		}
	}

	// t is the new node into which we insert the AABB.
	var t *cache.N
	aabb := data[x]

	s := c.GetOrDie(sibling(c, root, aabb))
	if s.IsLeaf() {
		// If the leaf is full, we need repartition the leaf and split
		// its children into a new node.
		if s.IsFull() {
			s.Data()[x] = struct{}{}

			t = expand(c, s)
			partition(s, t, vector.D(rand.Intn(int(c.K()))), data)
		} else {
			s.Data()[x] = struct{}{}
		}

		setAABB(data, s, expansion)
		setAABB(data, t, expansion)
	} else {
		t = expand(c, s)
		t.Data()[x] = struct{}{}

		setAABB(data, t, expansion)
	}

	// At this point in execution, nodes s and t have updated caches.

	var n *cache.N
	for n = t; n != nil; n = n.Parent() {
		if !n.IsLeaf() {
			n.AABB().Copy(n.Left().AABB().R())
			n.AABB().Union(n.Right().AABB().R())
		}
		// TODO: Rebalance and set height.
	}

	// If we created a new node, we need to broadcast any node changes to
	// the caller.
	updates := make([]Update, 0, c.LeafSize())
	if t != s {
		for n := range t.Data() {
			updates = append(updates, Update{
				ID:   n,
				Node: t.ID(),
			})
		}
	}

	// It is possible during repartitioning for the new object to be
	// inserted into the old node. We need to broadcast this change
	// as well.
	if _, ok := t.Data()[x]; !ok {
		updates = append(updates, Update{
			ID:   x,
			Node: s.ID(),
		})
	}

	return n.ID(), updates
}
