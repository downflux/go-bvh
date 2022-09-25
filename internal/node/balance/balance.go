package balance

import (
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/balance/aabb"
	"github.com/downflux/go-bvh/internal/node/balance/rotate"
	"github.com/downflux/go-bvh/internal/node/op/rotate/swap"
)

func Execute(n *node.N) *node.N {
	rotate.Execute(n)

	if r := aabb.Query(n); r != (aabb.R{}) {
		swap.Execute(r.Src, r.Target)
	}
	if n.IsRoot() {
		return n
	}
	return Execute(n.Parent())
}
