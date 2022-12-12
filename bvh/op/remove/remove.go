package remove

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/cache"
	"github.com/downflux/go-bvh/internal/cache/node"
	"github.com/downflux/go-bvh/internal/cache/op/balance"
	"github.com/downflux/go-bvh/internal/cache/op/unsafe"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	cid "github.com/downflux/go-bvh/internal/cache/id"
)

var (
	Default = O{
		Balance: balance.BrianNoyamaNoDF,
	}
)

type O struct {
	Balance balance.B
}

func Remove(c *cache.C, data map[id.ID]hyperrectangle.R, x cid.ID, y id.ID, tolerance float64) node.N {
	return Default.Remove(c, data, x, y, tolerance)
}

// Remove deletes a leaf from a node. If necessary, this function will merge
// leaf nodes.
//
// This function returns the new root node.
func (o O) Remove(c *cache.C, data map[id.ID]hyperrectangle.R, x cid.ID, y id.ID, tolerance float64) node.N {
	n, ok := c.Get(x)
	if !ok {
		return nil
	}

	delete(n.Leaves(), y)
	if len(n.Leaves()) == 0 {
		n = unsafe.Remove(c, n)
	}

	var root node.N
	for m := n; m != nil; m = m.Parent() {
		node.SetAABB(m, data, tolerance)
		node.SetHeight(m)

		if !n.IsRoot() {
			m = o.Balance(m)
		}

		if m.Parent() == nil {
			root = m
		}
	}

	return root
}
