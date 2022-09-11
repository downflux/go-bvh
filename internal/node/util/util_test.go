package util

import (
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-geometry/epsilon"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/internal/node/id"
)

func TestCost(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		want float64
	}

	configs := []config{
		{
			name: "Leaf",
			n: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {
						1: *hyperrectangle.New(
							[]float64{0, 0},
							[]float64{10, 10},
						),
					},
				},
				Nodes: map[nid.ID]N{
					100: N{},
				},
				Root: 100,
				Size: 1,
			}),
			want: 40,
		},
		{
			name: "Recursive",
			n: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {
						1: *hyperrectangle.New(
							[]float64{0, 0},
							[]float64{10, 10},
						),
					},
					103: {
						2: *hyperrectangle.New(
							[]float64{11, 11},
							[]float64{20, 20},
						),
					},
					104: {
						3: *hyperrectangle.New(
							[]float64{21, 21},
							[]float64{30, 30},
						),
					},
				},
				Nodes: map[nid.ID]N{
					100: N{Left: 101, Right: 102},
					101: N{Parent: 100},
					102: N{Left: 103, Right: 104},
					103: N{Parent: 102},
					104: N{Parent: 102},
				},
				Root: 100,
				Size: 1,
			}),
			want: 120 + 40 + 76 + 36 + 36,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := Cost(c.n); !epsilon.Within(c.want, got) {
				t.Errorf("Cost() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestMaxImbalance(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		want uint
	}

	configs := []config{
		{
			name: "Leaf",
			n: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: Interval(0, 100)},
				},
				Nodes: map[nid.ID]N{
					100: N{},
				},
				Root: 100,
				Size: 1,
			}),
			want: 0,
		},
		{
			name: "Recursive",
			n: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: Interval(0, 100)},
					103: {2: Interval(101, 200)},
					104: {3: Interval(201, 300)},
				},
				Nodes: map[nid.ID]N{
					100: N{Left: 101, Right: 102},
					101: N{Parent: 100},
					102: N{Left: 103, Right: 104},
					103: N{Parent: 102},
					104: N{Parent: 102},
				},
				Root: 100,
				Size: 1,
			}),
			want: 1,
		},
		{
			name: "Recursive/Commutative",
			n: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: Interval(0, 100)},
					103: {2: Interval(101, 200)},
					104: {3: Interval(201, 300)},
				},
				Nodes: map[nid.ID]N{
					100: N{Right: 101, Left: 102},
					101: N{Parent: 100},
					102: N{Left: 103, Right: 104},
					103: N{Parent: 102},
					104: N{Parent: 102},
				},
				Root: 100,
				Size: 1,
			}),
			want: 1,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := MaxImbalance(c.n); got != c.want {
				t.Errorf("MaxImbalance() = %v, want = %v", got, c.want)
			}
		})
	}
}

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
				Size: 1,
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
				Size: 1,
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
				Size: 2,
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
				Size: 2,
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
				Size: 1,
			}),
			b: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {2: Interval(0, 100)},
				},
				Nodes: map[nid.ID]N{
					100: N{},
				},
				Root: 100,
				Size: 1,
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
				Size: 1,
			}),
			b: New(T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					1001: {1: Interval(0, 100)},
				},
				Nodes: map[nid.ID]N{
					1001: N{},
				},
				Root: 1001,
				Size: 1,
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
				Size: 1,
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
				Size: 1,
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
				Size: 1,
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
				Size: 1,
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
				Size: 1,
			},
			want: node.New(node.O{
				Nodes: node.Cache(),
				ID:    100,
				Data: map[id.ID]hyperrectangle.R{
					1: Interval(0, 100),
				},
				Size: 1,
				K:    1,
			}),
		}, func() config {
			c := node.Cache()
			r := node.New(node.O{
				Nodes: c,
				ID:    100,
				Left:  101,
				Right: 102,
				Size:  1,
				K:     1,
			})
			node.New(node.O{
				Nodes:  c,
				ID:     101,
				Parent: 100,
				Data: map[id.ID]hyperrectangle.R{
					1: Interval(0, 100),
				},
				Size: 1,
				K:    1,
			})
			node.New(node.O{
				Nodes:  c,
				ID:     102,
				Parent: 100,
				Data: map[id.ID]hyperrectangle.R{
					2: Interval(101, 200),
				},
				Size: 1,
				K:    1,
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
					Size: 1,
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
