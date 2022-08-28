package util

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
)

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
				Data: map[NodeID]map[id.ID]hyperrectangle.R{
					101: {1: Interval(0, 100)},
				},
				Nodes: map[NodeID]N{},
				Root:  101,
			},
			want: node.New(node.O{
				Data: map[id.ID]hyperrectangle.R{
					1: Interval(0, 100),
				},
			}),
		},
		{
			name: "Children",
			t: T{
				Data: map[NodeID]map[id.ID]hyperrectangle.R{
					101: {1: Interval(0, 100)},
					102: {2: Interval(101, 200)},
				},
				Nodes: map[NodeID]N{
					100: N{Left: 101, Right: 102},
				},
				Root: 100,
			},
			want: node.New(node.O{
				Left: node.New(node.O{
					Data: map[id.ID]hyperrectangle.R{
						1: Interval(0, 100),
					},
				}),
				Right: node.New(node.O{
					Data: map[id.ID]hyperrectangle.R{
						2: Interval(101, 200),
					},
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
			data := map[NodeID]map[id.ID]hyperrectangle.R{
				101: {1: Interval(0, 100)},
				103: {2: Interval(101, 200)},
				104: {3: Interval(201, 300)},
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
