package swap

import (
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/internal/node/id"
)

func TestIsAncestor(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		m    *node.N
		want bool
	}

	configs := []config{
		func() config {
			root := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: util.Interval(0, 100)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{},
				},
				Root: 100,
				Size: 1,
			})
			return config{
				name: "TrivialRoot",
				n:    root,
				m:    root,
				want: true,
			}
		}(),
		func() config {
			root := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
					102: {2: util.Interval(101, 200)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 101, Right: 102},
					101: util.N{Parent: 100},
					102: util.N{Parent: 100},
				},
				Root: 100,
				Size: 1,
			})
			return config{
				name: "TrivialSibling",
				n:    root.Left(),
				m:    root.Right(),
				want: false,
			}
		}(),
		func() config {
			root := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
					102: {2: util.Interval(101, 200)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 101, Right: 102},
					101: util.N{Parent: 100},
					102: util.N{Parent: 100},
				},
				Root: 100,
				Size: 1,
			})
			return config{
				name: "ImmediateChild",
				n:    root.Left(),
				m:    root,
				want: false,
			}
		}(),
		func() config {
			root := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
					102: {2: util.Interval(101, 200)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 101, Right: 102},
					101: util.N{Parent: 100},
					102: util.N{Parent: 100},
				},
				Root: 100,
				Size: 1,
			})
			return config{
				name: "ImmediateParent",
				n:    root,
				m:    root.Left(),
				want: true,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := isAncestor(c.n, c.m); got != c.want {
				t.Errorf("isAncestor() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestSwap(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		m    *node.N
		want *node.N
	}

	configs := []config{
		func() config {
			root := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
					102: {2: util.Interval(101, 200)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 101, Right: 102},
					101: util.N{Parent: 100},
					102: util.N{Parent: 100},
				},
				Root: 100,
				Size: 1,
			})
			want := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
					102: {2: util.Interval(101, 200)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 102, Right: 101},
					101: util.N{Parent: 100},
					102: util.N{Parent: 100},
				},
				Root: 100,
				Size: 1,
			})
			return config{
				name: "Siblings",
				n:    root.Left(),
				m:    root.Right(),
				want: want,
			}
		}(),
		func() config {
			data := map[nid.ID]map[id.ID]hyperrectangle.R{
				101: {1: util.Interval(0, 100)},   // B
				103: {2: util.Interval(101, 200)}, // D
				104: {3: util.Interval(201, 300)}, // E
			}

			root := util.New(util.T{
				Data: data,
				//   A
				//  / \
				// B   C
				//    / \
				//   D   E
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 101, Right: 102},              // A
					101: util.N{Parent: 100},                        // B
					102: util.N{Left: 103, Right: 104, Parent: 100}, // C
					103: util.N{Parent: 102},                        // D
					104: util.N{Parent: 102},                        // E
				},
				Root: 100,
				Size: 1,
			})
			want := util.New(util.T{
				Data: data,
				//   A
				//  / \
				// D   C
				//    / \
				//   B   E
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 103, Right: 102},              // A
					101: util.N{Parent: 102},                        // B
					102: util.N{Left: 101, Right: 104, Parent: 100}, // C
					103: util.N{Parent: 100},                        // D
					104: util.N{Parent: 102},                        // E
				},
				Root: 100,
				Size: 1,
			})
			return config{
				name: "Internal/Leaf",
				n:    root.Left(),         // B
				m:    root.Right().Left(), // D
				want: want,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			Execute(c.n, c.m)
			if diff := cmp.Diff(
				c.want,
				c.n.Root(),
				cmp.Comparer(util.Equal),
			); diff != "" {
				t.Errorf("Execute() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
