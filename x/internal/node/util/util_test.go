package util

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
)

func interval(min, max float64) hyperrectangle.R {
	return *hyperrectangle.New([]float64{min}, []float64{max})
}

func TestNew(t *testing.T) {
	type config struct {
		name string
		t    T
		want *node.N
	}

	configs := []config{
		{
			name: "Trivial",
			t: T{
				Data: map[NodeID][]node.D{
					101: []node.D{{ID: 1, AABB: interval(0, 100)}},
				},
				Nodes: map[NodeID]N{},
				Root:  101,
			},
			want: node.New(node.O{
				Data: []node.D{{ID: 1, AABB: interval(0, 100)}},
			}),
		},
		{
			name: "Children",
			t: T{
				Data: map[NodeID][]node.D{
					101: []node.D{{ID: 1, AABB: interval(0, 100)}},
					102: []node.D{{ID: 2, AABB: interval(101, 200)}},
				},
				Nodes: map[NodeID]N{
					100: N{Left: 101, Right: 102},
				},
				Root: 100,
			},
			want: node.New(node.O{
				Left: node.New(node.O{
					Data: []node.D{{ID: 1, AABB: interval(0, 100)}},
				}),
				Right: node.New(node.O{
					Data: []node.D{{ID: 2, AABB: interval(101, 200)}},
				}),
			}),
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := New(c.t)
			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(
					node.N{},
					hyperrectangle.R{},
				),
			); diff != "" {
				t.Errorf("New() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestNodeSwap(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		m    *node.N
		want *node.N // root
	}

	configs := []config{
		func() config {
			data := map[NodeID][]node.D{
				101: []node.D{{ID: 1, AABB: interval(0, 100)}},
				103: []node.D{{ID: 2, AABB: interval(101, 200)}},
				104: []node.D{{ID: 3, AABB: interval(201, 300)}},
			}

			input := New(
				T{
					Data: data,
					Nodes: map[NodeID]N{
						100: N{Left: 101, Right: 102},
						102: N{Left: 103, Right: 104},
					},
					Root: 100,
				},
			)
			return config{
				name: "Simple",
				n:    input.Left(),
				m:    input.Right().Right(),
				want: New(
					T{
						Data: data,
						Nodes: map[NodeID]N{
							100: N{Left: 104, Right: 102},
							102: N{Left: 103, Right: 101},
						},
						Root: 100,
					},
				),
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			c.n.Swap(c.m)
			got := c.n.Root()

			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(
					node.N{},
					hyperrectangle.R{},
				),
			); diff != "" {
				t.Errorf("New() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
