package path

import (
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/cache"

type P struct {
	N *node.N
	B node.Branch
}

func (p P) Next(c *cache.C[*node.N]) P {
	if p.N.IsRoot(c) {
		return P{
			N: nil,
			B: node.BranchInvalid,
		}
	}

	
}
