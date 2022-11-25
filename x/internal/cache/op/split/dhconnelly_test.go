package split

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

var (
	_ S = DHConnelly
)

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
