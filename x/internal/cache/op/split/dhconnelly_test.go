package split

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

var (
	_ S = DHConnelly
)

func TestNext(t *testing.T) {
	type config struct {
		name   string
		data   map[id.ID]hyperrectangle.R
		leaves []id.ID
		n      node.N
		m      node.N
		want   int
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

			left.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}))
			right.AABB().Copy(*hyperrectangle.New(vector.V{10, 0}, vector.V{11, 1}))

			// Leaf 101 is directly adjacent to the left node, which
			// means the induced cost is much cheaper than if this
			// leaf were to be merged into the right node instead.
			// This is the crucial heuristic used by the next()
			// function.
			return config{
				name: "Adjacent",
				data: map[id.ID]hyperrectangle.R{
					100: *hyperrectangle.New(vector.V{5, 0}, vector.V{6, 1}),
					101: *hyperrectangle.New(vector.V{1, 0}, vector.V{2, 1}),
				},
				leaves: []id.ID{100, 101},
				n:      left,
				m:      right,
				want:   1,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			buf := hyperrectangle.New(vector.V(make([]float64, 2)), vector.V(make([]float64, 2))).M()
			if got := next(c.data, c.leaves, c.n, c.m, buf); got != c.want {
				t.Errorf("next() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestSeed(t *testing.T) {
	type w struct {
		l int
		r int
	}

	type config struct {
		name   string
		data   map[id.ID]hyperrectangle.R
		leaves []id.ID
		want   w
	}

	configs := []config{
		{
			name: "Trivial",
			data: map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
				101: *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
			},
			leaves: []id.ID{100, 101},
			want: w{
				l: 0,
				r: 1,
			},
		},
	}
	configs = append(configs, func() []config {
		data := map[id.ID]hyperrectangle.R{
			100: *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
			101: *hyperrectangle.New(vector.V{1, 0}, vector.V{2, 1}),
			102: *hyperrectangle.New(vector.V{10, 0}, vector.V{10, 1}),
		}
		return []config{
			// The largest waste of space in the following scenario is a box
			// drawn around 100 and 102. Check that this is handled
			// appropriately.
			{
				name:   "LargeBox",
				data:   data,
				leaves: []id.ID{100, 101, 102},
				want: w{
					l: 0,
					r: 2,
				},
			},
			// Check that the order of input leaves does not matter.
			{
				name:   "LargeBox/OrderInvariant",
				data:   data,
				leaves: []id.ID{101, 102, 100},
				want: w{
					l: 1,
					r: 2,
				},
			},
		}
	}()...)

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := &w{}
			got.l, got.r = seed(c.data, c.leaves, hyperrectangle.New(
				vector.V(make([]float64, 2)),
				vector.V(make([]float64, 2)),
			).M())

			if got.l != c.want.l {
				t.Errorf("l = %v, want = %v", got.l, c.want.l)
			}
			if got.r != c.want.r {
				t.Errorf("r = %v, want = %v", got.r, c.want.r)
			}
		})
	}
}
