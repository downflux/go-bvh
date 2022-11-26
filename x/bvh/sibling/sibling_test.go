package sibling

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/util/cmp"
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
		want node.N
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
				want: root,
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
				want: root,
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
				want: root,
			}
		}(),
		// Check that the child node which contains the input AABB
		// should be the optimal sibling.
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})
			root := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			left := c.GetOrDie(c.Insert(root.ID(), cid.IDInvalid, cid.IDInvalid, true))
			right := c.GetOrDie(c.Insert(root.ID(), cid.IDInvalid, cid.IDInvalid, true))

			root.SetHeight(1)

			root.SetLeft(left.ID())
			root.SetRight(right.ID())

			root.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{100, 100})))
			left.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{10, 10})))
			right.AABB().Copy(*hyperrectangle.New(vector.V([]float64{50, 50}), vector.V([]float64{100, 100})))

			return config{
				name: "ChildNode/Contains",
				c:    c,
				x:    root.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{1, 1}), vector.V([]float64{9, 9})),
				want: left,
			}
		}(),
		// Check that the child node which minimizes the total AABB SAH
		// change should be the optimal sibling.
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			root := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			left := c.GetOrDie(c.Insert(root.ID(), cid.IDInvalid, cid.IDInvalid, true))
			right := c.GetOrDie(c.Insert(root.ID(), cid.IDInvalid, cid.IDInvalid, true))

			root.SetHeight(1)

			root.SetLeft(left.ID())
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
				want: right,
			}
		}(),
		// Assert that the total size of the candidate node after a
		// merge event does not actually matter.
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			root := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			left := c.GetOrDie(c.Insert(root.ID(), cid.IDInvalid, cid.IDInvalid, true))
			right := c.GetOrDie(c.Insert(root.ID(), cid.IDInvalid, cid.IDInvalid, true))

			root.SetHeight(1)

			root.SetLeft(left.ID())
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
				want: left,
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
			na := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

			nb := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nc := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))

			nd := c.GetOrDie(c.Insert(nb.ID(), cid.IDInvalid, cid.IDInvalid, true))
			ne := c.GetOrDie(c.Insert(nb.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nf := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))
			ng := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))

			na.SetHeight(2)
			nb.SetHeight(1)
			nc.SetHeight(1)

			na.SetLeft(nb.ID())
			na.SetRight(nc.ID())
			nb.SetLeft(nd.ID())
			nb.SetRight(ne.ID())
			nc.SetLeft(nf.ID())
			nc.SetRight(ng.ID())

			na.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{100, 100})))

			nb.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{50, 50})))
			nc.AABB().Copy(*hyperrectangle.New(vector.V([]float64{90, 90}), vector.V([]float64{100, 100})))

			nd.AABB().Copy(*hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{10, 10})))
			ne.AABB().Copy(*hyperrectangle.New(vector.V([]float64{40, 40}), vector.V([]float64{50, 50})))
			nf.AABB().Copy(*hyperrectangle.New(vector.V([]float64{90, 90}), vector.V([]float64{91, 91})))
			ng.AABB().Copy(*hyperrectangle.New(vector.V([]float64{95, 95}), vector.V([]float64{100, 100})))

			// Check that the induced cost is accounted for by inner
			// tree nodes.
			return config{
				name: "InnerNode/Overlaps",
				c:    c,
				x:    na.ID(),
				aabb: *hyperrectangle.New(vector.V([]float64{49, 49}), vector.V([]float64{51, 51})),
				want: ne,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := Find(c.c, c.x, c.aabb); !cmp.Equal(got, c.want) {
				t.Errorf("sibling() = %v, want = %v", got, c.want)
			}
		})
	}
}
