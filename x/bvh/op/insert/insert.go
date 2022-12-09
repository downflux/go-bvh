package insert

import (
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/op/balance"
	"github.com/downflux/go-bvh/x/internal/cache/op/candidate"
	"github.com/downflux/go-bvh/x/internal/cache/op/split"
	"github.com/downflux/go-bvh/x/internal/cache/op/unsafe"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

var (
	Default = O{
		Candidate: candidate.Guttman,
		Split:     split.GuttmanLinear,
		Balance:   balance.BrianNoyamaNoDF,
	}
)

type O struct {
	Candidate candidate.C
	Split     split.S
	Balance   balance.B
}

func Insert(c *cache.C, rid cid.ID, data map[id.ID]hyperrectangle.R, x id.ID, tolerance float64) (node.N, []node.N) {
	return Default.Insert(c, rid, data, x, tolerance)
}

// insert adds a new AABB into a tree, and returns the new root, along with any
// object node updates.
//
// The input data cache is a read-only map within the insert function.
func (o O) Insert(c *cache.C, rid cid.ID, data map[id.ID]hyperrectangle.R, x id.ID, tolerance float64) (node.N, []node.N) {
	var mutations []node.N

	root, ok := c.Get(rid)
	if !ok {
		root = c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, false))
	}

	// s is a leaf node. This leaf node may be full.
	s := o.Candidate(c, root, data[x])
	s.Leaves()[x] = struct{}{}

	mutations = append(mutations, s)

	if len(s.Leaves()) > c.LeafSize() {
		t := unsafe.Expand(c, s)
		o.Split(c, data, s, t)

		mutations = append(mutations, t)

		s = t
	}

	for _, n := range mutations {
		node.SetAABB(n, data, tolerance)
		node.SetHeight(n)
	}

	// At this point in execution, all leaf nodes have updated caches. As we
	// traverse up to the root, we will incrementally rebalance the trees.
	for m := s; m != nil; m = m.Parent() {

		if !m.IsLeaf() {
			node.SetAABB(m, nil, 1)
			node.SetHeight(m)

			m = o.Balance(m)
		}

		if m.Parent() == nil {
			root = m
		}
	}

	return root, mutations
}
