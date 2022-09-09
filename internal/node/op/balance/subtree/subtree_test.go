package subtree

import (
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/internal/node/id"
)

func TestNew(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		want *T
	}

	configs := []config{
		func() config {
			n := util.New(util.T{
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
				name: "Leaf",
				n:    n,
				want: &T{A: n},
			}
		}(),
		func() config {
			const A = 100
			const B = 101
			const C = 102
			const D = 103
			const E = 104
			const F = 105
			const G = 106
			//      A
			//     / \
			//    /   \
			//   B     C
			//  / \   / \
			// D   E F   G
			n := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					D: {1: util.Interval(0, 100)},
					E: {2: util.Interval(101, 200)},
					F: {3: util.Interval(201, 300)},
					G: {4: util.Interval(301, 400)},
				},
				Nodes: map[nid.ID]util.N{
					A: util.N{Left: B, Right: C},
					B: util.N{Left: D, Right: E, Parent: A},
					C: util.N{Left: F, Right: G, Parent: A},
					D: util.N{Parent: B},
					E: util.N{Parent: B},
					F: util.N{Parent: C},
					G: util.N{Parent: C},
				},
				Root: A,
				Size: 1,
			})
			return config{
				name: "EqualHeight",
				n:    n,
				want: &T{A: n},
			}
		}(),
		func() config {
			const A = 100
			const B = 101
			const C = 102
			const F = 105
			const G = 106
			//   A
			//  / \
			// B   C
			//    / \
			//   F   G
			n := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					B: {1: util.Interval(0, 100)},
					F: {3: util.Interval(201, 300)},
					G: {5: util.Interval(401, 500)},
				},
				Nodes: map[nid.ID]util.N{
					A: util.N{Left: B, Right: C},
					B: util.N{Parent: A},
					C: util.N{Left: F, Right: G, Parent: A},
					F: util.N{Parent: C},
					G: util.N{Parent: C},
				},
				Root: A,
				Size: 1,
			})
			return config{
				name: "Deep/Right",
				n:    n,
				want: &T{
					A: n,
					B: n.Left(),
					C: n.Right(),
					F: n.Right().Right(),
					G: n.Right().Left(),
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := New(c.n)
			if diff := cmp.Diff(
				c.want,
				got,
				cmp.Comparer(util.Equal)); diff != "" {
				t.Errorf("New() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
