package node

import (
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBroadPhase(t *testing.T) {
	type config struct {
		name string
		n    *N
		q    hyperrectangle.R
		want []id.ID
	}

	configs := []config{
		func() config {
			c := Cache()
			root := New(O{
				Nodes: c,
				ID:    100,
				Left:  101,
				Right: 102,
			})
			New(O{
				Nodes:  c,
				ID:     101,
				Parent: 100,
				Data: map[id.ID]hyperrectangle.R{
					1: *hyperrectangle.New([]float64{0}, []float64{24}),
					2: *hyperrectangle.New([]float64{25}, []float64{49}),
				},
			})
			New(O{
				Nodes:  c,
				ID:     102,
				Parent: 100,
				Data: map[id.ID]hyperrectangle.R{
					3: *hyperrectangle.New([]float64{50}, []float64{74}),
					4: *hyperrectangle.New([]float64{75}, []float64{99}),
				},
			})

			return config{
				name: "Internal",
				n:    root,
				q:    *hyperrectangle.New([]float64{26}, []float64{73}),
				want: []id.ID{2, 3},
			}
		}(),
		{
			name: "Leaf/Overlaps",
			n: New(O{
				Nodes: Cache(),
				Data: map[id.ID]hyperrectangle.R{
					1: *hyperrectangle.New([]float64{51}, []float64{100}),
					2: *hyperrectangle.New([]float64{0}, []float64{50}),
				},
			}),
			q:    *hyperrectangle.New([]float64{1}, []float64{99}),
			want: []id.ID{1, 2},
		},
		{
			name: "Leaf/Disjoint",
			n: New(O{
				Nodes: Cache(),
				Data: map[id.ID]hyperrectangle.R{
					1: *hyperrectangle.New([]float64{51}, []float64{100}),
					2: *hyperrectangle.New([]float64{0}, []float64{50}),
				},
			}),
			q:    *hyperrectangle.New([]float64{100.1}, []float64{100.2}),
			want: []id.ID{},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := c.n.BroadPhase(c.q)
			if diff := cmp.Diff(
				c.want,
				got,
				cmpopts.SortSlices(
					func(a, b id.ID) bool { return a < b },
				),
			); diff != "" {
				t.Errorf("BroadPhase() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
