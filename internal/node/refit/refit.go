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

	// Walk back up the tree refitting AABBs.
	for p := node.Parent(nodes, n); p != nil; p = node.Parent(nodes, p) {
		var s *node.N
		if node.Left(nodes, p) == n {
			s = node.Right(nodes, p)
		} else {
			s = node.Left(nodes, p)
		}
		p.SetBound(hyperrectangle.Union(n.Bound(), s.Bound()))
		n = p
	}
}
