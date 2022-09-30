package path

import (
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/cache"

type P struct {
	P *node.N
	B node.Branch
}

func (p P) N(c *cache.C[*node.N]) *node.N { return p.P.Child(c, p.B) }

func (p P) Next(c *cache.C[*node.N]) P {
	n := p.N(c)
	if n.IsRoot(c) {
		panic("iterator error")
	}

	
	q := n.Parent(c)
	return P{
		P: q,
		B: qp.Branch(n),
	}
}
