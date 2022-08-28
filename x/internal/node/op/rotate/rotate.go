package rotate

import (
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/op/rotate/rotation"
)

func Execute(a *node.N) *node.N {
	if a == nil {
		panic("cannot rotate a nil node")
	}

	if !a.IsLeaf() {
		if r := rotation.Optimal(a); r != (rotation.R{}) {
			r.B.Swap(r.F)
		}
	}

	if a.IsRoot() {
		return a
	}
	return Execute(a.Parent())
}
