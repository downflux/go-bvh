package candidate

import (
	"testing"

	"github.com/downflux/go-bvh/internal/cache"
	"github.com/downflux/go-bvh/internal/cache/node"
	"github.com/downflux/go-bvh/internal/cache/node/impl"
	"github.com/downflux/go-bvh/internal/cache/node/util/cmp"
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/internal/cache/id"
)

var (
	_ C = Catto
)

func TestCatto(t *testing.T) {
	type config struct {
		name string
		c    *cache.C
		n    node.N
		aabb hyperrectangle.R
		want node.N
	}

	configs := []config{
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			root := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			root.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}))
			root.SetHeuristic(heuristic.H(root.AABB().R()))

			return config{
				name: "Root",
				c:    c,
				n:    root,
				aabb: *hyperrectangle.New(vector.V{1, 1}, vector.V{2, 1}),
				want: root,
			}
		}(),
		func() config {
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

			na.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{10, 1}))

			nb.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{2, 1}))
			nc.AABB().Copy(*hyperrectangle.New(vector.V{8, 0}, vector.V{10, 1}))

			na.SetHeuristic(heuristic.H(na.AABB().R()))
			nb.SetHeuristic(heuristic.H(nb.AABB().R()))
			nc.SetHeuristic(heuristic.H(nc.AABB().R()))

			want := impl.New(c, 4)
			want.Allocate(3, cid.IDInvalid, cid.IDInvalid)

			return config{
				name: "Expand",
				c:    c,
				n:    na,
				aabb: *hyperrectangle.New(vector.V{3, 0}, vector.V{4, 1}),
				want: want,
			}
		}(),
	}

	configs = append(configs, func() []config {
		c := cache.New(cache.O{
			LeafSize: 1,
			K:        2,
		})

		na := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

		nb := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))
		nc := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))

		nd := c.GetOrDie(c.Insert(nb.ID(), cid.IDInvalid, cid.IDInvalid, true))
		ne := c.GetOrDie(c.Insert(nb.ID(), cid.IDInvalid, cid.IDInvalid, true))

		nf := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))
		ng := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))

		na.SetLeft(nb.ID())
		na.SetRight(nc.ID())

		nb.SetLeft(nd.ID())
		nb.SetRight(ne.ID())

		nc.SetLeft(nf.ID())
		nc.SetRight(ng.ID())

		na.SetHeight(2)
		nb.SetHeight(1)
		nc.SetHeight(1)

		na.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{10, 1}))

		nb.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{3, 1}))
		nc.AABB().Copy(*hyperrectangle.New(vector.V{6, 0}, vector.V{10, 1}))

		nd.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}))
		ne.AABB().Copy(*hyperrectangle.New(vector.V{2, 0}, vector.V{3, 1}))

		nf.AABB().Copy(*hyperrectangle.New(vector.V{6, 0}, vector.V{8, 1}))
		ng.AABB().Copy(*hyperrectangle.New(vector.V{9, 0}, vector.V{10, 1}))

		na.SetHeuristic(heuristic.H(na.AABB().R()))
		nb.SetHeuristic(heuristic.H(nb.AABB().R()))
		nc.SetHeuristic(heuristic.H(nc.AABB().R()))
		nd.SetHeuristic(heuristic.H(nd.AABB().R()))
		ne.SetHeuristic(heuristic.H(ne.AABB().R()))
		nf.SetHeuristic(heuristic.H(nf.AABB().R()))
		ng.SetHeuristic(heuristic.H(ng.AABB().R()))

		return []config{
			{
				name: "Leaf/D",
				c:    c,
				n:    na,
				aabb: *hyperrectangle.New(vector.V{-1, 0}, vector.V{0, 1}),
				want: nd,
			},
			{
				name: "Leaf/E",
				c:    c,
				n:    na,
				aabb: *hyperrectangle.New(vector.V{3, 0}, vector.V{4, 1}),
				want: ne,
			},
			{
				name: "Leaf/F",
				c:    c,
				n:    na,
				aabb: *hyperrectangle.New(vector.V{5, 0}, vector.V{6, 1}),
				want: nf,
			},
			{
				name: "Leaf/G",
				c:    c,
				n:    na,
				aabb: *hyperrectangle.New(vector.V{10, 0}, vector.V{11, 1}),
				want: ng,
			},
		}
	}()...)

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := Catto(c.c, c.n, c.aabb); !cmp.Equal(got, c.want) {
				t.Errorf("Catto() = %v, want = %v", got, c.want)
			}
		})
	}
}
