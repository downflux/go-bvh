package rotation

import (
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/op/rotate/subtree"
)

type R struct {
	B, C, F, G *node.N
}

func Generate(n *node.N) []R {
	t := subtree.New(n)

	var rs []R
	if !t.A.Leaf() { // t.B and t.C are non-nil
		if !t.B.Leaf() { // t.F and t.G are non-nil
			rs = append(rs, R{
				B: t.B, C: t.C, F: t.F, G: t.G,
			}, R{
				B: t.B, C: t.C, F: t.F, G: t.G,
			})
		}
		if !t.C.Leaf() { // t.E and t.D are non-nil
			rs = append(rs, R{
				B: t.C, C: t.B, F: t.E, G: t.D,
			}, R{
				B: t.C, C: t.B, F: t.D, G: t.E,
			})
		}
	}

	return rs
}
