package aabb

import (
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/rotate/rotation"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/internal/node/id"
)

func TestList(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		want []rotation.R
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
					Size: 1,
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
				Size: 1,
			})
			return config{
				name: "CDE",
				n:    root,
				want: []rotation.R{
					rotation.R{B: root.Right(), C: root.Left(), F: root.Left().Left(), G: root.Left().Right()},
					rotation.R{B: root.Right(), C: root.Left(), F: root.Left().Right(), G: root.Left().Left()},
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
				Size: 1,
			})
			return config{
				name: "BFG",
				n:    root,
				want: []rotation.R{
					rotation.R{B: root.Left(), C: root.Right(), F: root.Right().Left(), G: root.Right().Right()},
					rotation.R{B: root.Left(), C: root.Right(), F: root.Right().Right(), G: root.Right().Left()},
				},
			}
		}(),
	}

	for _, c := range configs {
		got := list(c.n)
		t.Run(c.name, func(t *testing.T) {
			if diff := cmp.Diff(
				c.want,
				got,
				cmp.Comparer(util.Equal),
			); diff != "" {
				t.Errorf("list() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		want rotation.R
	}

	configs := []config{
		{
			name: "Trivial",
			n: util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: util.Interval(0, 100)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{},
				},
				Root: 100,
				Size: 1,
			}),
			want: rotation.R{},
		},
		func() config {
			n := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {
						1: *hyperrectangle.New(
							[]float64{49, 49},
							[]float64{50, 50},
						),
					},
					103: {
						2: *hyperrectangle.New(
							[]float64{0, 0},
							[]float64{1, 1},
						),
					},
					104: {
						3: *hyperrectangle.New(
							[]float64{50, 50},
							[]float64{100, 100},
						),
					},
				},
				Nodes: map[nid.ID]util.N{
					100: {Left: 101, Right: 102},
					101: {Parent: 100},
					102: {Left: 103, Right: 104, Parent: 100},
					103: {Parent: 102},
					104: {Parent: 102},
				},
				Root: 100,
				Size: 1,
			})
			return config{
				name: "Rotate",
				n:    n,
				want: rotation.R{
					B: n.Left(),
					C: n.Right(),
					F: n.Right().Left(),
					G: n.Right().Right(),
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Generate(c.n)
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
