package util

import (
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/internal/node/id"
)

func TestEqual(t *testing.T) {
	type config struct {
		name string
		a    *node.N
		b    *node.N
		want bool
	}

	configs := []config{
		{
			name: "Trivial",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "Leaf/NoLeaf",
			a: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: Interval(0, 100)},
				},
				Nodes: map[nid.ID]N{
					100: N{},
				},
				Root: 100,
			}),
			b: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: Interval(0, 100)},
					102: {2: Interval(101, 200)},
				},
				Nodes: map[nid.ID]N{
					100: N{Left: 101, Right: 102},
					101: N{Parent: 100},
					102: N{Parent: 100},
				},
				Root: 100,
			}),
			want: false,
		},
		{
			name: "Leaf/NoEqual/Data",
			a: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: Interval(0, 100)},
				},
				Nodes: map[nid.ID]N{
					100: N{},
				},
				Root: 100,
			}),
			b: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {
						1: Interval(0, 100),
						2: Interval(101, 200),
					},
				},
				Nodes: map[nid.ID]N{
					100: N{},
				},
				Root: 100,
			}),
			want: false,
		},
		{
			name: "Leaf/NoEqual",
			a: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: Interval(0, 100)},
				},
				Nodes: map[nid.ID]N{
					100: N{},
				},
				Root: 100,
			}),
			b: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {2: Interval(0, 100)},
				},
				Nodes: map[nid.ID]N{
					100: N{},
				},
				Root: 100,
			}),
			want: false,
		},
		{
			name: "Leaf/NIDInvariant",
			a: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: Interval(0, 100)},
				},
				Nodes: map[nid.ID]N{
					100: N{},
				},
				Root: 100,
			}),
			b: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					1001: {1: Interval(0, 100)},
				},
				Nodes: map[nid.ID]N{
					1001: N{},
				},
				Root: 1001,
			}),
			want: true,
		},
		{
			name: "Internal/BadLeaf",
			a: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: Interval(0, 100)},
					102: {2: Interval(101, 200)},
				},
				Nodes: map[nid.ID]N{
					100: N{Left: 101, Right: 102},
					101: N{Parent: 100},
					102: N{Parent: 100},
				},
				Root: 100,
			}),
			b: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: Interval(0, 100)},
					102: {2: Interval(101, 200)},
				},
				Nodes: map[nid.ID]N{
					100: N{Left: 102, Right: 101},
					101: N{Parent: 100},
					102: N{Parent: 100},
				},
				Root: 100,
			}),
			want: false,
		},
		{
			name: "Internal/NIDInvariant",
			a: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: Interval(0, 100)},
					102: {2: Interval(101, 200)},
				},
				Nodes: map[nid.ID]N{
					100: N{Left: 101, Right: 102},
					101: N{Parent: 100},
					102: N{Parent: 100},
				},
				Root: 100,
			}),
			b: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					1001: {1: Interval(0, 100)},
					1002: {2: Interval(101, 200)},
				},
				Nodes: map[nid.ID]N{
					1000: N{Left: 1001, Right: 1002},
					1001: N{Parent: 1000},
					1002: N{Parent: 1000},
				},
				Root: 1000,
			}),
			want: true,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := Equal(c.a, c.b); got != c.want {
				t.Errorf("Equal() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type config struct {
		name string
		t    T
		want *node.N
	}

	configs := []config{
		{
			name: "Leaf",
			t: T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: Interval(0, 100)},
				},
				Nodes: map[nid.ID]N{
					100: {},
				},
				Root: 100,
			},
			want: node.New(node.O{
				Nodes: node.Cache(),
				ID:    100,
				Data: map[id.ID]hyperrectangle.R{
					1: Interval(0, 100),
				},
			}),
		}, func() config {
			c := node.Cache()
			r := node.New(node.O{
				Nodes: c,
				ID:    100,
				Left:  101,
				Right: 102,
			})
			node.New(node.O{
				Nodes:  c,
				ID:     101,
				Parent: 100,
				Data: map[id.ID]hyperrectangle.R{
					1: Interval(0, 100),
				},
			})
			node.New(node.O{
				Nodes:  c,
				ID:     102,
				Parent: 100,
				Data: map[id.ID]hyperrectangle.R{
					2: Interval(101, 200),
				},
			})

			return config{
				name: "Root",
				t: T{
					Data: map[nid.ID]map[id.ID]hyperrectangle.R{
						101: {1: Interval(0, 100)},
						102: {2: Interval(101, 200)},
					},
					Nodes: map[nid.ID]N{
						100: {Left: 101, Right: 102},
						101: {Parent: 100},
						102: {Parent: 100},
					},
					Root: 100,
				},
				want: r,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := New(c.t)
			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(
					node.N{},
					node.C{},
					hyperrectangle.R{},
				),
			); diff != "" {
				t.Errorf("New() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}