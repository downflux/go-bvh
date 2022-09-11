package rotate

import (
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/rotate/rotation/aabb"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/internal/node/id"
)

func TestExecute(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		f    R
		want *node.N // root
	}

	configs := []config{
		func() config {
			data := map[nid.ID]map[id.ID]hyperrectangle.R{
				100: {1: util.Interval(0, 100)},
			}
			return config{
				name: "NoRotate/Root",
				n: util.New(util.T{
					Data: data,
					Nodes: map[nid.ID]util.N{
						100: util.N{},
					},
					Root: 100,
					Size: 1,
				}),
				want: util.New(util.T{
					Data: data,
					Nodes: map[nid.ID]util.N{
						100: util.N{},
					},
					Root: 100,
					Size: 1,
				}),
			}
		}(),
		func() config {
			data := map[nid.ID]map[id.ID]hyperrectangle.R{
				101: {1: util.Interval(0, 100)},
				102: {2: util.Interval(101, 200)},
			}
			return config{
				name: "NoRotate/NoGrandchildren",
				n: util.New(util.T{
					Data: data,
					Nodes: map[nid.ID]util.N{
						100: util.N{Left: 101, Right: 102},
						101: util.N{Parent: 100},
						102: util.N{Parent: 100},
					},
					Root: 100,
					Size: 1,
				}),
				want: util.New(util.T{
					Data: data,
					Nodes: map[nid.ID]util.N{
						100: util.N{Left: 101, Right: 102},
						101: util.N{Parent: 100},
						102: util.N{Parent: 100},
					},
					Root: 100,
					Size: 1,
				}),
				f: aabb.Generate,
			}
		}(),
		func() config {
			data := map[nid.ID]map[id.ID]hyperrectangle.R{
				101: {1: *hyperrectangle.New([]float64{1, 1}, []float64{2, 2})},       // B
				103: {1: *hyperrectangle.New([]float64{99, 99}, []float64{100, 100})}, // F
				104: {1: *hyperrectangle.New([]float64{0, 0}, []float64{1, 1})},       // G
			}
			return config{
				name: "Rotate/BF",
				//   A
				//  / \
				// B   C
				//    / \
				//   F   G
				n: util.New(util.T{
					Data: data,
					Nodes: map[nid.ID]util.N{
						100: util.N{Left: 101, Right: 102},
						102: util.N{Left: 103, Right: 104, Parent: 100},
						101: util.N{Parent: 100},
						103: util.N{Parent: 102},
						104: util.N{Parent: 102},
					},
					Root: 100,
					Size: 1,
				}),
				//   A
				//  / \
				// F   C
				//    / \
				//   B   G
				want: util.New(util.T{
					Data: data,
					Nodes: map[nid.ID]util.N{
						100: util.N{Left: 103, Right: 102},
						102: util.N{Left: 101, Right: 104},
						103: util.N{Parent: 100},
						101: util.N{Parent: 102},
						104: util.N{Parent: 102},
					},
					Root: 100,
					Size: 1,
				}),
				f: aabb.Generate,
			}
		}(),
		func() config {
			data := map[nid.ID]map[id.ID]hyperrectangle.R{
				101: {1: util.Interval(99, 100)}, // B
				103: {2: util.Interval(1, 2)},    // F
				104: {3: util.Interval(0, 1)},    // G
			}

			//   A
			//  / \
			// B   C
			//    / \
			//   F   G
			o := util.T{
				Data: data,
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 101, Right: 102},
					101: util.N{Parent: 100},
					102: util.N{Left: 103, Right: 104},
					103: util.N{Parent: 102},
					104: util.N{Parent: 102},
				},
				Root: 100,
				Size: 1,
			}
			return config{
				name: "NoRotate",
				n:    util.New(o),
				want: util.New(o),
				f:    aabb.Generate,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Execute(c.n, c.f)
			if diff := cmp.Diff(
				c.want,
				got,
				cmp.Comparer(util.Equal),
			); diff != "" {
				t.Errorf("Execute() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
