package remove

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
		want *node.N
	}

	configs := []config{
		{
			name: "Root",
			n: util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: util.Interval(0, 100)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{},
				},
				Root: 100,
			}),
			want: nil,
		},
		{
			name: "Sibling",
			//   A
			//  / \
			// B   C
			n: util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},   // B
					102: {2: util.Interval(101, 200)}, // C
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 101, Right: 102}, // A
					101: util.N{Parent: 100},           // B
					102: util.N{Parent: 100},           // C
				},
				Root: 100,
			}).Right(), // C
			want: util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
				},
				Nodes: map[nid.ID]util.N{
					101: util.N{},
				},
				Root: 101,
			}), // B
		},
		{
			name: "Ancestor",
			//   A
			//  / \
			// B   C
			//    / \
			//   F   G
			n: util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},   // B
					103: {2: util.Interval(101, 200)}, // F
					104: {3: util.Interval(201, 300)}, // G
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 101, Right: 102},              // A
					101: util.N{Parent: 100},                        // B
					102: util.N{Left: 103, Right: 104, Parent: 100}, // C
					103: util.N{Parent: 102},                        // F
					104: util.N{Parent: 102},                        // G
				},
				Root: 100,
			}).Right().Right(), // G
			//   A
			//  / \
			// B   F
			want: util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},   // B
					103: {2: util.Interval(101, 200)}, // F
				},
				Nodes: map[nid.ID]util.N{
					100: {Left: 101, Right: 103},
					101: {Parent: 100},
					103: {Parent: 100},
				},
				Root: 100,
			}).Right(), // F
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Execute(c.n)
			if diff := cmp.Diff(c.want, got, cmp.Comparer(util.Equal)); diff != "" {
				t.Errorf("Execute() mismatch (-want +got):\n%v", diff)
			}
			if c.want != nil {
				// Ensure we are actually deleting data from the
				// lookup table to ensure no memory leaks.
				if diff := cmp.Diff(
					c.want.Cache(),
					got.Cache(),
					cmp.AllowUnexported(
						node.C{},
						node.N{},
						hyperrectangle.R{},
					),
				); diff != "" {
					t.Errorf("Cache() mismatch (-want +got):\n%v", diff)
				}
			}
		})
	}
}
