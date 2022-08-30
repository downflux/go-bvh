package insert

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
)

func TestParent(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		pid  id.ID
		aabb hyperrectangle.R
		want *node.N
	}

	configs := []config{
		{
			name: "NewRoot",
			n: util.New(util.T{
				Data: map[util.NodeID]map[id.ID]hyperrectangle.R{
					100: {1: util.Interval(0, 100)},
				},
				Root: 100,
			}),
			pid:  2,
			aabb: util.Interval(101, 200),
			want: util.New(util.T{
				Data: map[util.NodeID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
					102: {2: util.Interval(101, 200)},
				},
				Nodes: map[util.NodeID]util.N{
					100: {Left: 102, Right: 101},
				},
				Root: 100,
			}),
		},
		{
			name: "SwapRoot",
			n: util.New(util.T{
				Data: map[util.NodeID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
					102: {2: util.Interval(101, 200)},
				},
				Nodes: map[util.NodeID]util.N{
					100: {Left: 101, Right: 102},
				},
				Root: 100,
			}),
			pid:  3,
			aabb: util.Interval(201, 300),
			want: util.New(util.T{
				Data: map[util.NodeID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
					102: {2: util.Interval(101, 200)},
					103: {3: util.Interval(201, 300)},
				},
				Nodes: map[util.NodeID]util.N{
					100: {Left: 103, Right: 104},
					104: {Left: 102, Right: 101},
				},
				Root: 100,
			}),
		},
		func() config {
			data := map[util.NodeID]map[id.ID]hyperrectangle.R{
				101: {1: util.Interval(0, 100)},   // w
				102: {2: util.Interval(101, 200)}, // x
				103: {3: util.Interval(201, 300)}, // y
				104: {4: util.Interval(301, 400)}, // z
			}
			return config{
				name: "SwapInternal",
				//   A
				//  / \
				// w   C
				//    / \
				//   x   y
				n: util.New(util.T{
					Data: data,
					Nodes: map[util.NodeID]util.N{
						100: {Left: 101, Right: 105}, // A
						105: {Left: 102, Right: 103}, // C
					},
					Root: 100,
				}).Right(), // C
				pid:  4, // z
				aabb: data[104][4],
				//   A
				//  / \
				// w   D
				//    / \
				//   z   C
				//      / \
				//     x   y
				want: util.New(util.T{
					Data: data,
					Nodes: map[util.NodeID]util.N{
						100: {Left: 101, Right: 105}, // A
						105: {Left: 104, Right: 106}, // D
						106: {Left: 102, Right: 103}, // C
					},
					Root: 100,
				}).Right(), // D
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := parent(c.n, c.pid, c.aabb)
			if diff := cmp.Diff(
				c.want.Root(),
				got.Root(),
				cmp.Comparer(util.Equal),
			); diff != "" {
				t.Errorf("parent() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
