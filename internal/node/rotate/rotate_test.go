package rotate

import (
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/internal/node/id"
)

func TestRebalance(t *testing.T) {
	type config struct {
		name string
		z    *node.N
		want *node.N
	}

	configs := []config{
		{
			name: "Root",
			z: util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: util.Interval(0, 100)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{},
				},
				Root: 100,
				Size: 1,
			}),
			want: nil,
		},
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
				name: "Leaf",
				z:    root.Left(),
				want: root,
			}
		}(),
		func() config {
			root := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
					102: {2: util.Interval(101, 200)},
					103: {3: util.Interval(201, 300)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 104, Right: 103},
					101: util.N{Parent: 104},
					102: util.N{Parent: 104},
					103: util.N{Parent: 100},
					104: util.N{Left: 101, Right: 102, Parent: 100},
				},
				Root: 100,
				Size: 1,
			})
			want := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
					102: {2: util.Interval(101, 200)},
					103: {3: util.Interval(201, 300)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 104, Right: 103},
					101: util.N{Parent: 104},
					102: util.N{Parent: 104},
					103: util.N{Parent: 100},
					104: util.N{Left: 101, Right: 102, Parent: 100},
				},
				Root: 100,
				Size: 1,
			})
			return config{
				name: "NoSwap/ZLeft/YLeft",
				z:    root.Left().Left(),
				want: want.Left(),
			}
		}(),
		func() config {
			root := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
					102: {2: util.Interval(101, 200)},
					103: {3: util.Interval(201, 300)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 103, Right: 104},
					101: util.N{Parent: 104},
					102: util.N{Parent: 104},
					103: util.N{Parent: 100},
					104: util.N{Left: 101, Right: 102, Parent: 100},
				},
				Root: 100,
				Size: 1,
			})
			want := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
					102: {2: util.Interval(101, 200)},
					103: {3: util.Interval(201, 300)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 103, Right: 104},
					101: util.N{Parent: 104},
					102: util.N{Parent: 104},
					103: util.N{Parent: 100},
					104: util.N{Left: 101, Right: 102, Parent: 100},
				},
				Root: 100,
				Size: 1,
			})
			return config{
				name: "NoSwap/ZRight/YRight",
				z:    root.Right().Right(),
				want: want.Right(),
			}
		}(),
		func() config {
			data := map[nid.ID]map[id.ID]hyperrectangle.R{
				101: {1: util.Interval(0, 100)},   // B
				102: {2: util.Interval(101, 200)}, // D
				103: {3: util.Interval(201, 300)}, // F
				104: {4: util.Interval(301, 400)}, // G
			}

			//    A
			//   / \
			//  B   C
			//     / \
			//    D   E
			//       / \
			//      F   G
			root := util.New(util.T{
				Data: data,
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 101, Right: 105},
					101: util.N{Parent: 100},
					102: util.N{Parent: 105},
					103: util.N{Parent: 106},
					104: util.N{Parent: 106},
					105: util.N{Left: 102, Right: 106, Parent: 100}, // C
					106: util.N{Left: 103, Right: 104, Parent: 105}, // E
				},
				Root: 100,
				Size: 1,
			})

			//       A
			//      / \
			//     /   \
			//    E     C
			//   / \   / \
			//  F   G D   B
			want := util.New(util.T{
				Data: data,
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 105, Right: 106},
					101: util.N{Parent: 106},
					102: util.N{Parent: 106},
					103: util.N{Parent: 105},
					104: util.N{Parent: 105},
					105: util.N{Left: 103, Right: 104, Parent: 100}, // E
					106: util.N{Left: 102, Right: 101, Parent: 100}, // C
				},
				Root: 100,
				Size: 1,
			})
			return config{
				name: "Swap/ZRight/YRight",
				z:    root.Right().Right(), // E
				want: want.Right(),
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Rebalance(c.z)
			if diff := cmp.Diff(c.want, got, cmp.Comparer(util.Equal)); diff != "" {
				t.Errorf("Rebalance() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
