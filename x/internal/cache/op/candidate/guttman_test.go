package candidate

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

var (
	_ C = Guttman
)

func TestGuttman(t *testing.T) {
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
				aabb: *hyperrectangle.New(vector.V{2, 2}, vector.V{3, 3}),
				want: root,
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
		left.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}))
		right.AABB().Copy(*hyperrectangle.New(vector.V{8, 0}, vector.V{10, 1}))

		return []config{
			config{
				name: "Left",
				c:    c,
				n:    root,
				aabb: *hyperrectangle.New(vector.V{1, 0}, vector.V{2, 1}),
				want: left,
			},
			// Check that in the case that if the proposed heuristic
			// delta is increased by the same amount between the
			// left and right nodes, the node with the lesser total
			// heuristic after the merge is chosen.
			config{
				name: "Left/DeltaTie",
				c:    c,
				n:    root,
				aabb: *hyperrectangle.New(vector.V{4, 0}, vector.V{5, 1}),
				want: left,
			},
			config{
				name: "Right",
				c:    c,
				n:    root,
				aabb: *hyperrectangle.New(vector.V{7, 0}, vector.V{8, 1}),
				want: right,
			},
		}
	}()...)

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := Guttman(c.c, c.n, c.aabb); !node.Equal(got, c.want) {
				t.Errorf("Guttman() = %v, want = %v", got, c.want)
			}
		})
	}
}
