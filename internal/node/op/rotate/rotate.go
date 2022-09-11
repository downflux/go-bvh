package rotate

import (
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/rotate/rotation"
	"github.com/downflux/go-bvh/internal/node/op/rotate/swap"
)

type R func(n *node.N) rotation.R

func Execute(a *node.N, f R) *node.N {
	if a == nil {
		panic("cannot rotate a nil node")
	}

	if !a.IsLeaf() {
		if r := f(a); r != (rotation.R{}) {
			swap.Execute(r.B, r.F)
		}
	}

	if a.IsRoot() {
		return a
	}
	return Execute(a.Parent(), f)
}
