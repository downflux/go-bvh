package split

import (
	"math"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/cache"
	"github.com/downflux/go-bvh/internal/cache/node"
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

// GuttmanLinear implements the linear split algorithm as defined in Guttman
// 1984, section 3.5.3.
//
// Here, n is the source (i.e. full) node, and m is the destination (empty)
// node.
//
// N.B.: This will leave the nodes n and m in an inconsistent state. The caller
// will be responsible for updating the bounding box of these leaves via
// node.SetAABB().
func GuttmanLinear(c *cache.C, data map[id.ID]hyperrectangle.R, n node.N, m node.N) {
	if c.LeafSize() == 1 {
		for x := range n.Leaves() {
			m.Leaves()[x] = struct{}{}
			delete(n.Leaves(), x)
			break
		}
		return
	}

	// Use the source node AABB as a scratch space to calculate the
	// tightly-bound AABB.
	node.SetAABB(n, data, 1)
	buf := hyperrectangle.New(
		vector.V(make([]float64, c.K())),
		vector.V(make([]float64, c.K())),
	).M()
	buf.Copy(n.AABB().R())

	// Reset the leaves within the source node, as data will be copied into
	// here.
	nodes := make([]id.ID, 0, len(n.Leaves()))
	for x := range n.Leaves() {
		nodes = append(nodes, x)
		delete(n.Leaves(), x)
	}

	// separation tracks the normalized maximum separation factor between
	// AABBs across all dimensions, where separation is defined as the
	// distance between the maximum of the lower bounds and the minimum of
	// the upper bounds of AABBs along that dimension.
	separation := math.Inf(-1)

	// WLOG left tracks the AABB which will be the seed for n, i.e. the AABB
	// which contributes the maximum lower bound of the separation factor.
	var left id.ID
	var right id.ID

	// Pick node seeds.
	for i := vector.D(0); i < c.K(); i++ {
		var kl id.ID
		var kr id.ID

		if data[nodes[0]].Max().X(i) < data[nodes[1]].Max().X(i) {
			kl, kr = nodes[0], nodes[1]
		} else {
			kl, kr = nodes[1], nodes[0]
		}

		// WLOG klower tracks the maximum lower bound.
		klower := math.Max(data[kl].Min().X(i), data[kr].Min().X(i))
		kupper := math.Min(data[kl].Max().X(i), data[kr].Max().X(i))

		for _, x := range nodes[2:] {
			aabb := data[x]

			if k := aabb.Max().X(i); k < kupper {
				kupper = k
				kl = x
			} else if k := aabb.Min().X(i); k > klower && x != kl {
				klower = k
				kr = x
			}
		}

		if s := (klower - kupper) / (buf.Max().X(i) - buf.Min().X(i)); s > separation {
			separation = s
			left = kl
			right = kr
		}
	}

	// Set node seeds.
	n.Leaves()[left] = struct{}{}
	m.Leaves()[right] = struct{}{}

	n.AABB().Copy(data[left])
	m.AABB().Copy(data[right])

	// Set AABBs based on the smallest net increase in node size.
	for _, x := range nodes {
		aabb := data[x]

		if x == left || x == right {
			continue
		}

		lh := heuristic.H(n.AABB().R())
		rh := heuristic.H(m.AABB().R())

		buf.Copy(n.AABB().R())
		buf.Union(aabb)
		dlh := heuristic.H(buf.R()) - lh

		buf.Copy(m.AABB().R())
		buf.Union(aabb)
		drh := heuristic.H(buf.R()) - rh

		if dlh < drh {
			n.Leaves()[x] = struct{}{}
		} else {
			m.Leaves()[x] = struct{}{}
		}
	}
}
