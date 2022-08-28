package rotate

import (
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/op/rotate/rotation"
)

func Execute(n *node.N) *node.N {
	if n == nil {
		panic("cannot rotate a nil node")
	}

	if n.IsLeaf() {
		return Execute(n.Parent())
	}

	r := rotation.Optimal(n)
	if r == (rotation.R{}) {
		if n.IsRoot() {
			return n
		}
		return Execute(n.Parent())
	}

	return r.B
}
