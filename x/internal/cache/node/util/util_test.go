package util

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

func TestValidate(t *testing.T) {
	type config struct {
		name    string
		data    map[id.ID]hyperrectangle.R
		n       node.N
		success bool
	}

	configs := []config{
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
			}

			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			root := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			root.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}))
			root.Leaves()[100] = struct{}{}

			return config{
				name:    "Leaf",
				data:    data,
				n:       root,
				success: true,
			}
		}(),
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
			}

			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			root := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			root.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}))

			return config{
				name:    "Leaf/NoData",
				data:    data,
				n:       root,
				success: false,
			}
		}(),
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
			}

			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			root := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			root.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{0, 0}))
			root.Leaves()[100] = struct{}{}

			return config{
				name:    "Leaf/NoEncapsulate",
				data:    data,
				n:       root,
				success: false,
			}
		}(),
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

			na.SetHeight(1000)

			na.AABB().Copy(hyperrectangle.Union(data[100], data[101]))

			nb.AABB().Copy(data[100])
			nc.AABB().Copy(data[101])

			nb.Leaves()[100] = struct{}{}
			nc.Leaves()[101] = struct{}{}

			return config{
				name:    "HeightMismatch",
				data:    data,
				n:       na,
				success: false,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if err := Validate(c.data, c.n); c.success && err != nil {
				t.Errorf("Validate() encountered an unexpected error: %v", err)
			} else if !c.success && err == nil {
				t.Errorf("Validate() unexpectedly succeeded")
			}
		})
	}
}
