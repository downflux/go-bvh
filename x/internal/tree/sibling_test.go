package tree

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

func TestSibling(t *testing.T) {
	type config struct {
		name string
		c    *cache.C
		x    cache.ID
		aabb hyperrectangle.R
		want cache.ID
	}

	configs := []config{
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        1,
			})
			root := c.GetOrDie(
				c.Insert(
					cache.IDInvalid,
					cache.IDInvalid,
					cache.IDInvalid,
					/* validate = */ true,
				),
			)

			root.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0}), vector.V([]float64{10})))
			return config{
				name: "SingleNode/Contained",
				c:    c,
				x:    root.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{1}), vector.V([]float64{2})),
				want: root.ID(),
			}
		}(),
		// Ensure that we will always get a sibling node as long as the
		// cache is not empty.
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        1,
			})
			root := c.GetOrDie(
				c.Insert(
					cache.IDInvalid,
					cache.IDInvalid,
					cache.IDInvalid,
					/* validate = */ true,
				),
			)

			root.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0}), vector.V([]float64{10})))
			return config{
				name: "SingleNode/NotContained",
				c:    c,
				x:    root.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{-1}), vector.V([]float64{100})),
				want: root.ID(),
			}
		}(),
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        1,
			})
			root := c.GetOrDie(
				c.Insert(
					cache.IDInvalid,
					cache.IDInvalid,
					cache.IDInvalid,
					/* validate = */ true,
				),
			)

			root.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0}), vector.V([]float64{10})))
			return config{
				name: "SingleNode/Overlap",
				c:    c,
				x:    root.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{1}), vector.V([]float64{100})),
				want: root.ID(),
			}
		}(),
		// Check that the child node which contains the input AABB
		// should be the optimal sibling.
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        1,
			})
			root := c.GetOrDie(
				c.Insert(
					cache.IDInvalid,
					cache.IDInvalid,
					cache.IDInvalid,
					true,
				),
			)
			left := c.GetOrDie(
				c.Insert(
					root.ID(),
					cache.IDInvalid,
					cache.IDInvalid,
					true,
				),
			)
			root.SetLeft(left.ID())
			right := c.GetOrDie(
				c.Insert(
					root.ID(),
					cache.IDInvalid,
					cache.IDInvalid,
					true,
				),
			)
			root.SetRight(right.ID())

			root.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0}), vector.V([]float64{100})))
			left.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0}), vector.V([]float64{10})))
			right.AABB().Copy(*hyperrectangle.New(vector.V([]float64{50}), vector.V([]float64{100})))

			return config{
				name: "ChildNode/Contains",
				c:    c,
				x:    root.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{1}), vector.V([]float64{9})),
				want: left.ID(),
			}
		}(),
		// Check that the child node which minimizes the total AABB SAH
		// change should be the optimal sibling.
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        1,
			})
			root := c.GetOrDie(
				c.Insert(
					cache.IDInvalid,
					cache.IDInvalid,
					cache.IDInvalid,
					true,
				),
			)
			left := c.GetOrDie(
				c.Insert(
					root.ID(),
					cache.IDInvalid,
					cache.IDInvalid,
					true,
				),
			)
			root.SetLeft(left.ID())
			right := c.GetOrDie(
				c.Insert(
					root.ID(),
					cache.IDInvalid,
					cache.IDInvalid,
					true,
				),
			)
			root.SetRight(right.ID())

			root.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0}), vector.V([]float64{100})))
			left.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0}), vector.V([]float64{10})))
			right.AABB().Copy(*hyperrectangle.New(vector.V([]float64{50}), vector.V([]float64{100})))

			// The child node which contains the input AABB should
			// be the optimal sibling.
			return config{
				name: "ChildNode/Overlaps",
				c:    c,
				x:    root.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{49}), vector.V([]float64{51})),
				want: right.ID(),
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := sibling(c.c, c.x, c.aabb); got != c.want {
				t.Errorf("sibling() = %v, want = %v", got, c.want)
			}
		})
	}
}
