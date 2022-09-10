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
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Execute(c.root, c.size, c.id, c.aabb).Root()
			if diff := cmp.Diff(c.want, got, cmp.Comparer(util.Equal)); diff != "" {
				t.Errorf("Execute() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
