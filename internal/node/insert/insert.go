package insert

import (
	"fmt"

	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-pq/pq"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

// Execute adds the given point into the tree. If a new node is created, it will
// be created with a new index.
//
// This function will return the new root.
func Execute(nodes allocation.C[*node.N], root allocation.ID, id point.ID, bound hyperrectangle.R) allocation.ID {
	if _, ok := nodes[root]; !ok {
		nid := nodes.Allocate()
		n := node.New(node.O{
			ID:    id,
			Index: nid,
			Bound: bound,
		})
		if err := nodes.Insert(nid, n); err != nil {
			panic(fmt.Sprintf("cannot insert node: %s", err))
		}
		return n.Index()
	}

	// Find best new sibling for the new leaf.
	cid := findSibling(nodes, root, bound)

	// Create new parent.
	pid := createParent(nodes, cid, id, bound)

	// Rebalance the tree up to the root.
	rotate(nodes, nodes[pid])

	return pid
}

// rotate finds the configuration of the children and grandchildren which
// minimizes the local heuristic, and rotates the tree accordingly. See Catto
// 2019 for more information. As we are applyiing potential rotations on the
// descendents of the input node, neither the bounds of the parent nor its own
// parent change.
//
// We refer to the local subtree using the same labels as the Catto
// presentation, i.e.
//
//        A
//       / \
//      /   \
//     B     C
//    / \   / \
//   D   E F   G
//
// We will attempt to apply a rotation on the node pairs
//
//   BF, BG, CD, and CE
//
// and choose the configuration with the minimal "surface area" heuristic.
//
// N.B.: The Box2D implementation uses pure tree heights as the balancing
// heuristic.
//
// TODO(minkezhang): Investigate relative performance characteristics of
// choosing between these two balancing modes.
func rotate(nodes allocation.C[*node.N], a *node.N) {
	var b, c, d, e, f, g *node.N

	if a == nil {
		return
	}

	b = node.Left(nodes, a)
	c = node.Right(nodes, a)

	if b != nil {
		d = node.Left(nodes, b)
		e = node.Right(nodes, b)
	}
	if c != nil {
		f = node.Left(nodes, b)
		g = node.Right(nodes, b)
	}

	// We are using the BF rotation syntax as explored in the Catto
	// presentation. Here, we want to check if the b node should be swapped
	// with the f node (i.e. from b to f). The c node is therefore the
	// sibling of the b node, and g is the sibling of the f node (i.e. the
	// other child of c).
	type r struct {
		b *node.N

		c *node.N
		// Note that if f is non-nil, then both the parent c and sibling
		// g are also non-nil.
		f *node.N
		g *node.N
	}
	rs := []r{
		{b: b, c: c, f: f, g: g},
		{b: b, c: c, f: g, g: f},
		{b: c, c: b, f: e, g: d},
		{b: c, c: b, f: d, g: e},
	}

	var optimal r
	h := heuristic.Heuristic(b.Bound()) + heuristic.Heuristic(c.Bound())
	for _, rotation := range rs {
		if rotation.b != nil && rotation.f != nil {
			swapped := heuristic.Heuristic(rotation.f.Bound()) + heuristic.Heuristic(
				bhr.Union(
					rotation.b.Bound(),
					rotation.g.Bound(),
				),
			)
			if swapped < h {
				optimal = rotation
				h = swapped
			}
		}
	}

	// Apply a rotation if there is a configuration which has a lower
	// surface area heuristic.
	if optimal.b != nil && optimal.f != nil {
		if node.Left(nodes, a) == optimal.b {
			a.SetLeft(optimal.f.Index())
		} else {
			a.SetRight(optimal.f.Index())
		}

		if node.Left(nodes, optimal.c) == optimal.f {
			optimal.c.SetLeft(optimal.b.Index())
		} else {
			optimal.c.SetRight(optimal.b.Index())
		}

		optimal.b.SetParent(optimal.c.Index())
		optimal.f.SetParent(a.Index())
	}

	rotate(nodes, node.Parent(nodes, a))
}

// createParent creates a new parent node for a candidate r. This parent will
// have have the old node and a newly-created node with the given bounds. The
// sibling ID rid must exist and refer to a non-nil node.
//
// This function will modify the allocation table as a side-effect.
//
// This function returns the allocation ID of the newly-created parent.
func createParent(nodes allocation.C[*node.N], rid allocation.ID, id point.ID, bound hyperrectangle.R) allocation.ID {
	r := nodes[rid]
	if r == nil {
		panic(fmt.Sprintf("given right allocation ID does not exist: %v", rid))
	}

	pid := nodes.Allocate()
	lid := nodes.Allocate()

	l := node.New(node.O{
		ID: id,

		Index:  lid,
		Parent: pid,

		Bound: bound,
	})

	var aid allocation.ID
	if node.Parent(nodes, r) != nil {
		aid = node.Parent(nodes, r).Index()
	}

	p := node.New(node.O{
		Index:  pid,
		Parent: aid,
		Left:   lid,
		Right:  rid,

		Bound: bhr.Union(bound, r.Bound()),
	})
	nodes.Insert(lid, l)
	nodes.Insert(pid, p)
	if node.Parent(nodes, r) != nil {
		if node.Left(nodes, node.Parent(nodes, r)) == r {
			node.Parent(nodes, r).SetLeft(pid)
		} else {
			node.Parent(nodes, r).SetRight(pid)
		}
	}
	r.SetParent(pid)
	node.Left(nodes, p).SetParent(pid)

	// Walk back up the tree refitting AABBs.
	for n := nodes[pid]; n != nil; n = node.Parent(nodes, n) {
		n.SetBound(bhr.Union(bound, n.Bound()))
	}

	return pid
}

// findSibling finds the node to which an object with the given bound will
// siblings. A new parent node will be created above both the sibling and the
// input bound. This is based on the branch-and-bound algorithm (Catto 2019).
func findSibling(nodes allocation.C[*node.N], rid allocation.ID, bound hyperrectangle.R) allocation.ID {
	root, ok := nodes[rid]
	if !ok {
		panic("cannot find candidate for an empty root node")
	}

	q := pq.New[*node.N](0)

	// Note that the priority queue is a max-heap, so we will need to flip
	// the heuristic signs.
	q.Push(root, -float64(heuristic.Direct(root, bound)))

	c := root
	d := -q.Priority()

	for q.Len() > 0 {
		n := q.Pop()

		inherited := heuristic.Inherited(nodes, n, bound)

		// Check if the current node is a better insertion candidate.
		actual := float64(heuristic.Direct(n, bound) + inherited)
		if actual < d {
			c = n
			d = actual
		}

		if !node.Leaf(nodes, n) {
			// Append queue children to the queue if the lower bound for
			// inserting into the child is less than the current minimum
			// (i.e. there's room for optimization).
			//
			// Note that the inherited heuristic is the same between the
			// left and right children.
			inherited += heuristic.Heuristic(bhr.Union(bound, n.Bound())) - heuristic.Heuristic(n.Bound())
			estimated := heuristic.Heuristic(bound) + inherited

			if float64(estimated) < d {
				q.Push(node.Left(nodes, n), -float64(estimated))
				q.Push(node.Right(nodes, n), -float64(estimated))
			}
		}
	}
	return c.Index()
}
