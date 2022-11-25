package split

import (
	"math"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

// DHConnelly implements the leaf node splitting function as used in
// github.com/dhconnelly/rtreego.
//
// N.B.: The rtreego implementation considers a minimum leaf size; this is not a
// feature we support, and therefore the (minimal) corresponding logic in
// rtreego which takes this into consideration is not implemented here.
func DHConnelly(c *cache.C, data map[id.ID]hyperrectangle.R, n node.N, m node.N) {
	if c.LeafSize() == 1 {
		for x := range n.Leaves() {
			m.Leaves()[x] = struct{}{}
			delete(n.Leaves(), x)
			break
		}
		return
	}

	buf := hyperrectangle.New(
		vector.V(make([]float64, c.K())),
		vector.V(make([]float64, c.K())),
	).M()

	leaves := make([]id.ID, 0, len(n.Leaves()))
	for x := range n.Leaves() {
		leaves = append(leaves, x)
	}

	// Reset the source node's leaves. We are skipping the same operation
	// for the destination node as we expect that node to be empty.
	for _, x := range leaves {
		delete(n.Leaves(), x)
	}

	li, ri := seed(data, leaves, buf)

	// Use the AABB objects from the n and m nodes as buffers here.
	n.Leaves()[leaves[li]] = struct{}{}
	m.Leaves()[leaves[ri]] = struct{}{}
	n.AABB().Copy(data[leaves[li]])
	m.AABB().Copy(data[leaves[ri]])
	n.SetHeuristic(heuristic.H(n.AABB().R()))
	m.SetHeuristic(heuristic.H(m.AABB().R()))

	remaining := append(leaves[:li], leaves[li+1:ri]...)
	remaining = append(remaining, leaves[ri+1:]...)

	for len(remaining) > 0 {
		ni := next(data, remaining, n, m, buf)
		aabb := data[remaining[ni]]

		p := group(aabb, n, m, buf)
		p.Leaves()[remaining[ni]] = struct{}{}
		p.AABB().Union(aabb)
		p.SetHeuristic(heuristic.H(p.AABB().R()))

		remaining = append(remaining[:ni], remaining[ni+1:]...)
	}

	return
}

// group determines if an input AABB object should be put into the left or right
// node. This function assumes the node bounding box (i.e. n.AABB()) and node
// heuristic cache (n.Heuristic()) are up-to-date.
func group(aabb hyperrectangle.R, n node.N, m node.N, buf hyperrectangle.M) node.N {
	buf.Copy(aabb)
	buf.Union(n.AABB().R())
	lh := heuristic.H(buf.R())

	buf.Copy(aabb)
	buf.Union(m.AABB().R())
	rh := heuristic.H(buf.R())

	ld := lh - n.Heuristic()
	rd := rh - m.Heuristic()

	if ld < rd {
		return n
	}
	if ld > rd {
		return m
	}

	// In the case the increase in heuristic is equal, use the node with
	// less total surface area.
	if n.Heuristic() < m.Heuristic() {
		return n
	}
	if n.Heuristic() > m.Heuristic() {
		return m
	}

	// In the case this too is equal, choose the node with the least amount
	// of elements.
	if len(n.Leaves()) < len(m.Leaves()) {
		return n
	}
	return m
}

// next picks a leaf object to be considered for the left / right node
// placement. This function assumes that the node bounding box and heuristic
// cache are valid.
func next(data map[id.ID]hyperrectangle.R, leaves []id.ID, n node.N, m node.N, buf hyperrectangle.M) int {
	var next int

	d := math.Inf(-1)
	for i, x := range leaves {
		aabb := data[x]

		buf.Copy(aabb)
		buf.Union(n.AABB().R())
		ld := heuristic.H(buf.R()) - n.Heuristic()

		buf.Copy(aabb)
		buf.Union(m.AABB().R())
		rd := heuristic.H(buf.R()) - m.Heuristic()

		if e := math.Abs(ld - rd); e > d {
			d = e
			next = i
		}
	}

	return next
}

// seed picks two initial leaf objects to be placed in the left and right nodes.
func seed(data map[id.ID]hyperrectangle.R, leaves []id.ID, buf hyperrectangle.M) (int, int) {
	var l, r int
	h := math.Inf(-1)

	for i, x := range leaves {
		for j, y := range leaves[i+1:] {
			buf.Copy(data[x])
			buf.Union(data[y])

			if g := heuristic.H(buf.R()) - (heuristic.H(data[x]) + heuristic.H(data[y])); g > h {
				h = g
				l = i
				r = j + i + 1
			}
		}
	}

	return l, r
}
