package insert

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/x/internal/node/id"
)

func TestExecute(t *testing.T) {
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
				},
				Nodes: map[nid.ID]util.N{
					101: util.N{},
				},
				Root: 101,
			})
			m := node.New(node.O{
				Nodes: root.Cache(),
				Data:  map[id.ID]hyperrectangle.R{2: util.Interval(101, 200)},
			})

			return config{
				name: "Root",
				n:    root,
				m:    m,
				want: util.New(util.T{
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
				}),
			}
		}(),
		func() config {
			root := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},   // B
					102: {2: util.Interval(101, 200)}, // C
				},
				//   A
				//  / \
				// B   C
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 101, Right: 102}, // A
					101: util.N{Parent: 100},           // B
					102: util.N{Parent: 100},           // C
				},
				Root: 100,
			})
			m := node.New(node.O{
				Nodes: root.Cache(),
				Data:  map[id.ID]hyperrectangle.R{3: util.Interval(201, 300)},
			})

			return config{
				name: "Sibling",
				n:    root.Right(), // C
				m:    m,
				//   A
				//  / \
				// B   D
				//    / \
				//   E   C
				want: util.New(util.T{
					Data: map[nid.ID]map[id.ID]hyperrectangle.R{
						101: {1: util.Interval(0, 100)},   // B
						103: {2: util.Interval(101, 200)}, // C
						104: {3: util.Interval(201, 300)}, // E
					},
					Nodes: map[nid.ID]util.N{
						100: util.N{Left: 101, Right: 102},              // A
						101: util.N{Parent: 100},                        // B
						102: util.N{Left: 104, Right: 103, Parent: 100}, // D
						103: util.N{Parent: 102},                        // C
						104: util.N{Parent: 102},                        // E
					},
					Root: 100,
				}),
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Execute(c.n, c.m)
			if diff := cmp.Diff(c.want, got, cmp.Comparer(util.Equal)); diff != "" {
				t.Errorf("Execute() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
