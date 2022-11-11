package tree

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

func TestSibling(t *testing.T) {
	type config struct {
		name string
		c    *cache.C
		x    cid.ID
		aabb hyperrectangle.R
		want cid.ID
	}

	configs := []config{
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})
			root := c.GetOrDie(
				c.Insert(
					cid.IDInvalid,
					cid.IDInvalid,
					cid.IDInvalid,
					/* validate = */ true,
				),
			)

			root.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{10, 10})))
			return config{
				name: "SingleNode/Contained",
				c:    c,
				x:    root.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{1, 1}), vector.V([]float64{2, 2})),
				want: root.ID(),
			}
		}(),
		// Ensure that we will always get a sibling node as long as the
		// cache is not empty.
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})
			root := c.GetOrDie(
				c.Insert(
					cid.IDInvalid,
					cid.IDInvalid,
					cid.IDInvalid,
					/* validate = */ true,
				),
			)

			root.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{10, 10})))
			return config{
				name: "SingleNode/NotContained",
				c:    c,
				x:    root.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{-1, -1}), vector.V([]float64{100, 100})),
				want: root.ID(),
			}
		}(),
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})
			root := c.GetOrDie(
				c.Insert(
					cid.IDInvalid,
					cid.IDInvalid,
					cid.IDInvalid,
					/* validate = */ true,
				),
			)

			root.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{10, 10})))
			return config{
				name: "SingleNode/Overlap",
				c:    c,
				x:    root.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{1, 1}), vector.V([]float64{100, 100})),
				want: root.ID(),
			}
		}(),
		// Check that the child node which contains the input AABB
		// should be the optimal sibling.
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})
			root := c.GetOrDie(
				c.Insert(
					cid.IDInvalid,
					cid.IDInvalid,
					cid.IDInvalid,
					true,
				),
			)
			left := c.GetOrDie(
				c.Insert(
					root.ID(),
					cid.IDInvalid,
					cid.IDInvalid,
					true,
				),
			)
			root.SetLeft(left.ID())
			right := c.GetOrDie(
				c.Insert(
					root.ID(),
					cid.IDInvalid,
					cid.IDInvalid,
					true,
				),
			)
			root.SetRight(right.ID())

			root.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{100, 100})))
			left.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{10, 10})))
			right.AABB().Copy(*hyperrectangle.New(vector.V([]float64{50, 50}), vector.V([]float64{100, 100})))

			return config{
				name: "ChildNode/Contains",
				c:    c,
				x:    root.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{1, 1}), vector.V([]float64{9, 9})),
				want: left.ID(),
			}
		}(),
		// Check that the child node which minimizes the total AABB SAH
		// change should be the optimal sibling.
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})
			root := c.GetOrDie(
				c.Insert(
					cid.IDInvalid,
					cid.IDInvalid,
					cid.IDInvalid,
					true,
				),
			)
			left := c.GetOrDie(
				c.Insert(
					root.ID(),
					cid.IDInvalid,
					cid.IDInvalid,
					true,
				),
			)
			root.SetLeft(left.ID())
			right := c.GetOrDie(
				c.Insert(
					root.ID(),
					cid.IDInvalid,
					cid.IDInvalid,
					true,
				),
			)
			root.SetRight(right.ID())

			root.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{100, 100})))
			left.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{50, 50})))
			right.AABB().Copy(*hyperrectangle.New(vector.V([]float64{90, 90}), vector.V([]float64{100, 100})))

			// Check that the direct cost of node construction is
			// accounted for. That is, we select the node with the
			// minimum SAH after merging with the input AABB.
			return config{
				name: "ChildNode/Overlaps",
				c:    c,
				x:    root.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{89, 89}), vector.V([]float64{91, 91})),
				want: right.ID(),
			}
		}(),
		// Assert that the total size of the candidate node after a
		// merge event does not actually matter.
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})
			root := c.GetOrDie(
				c.Insert(
					cid.IDInvalid,
					cid.IDInvalid,
					cid.IDInvalid,
					true,
				),
			)
			left := c.GetOrDie(
				c.Insert(
					root.ID(),
					cid.IDInvalid,
					cid.IDInvalid,
					true,
				),
			)
			root.SetLeft(left.ID())
			right := c.GetOrDie(
				c.Insert(
					root.ID(),
					cid.IDInvalid,
					cid.IDInvalid,
					true,
				),
			)
			root.SetRight(right.ID())

			root.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{100, 100})))
			left.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{50, 50})))
			right.AABB().Copy(*hyperrectangle.New(vector.V([]float64{90, 90}), vector.V([]float64{100, 100})))

			// Check that the direct cost of node construction is
			// accounted for. That is, we select the node with the
			// minimum SAH after merging with the input AABB.
			return config{
				name: "ChildNode/Overlaps/NoPreference",
				c:    c,
				x:    root.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{49, 49}), vector.V([]float64{51, 51})),
				want: left.ID(),
			}
		}(),
		// Assert that the induced cost matters is accounted for within
		// inner nodes.
		//
		//      A
		//     / \
		//    /   \
		//   B     C
		//  / \   / \
		// D   E F   G
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})
			nA := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

			nB := c.GetOrDie(c.Insert(nA.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nC := c.GetOrDie(c.Insert(nA.ID(), cid.IDInvalid, cid.IDInvalid, true))

			nD := c.GetOrDie(c.Insert(nB.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nE := c.GetOrDie(c.Insert(nB.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nF := c.GetOrDie(c.Insert(nC.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nG := c.GetOrDie(c.Insert(nC.ID(), cid.IDInvalid, cid.IDInvalid, true))

			nA.SetLeft(nB.ID())
			nA.SetRight(nC.ID())
			nB.SetLeft(nD.ID())
			nB.SetRight(nE.ID())
			nC.SetLeft(nF.ID())
			nC.SetRight(nG.ID())

			nA.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{100, 100})))

			nB.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{50, 50})))
			nC.AABB().Copy(*hyperrectangle.New(vector.V([]float64{90, 90}), vector.V([]float64{100, 100})))

			nD.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{10, 10})))
			nE.AABB().Copy(*hyperrectangle.New(vector.V([]float64{40, 40}), vector.V([]float64{50, 50})))
			nF.AABB().Copy(*hyperrectangle.New(vector.V([]float64{90, 90}), vector.V([]float64{91, 91})))
			nG.AABB().Copy(*hyperrectangle.New(vector.V([]float64{95, 95}), vector.V([]float64{100, 100})))

			// Check that the induced cost is accounted for by inner
			// tree nodes.
			return config{
				name: "InnerNode/Overlaps",
				c:    c,
				x:    nA.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{49, 49}), vector.V([]float64{51, 51})),
				want: nE.ID(),
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
