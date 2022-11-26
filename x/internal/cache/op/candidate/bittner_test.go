package candidate

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/impl"
	"github.com/downflux/go-bvh/x/internal/cache/node/util/cmp"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

var (
	_ C = Bittner
)

func TestBittner(t *testing.T) {
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
			left := c.GetOrDie(c.Insert(root.ID(), cid.IDInvalid, cid.IDInvalid, true))
			right := c.GetOrDie(c.Insert(root.ID(), cid.IDInvalid, cid.IDInvalid, true))

			root.SetLeft(left.ID())
			root.SetRight(right.ID())

			root.SetHeight(1)

			root.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{10, 1}))
			left.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{3, 1}))
			right.AABB().Copy(*hyperrectangle.New(vector.V{7, 0}, vector.V{10, 1}))

			root.SetHeuristic(heuristic.H(root.AABB().R()))
			left.SetHeuristic(heuristic.H(left.AABB().R()))
			right.SetHeuristic(heuristic.H(right.AABB().R()))

			want := impl.New(c, 4)
			want.Allocate(3, cid.IDInvalid, cid.IDInvalid)

			return config{
				name: "Expand",
				c:    c,
				n:    root,
				aabb: *hyperrectangle.New(vector.V{0, 2}, vector.V{10, 3}),
				want: want,
			}
		}(),
	}

	configs = append(configs, func() []config {
		c := cache.New(cache.O{
			LeafSize: 1,
			K:        2,
		})

		root := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
		left := c.GetOrDie(c.Insert(root.ID(), cid.IDInvalid, cid.IDInvalid, true))
		right := c.GetOrDie(c.Insert(root.ID(), cid.IDInvalid, cid.IDInvalid, true))

		root.SetLeft(left.ID())
		root.SetRight(right.ID())

		root.SetHeight(1)

		root.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{10, 1}))
		left.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{3, 1}))
		right.AABB().Copy(*hyperrectangle.New(vector.V{7, 0}, vector.V{10, 1}))

		root.SetHeuristic(heuristic.H(root.AABB().R()))
		left.SetHeuristic(heuristic.H(left.AABB().R()))
		right.SetHeuristic(heuristic.H(right.AABB().R()))

		return []config{
			{
				name: "Leaf/Left",
				c:    c,
				n:    root,
				aabb: *hyperrectangle.New(vector.V{-1, 0}, vector.V{0, 1}),
				want: left,
			},
			{
				name: "Leaf/Right",
				c:    c,
				n:    root,
				aabb: *hyperrectangle.New(vector.V{10, 0}, vector.V{11, 1}),
				want: right,
			},
		}
	}()...)

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := Bittner(c.c, c.n, c.aabb); !cmp.Equal(got, c.want) {
				t.Errorf("Bittner() = %v, want = %v", got, c.want)
			}
		})
	}
}
