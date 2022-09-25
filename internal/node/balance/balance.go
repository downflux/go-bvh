package balance

import (
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/balance/aabb"
	"github.com/downflux/go-bvh/internal/node/balance/rotate"
	"github.com/downflux/go-bvh/internal/node/op/rotate/swap"
)

func Execute(n *node.N) *node.N {
	rotate.Rebalance(n)
	if n.IsRoot() {
		return n
	}

	p := n.Parent()
	r := aabb.Query(p)
	if r != (aabb.R{}) {
		swap.Execute(r.Src, r.Target)
	}
	if p.IsRoot() {
		return n
	}
	return Execute(p)
}
