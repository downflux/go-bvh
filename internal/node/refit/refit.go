package refit

import (
	"github.com/downflux/go-bvh/hyperrectangle"
	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/allocation/id"
	"github.com/downflux/go-bvh/internal/node"
)

func Execute(nodes allocation.C[*node.N], nid id.ID) {
	n := nodes[nid]
	if n == nil {
		return
	}

	bound := n.Bound()
	// Walk back up the tree refitting AABBs.
	for n := node.Parent(nodes, n); n != nil; n = node.Parent(nodes, n) {
		n.SetBound(hyperrectangle.Union(bound, n.Bound()))
	}
}
