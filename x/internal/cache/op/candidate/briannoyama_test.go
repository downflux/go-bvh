package candidate

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/util/cmp"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

var (
	_ C = BrianNoyama
)

func TestBrianNoyama(t *testing.T) {
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

			return config{
				name: "Root",
				c:    c,
				n:    root,
				aabb: *hyperrectangle.New(vector.V{100, 100}, vector.V{101, 101}),
				want: root,
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

		na.SetHeight(1)

		na.SetLeft(nb.ID())
		na.SetRight(nc.ID())

		na.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{10, 1}))

		nb.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{3, 1}))
		nc.AABB().Copy(*hyperrectangle.New(vector.V{7, 0}, vector.V{10, 1}))

		return []config{
			{
				name: "Left",
				c:    c,
				n:    na,
				aabb: *hyperrectangle.New(vector.V{-1, 0}, vector.V{0, 1}),
				want: nb,
			},
			{
				name: "Right",
				c:    c,
				n:    na,
				aabb: *hyperrectangle.New(vector.V{10, 0}, vector.V{11, 1}),
				want: nc,
			},
		}
	}()...)

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := BrianNoyama(c.c, c.n, c.aabb); !cmp.Equal(got, c.want) {
				t.Errorf("BrianNoyama() = %v, want = %v", got, c.want)
			}
		})
	}
}
