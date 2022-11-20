package split

import (
	"fmt"
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/impl"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

var (
	_ S = GuttmanLinear
)

func TestGuttmanLinear(t *testing.T) {
	type w struct {
		n node.N
		m node.N
	}

	type config struct {
		name string
		c    *cache.C
		data map[id.ID]hyperrectangle.R
		n    node.N
		m    node.N
		want w
	}

	configs := []config{
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
				101: *hyperrectangle.New(vector.V{9, 0}, vector.V{10, 1}),
			}

			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			na := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			nb := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nc := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))

			na.SetLeft(nb.ID())
			na.SetRight(nc.ID())

			na.SetHeight(1)

			// Overload a node.
			nb.Leaves()[100] = struct{}{}
			nb.Leaves()[101] = struct{}{}

			wn := impl.New(c, nb.ID())
			wm := impl.New(c, nc.ID())

			wn = impl.New(c, nb.ID())
			wm = impl.New(c, nc.ID())

			wn.Allocate(na.ID(), cid.IDInvalid, cid.IDInvalid)
			wm.Allocate(na.ID(), cid.IDInvalid, cid.IDInvalid)

			wn.Leaves()[100] = struct{}{}
			wm.Leaves()[101] = struct{}{}

			return config{
				name: "LeafSize=1",
				c:    c,
				data: data,
				n:    nb,
				m:    nc,
				want: w{
					n: wn,
					m: wm,
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			node.SetAABB(c.n, c.data, 1)
			h := heuristic.H(c.n.AABB().R())
			nodes := []id.ID{}
			for x := range c.n.Leaves() {
				nodes = append(nodes, x)
			}

			GuttmanLinear(c.c, c.data, c.n, c.m)

			t.Run(fmt.Sprintf("%s/Nodes", c.name), func(t *testing.T) {
				for _, x := range nodes {
					if _, ok := c.n.Leaves()[x]; !ok {
						if _, ok := c.m.Leaves()[x]; !ok {
							t.Errorf("cannot find node %v in output", x)
						}
					}
				}
			})

			t.Run(fmt.Sprintf("%s/H", c.name), func(t *testing.T) {
				node.SetAABB(c.n, c.data, 1)
				node.SetAABB(c.m, c.data, 1)

				if got := heuristic.H(c.n.AABB().R()) + heuristic.H(c.m.AABB().R()); got > h {
					t.Errorf("GuttmanLinear() did not decrease overall heuristic")
				}
			})
		})
	}
}
