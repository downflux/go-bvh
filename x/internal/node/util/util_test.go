package util

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/x/internal/node/id"
)

func TestNew(t *testing.T) {
	type config struct {
		name string
		t    T
		want *node.N
	}

	configs := []config{
		{
			name: "Leaf",
			t: T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: Interval(0, 100)},
				},
				Nodes: map[nid.ID]N{
					100: {},
				},
				Root: 100,
			},
			want: node.New(node.O{
				Nodes: node.Cache(),
				ID:    100,
				Data: map[id.ID]hyperrectangle.R{
					1: Interval(0, 100),
				},
			}),
		}, func() config {
			c := node.Cache()
			r := node.New(node.O{
				Nodes: c,
				ID:    100,
				Left:  101,
				Right: 102,
			})
			node.New(node.O{
				Nodes:  c,
				ID:     101,
				Parent: 100,
				Data: map[id.ID]hyperrectangle.R{
					1: Interval(0, 100),
				},
			})
			node.New(node.O{
				Nodes:  c,
				ID:     102,
				Parent: 100,
				Data: map[id.ID]hyperrectangle.R{
					2: Interval(101, 200),
				},
			})

			return config{
				name: "Root",
				t: T{
					Data: map[nid.ID]map[id.ID]hyperrectangle.R{
						101: {1: Interval(0, 100)},
						102: {2: Interval(101, 200)},
					},
					Nodes: map[nid.ID]N{
						100: {Left: 101, Right: 102},
						101: {Parent: 100},
						102: {Parent: 100},
					},
					Root: 100,
				},
				want: r,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := New(c.t)
			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(
					node.N{},
					node.C{},
					hyperrectangle.R{},
				),
			); diff != "" {
				t.Errorf("New() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
