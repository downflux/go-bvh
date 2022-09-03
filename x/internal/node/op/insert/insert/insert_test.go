package insert

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
		m    *node.N
		want *node.N
	}

	configs := []config{
		func() config {
			c := node.Cache()
			n := node.New(node.O{
				Nodes: c,
				Data:  map[id.ID]hyperrectangle.R{1: util.Interval(0, 100)},
			})
			m := node.New(node.O{
				Nodes: c,
				Data:  map[id.ID]hyperrectangle.R{2: util.Interval(101, 200)},
			})
			return config{
				name: "Root",
				n:    n,
				m:    m,
				want: util.New(util.T{
					Data: map[nid.ID]map[id.ID]hyperrectangle.R{
						101: {1: util.Interval(0, 100)},
						102: {2: util.Interval(101, 200)},
					},
					Nodes: map[nid.ID]util.N{
						100: util.N{Left: 102, Right: 101},
						101: util.N{Parent: 100},
						102: util.N{Parent: 100},
					},
					Root: 100,
				}),
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Execute(c.n, c.m)
			if diff := cmp.Diff(c.want, got, cmp.Comparer(util.Equal)); diff != "" {
				t.Errorf("Execute() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
