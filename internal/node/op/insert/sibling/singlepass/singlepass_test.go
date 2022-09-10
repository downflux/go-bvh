package singlepass

import (
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/internal/node/id"
)

func TestExecute(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		aabb hyperrectangle.R
		want *node.N
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
				aabb: util.Interval(0, 1),
				want: n,
			}
		}(),
		// Based on experimental results, we want to validate the
		// branching algorithm is acting in an intuitive manner.
		func() config {
			n := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {
						1: *hyperrectangle.New(
							[]float64{346, 0}, []float64{347, 1},
						),
					},
					102: {
						2: *hyperrectangle.New(
							[]float64{239, 0}, []float64{240, 1},
						),
					},
					103: {
						3: *hyperrectangle.New(
							[]float64{896, 0}, []float64{897, 1},
						),
					},
				},
				Nodes: map[nid.ID]util.N{
					100: {Left: 103, Right: 104},
					101: {Parent: 104},
					102: {Parent: 104},
					103: {Parent: 100},
					104: {Left: 102, Right: 101, Parent: 100},
				},
				Root: 100,
				Size: 1,
			})
			return config{
				name: "Experimental",
				n:    n,
				aabb: *hyperrectangle.New(
					[]float64{826, 0},
					[]float64{827, 1},
				),
				want: n.Left(),
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Execute(c.n, c.aabb)
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
