package insert

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
		root *node.N
		size uint
		id   id.ID
		aabb hyperrectangle.R
		want *node.N
	}

	configs := []config{
		{
			name: "Trivial",
			root: nil,
			size: 1,
			id:   1,
			aabb: util.Interval(0, 100),
			want: util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: util.Interval(0, 100)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{},
				},
				Root: 100,
				Size: 1,
			}),
		},
		func() config {
			return config{
				name: "Root/Split",
				root: util.New(util.T{
					Data: map[nid.ID]map[id.ID]hyperrectangle.R{
						100: {1: util.Interval(0, 10)},
					},
					Nodes: map[nid.ID]util.N{
						100: util.N{},
					},
					Root: 100,
					Size: 1,
				}),
				size: 1,
				id:   2,
				aabb: util.Interval(20, 50),
				want: util.New(util.T{
					Data: map[nid.ID]map[id.ID]hyperrectangle.R{
						101: {1: util.Interval(0, 10)},
						102: {2: util.Interval(20, 50)},
					},
					Nodes: map[nid.ID]util.N{
						100: util.N{Left: 102, Right: 101},
						101: util.N{Parent: 100},
						102: util.N{Parent: 100},
					},
					Root: 100,
					Size: 1,
				}),
			}
		}(),
		func() config {
			return config{
				name: "Root/Insert",
				root: util.New(util.T{
					Data: map[nid.ID]map[id.ID]hyperrectangle.R{
						100: {1: util.Interval(0, 10)},
					},
					Nodes: map[nid.ID]util.N{
						100: util.N{},
					},
					Root: 100,
					Size: 2,
				}),
				size: 1,
				id:   2,
				aabb: util.Interval(20, 50),
				want: util.New(util.T{
					Data: map[nid.ID]map[id.ID]hyperrectangle.R{
						100: {
							1: util.Interval(0, 10),
							2: util.Interval(20, 50),
						},
					},
					Nodes: map[nid.ID]util.N{
						100: util.N{},
					},
					Root: 100,
					Size: 2,
				}),
			}
		}(),
		// Based on experimental results, we want to validate the
		// branching algorithm is acting in an intuitive manner.
		func() config {
			data := map[nid.ID]map[id.ID]hyperrectangle.R{
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
				104: {
					4: *hyperrectangle.New(
						[]float64{826, 0}, []float64{827, 1},
					),
				},
			}
			root := util.New(util.T{
				Data: data,
				Nodes: map[nid.ID]util.N{
					100: {Left: 103, Right: 105},
					101: {Parent: 105},
					102: {Parent: 105},
					103: {Parent: 100},
					105: {Left: 102, Right: 101, Parent: 100},
				},
				Root: 100,
				Size: 1,
			})
			want := util.New(util.T{
				Data: data,
				Nodes: map[nid.ID]util.N{
					100: {Left: 105, Right: 106},
					101: {Parent: 106},
					102: {Parent: 106},
					103: {Parent: 105},
					104: {Parent: 105},
					105: {Left: 104, Right: 103, Parent: 100},
					106: {Left: 102, Right: 101, Parent: 100},
				},
				Root: 100,
				Size: 1,
			})
			return config{
				name: "Experimental",
				root: root,
				size: 1,
				id:   4,
				aabb: *hyperrectangle.New(
					[]float64{826, 0},
					[]float64{827, 1},
				),
				want: want,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Execute(c.root, c.size, c.id, c.aabb)
			if diff := cmp.Diff(c.want, got.Root(), cmp.Comparer(util.Equal)); diff != "" {
				t.Errorf("Execute() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
