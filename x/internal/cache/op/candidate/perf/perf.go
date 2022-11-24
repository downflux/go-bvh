package perf

import (
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

type G func() (*cache.C, node.N)

func Trivial() (*cache.C, node.N) {
	c := cache.New(cache.O{
		LeafSize: 1,
		K:        2,
	})

	root := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
	root.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}))

	return c, root
}

func Balanced(n int) (*cache.C, node.N) {
	horizon := make([]node.N, 0, n)

	c := cache.New(cache.O{
		LeafSize: 1,
		K:        2,
	})

	for i := 0; i < n; i++ {
		l := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
		l.AABB().Copy(*hyperrectangle.New(vector.V{float64(i), 0}, vector.V{float64(i) + 1, 1}))

		horizon = append(horizon, l)
	}

	for len(horizon) > 1 {
		var next []node.N
		for i := 0; i < len(horizon)-1; i += 2 {
			l := horizon[i]
			r := horizon[i+1]

			p := c.GetOrDie(c.Insert(cid.IDInvalid, l.ID(), r.ID(), true))
			node.SetAABB(p, nil, 1)
			node.SetHeight(p)

			l.SetParent(p.ID())
			r.SetParent(p.ID())

			next = append(next, p)
		}
		if len(horizon)%2 == 1 {
			next = append(next, horizon[len(horizon)-1])
		}
		horizon = next
	}

	return c, horizon[0]
}
