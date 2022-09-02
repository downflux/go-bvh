package rotation

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/x/internal/node/id"
)

func TestGenerate(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		want []R
	}

	configs := []config{
		func() config {
			data := map[nid.ID]map[id.ID]hyperrectangle.R{
				100: {1: util.Interval(0, 100)},
			}
			return config{
				name: "Leaf",
				n: util.New(util.T{
					Data: data,
					Nodes: map[nid.ID]util.N{
						100: util.N{},
					},
					Root: 100,
				}),
				want: nil,
			}
		}(),
		func() config {
			data := map[nid.ID]map[id.ID]hyperrectangle.R{
				102: {1: util.Interval(0, 100)},
				103: {2: util.Interval(101, 200)},
				104: {3: util.Interval(201, 300)},
			}
			//     A
			//    / \
			//   B   C
			//  / \
			// D   E
			root := util.New(util.T{
				Data: data,
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 101, Right: 102},
					101: util.N{Left: 103, Right: 104, Parent: 100},
					102: util.N{Parent: 100},
					103: util.N{Parent: 101},
					104: util.N{Parent: 104},
				},
				Root: 100,
			})
			return config{
				name: "CDE",
				n:    root,
				want: []R{
					R{B: root.Right(), C: root.Left(), F: root.Left().Left(), G: root.Left().Right()},
					R{B: root.Right(), C: root.Left(), F: root.Left().Right(), G: root.Left().Left()},
				},
			}
		}(),
		func() config {
			data := map[nid.ID]map[id.ID]hyperrectangle.R{
				101: {1: util.Interval(0, 100)},
				104: {2: util.Interval(101, 200)},
				105: {3: util.Interval(201, 300)},
			}
			//   A
			//  / \
			// B   C
			//    / \
			//   F   G
			root := util.New(util.T{
				Data: data,
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 101, Right: 102},
					101: util.N{Parent: 100},
					102: util.N{Left: 104, Right: 105, Parent: 100},
					104: util.N{Parent: 102},
					105: util.N{Parent: 105},
				},
				Root: 100,
			})
			return config{
				name: "BFG",
				n:    root,
				want: []R{
					R{B: root.Left(), C: root.Right(), F: root.Right().Left(), G: root.Right().Right()},
					R{B: root.Left(), C: root.Right(), F: root.Right().Right(), G: root.Right().Left()},
				},
			}
		}(),
	}

	for _, c := range configs {
		got := Generate(c.n)
		t.Run(c.name, func(t *testing.T) {
			if diff := cmp.Diff(
				c.want,
				got,
				cmp.Comparer(util.Equal),
			); diff != "" {
				t.Errorf("Generate() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
