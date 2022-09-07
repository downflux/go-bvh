package split

import (
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/internal/node/id"
)

func TestExecute(t *testing.T) {
	type result struct {
		n *node.N
		m *node.N
	}

	type config struct {
		name string
		p    P
		n    *node.N
		want result
	}

	configs := []config{
		{
			name: "Split",
			n: util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {
						1: util.Interval(0, 100),
						2: util.Interval(51, 100),
					},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{},
				},
				Root: 100,
				Size: 10,
			}),
			p: func(data map[id.ID]hyperrectangle.R) (map[id.ID]hyperrectangle.R, map[id.ID]hyperrectangle.R) {
				return partition(data, vector.AXIS_X, 50)
			},
			want: result{
				n: util.New(util.T{
					Data: map[nid.ID]map[id.ID]hyperrectangle.R{
						101: {1: util.Interval(0, 100)},
					},
					Nodes: map[nid.ID]util.N{
						101: util.N{},
					},
					Root: 101,
					Size: 10,
				}),
				m: util.New(util.T{
					Data: map[nid.ID]map[id.ID]hyperrectangle.R{
						101: {2: util.Interval(51, 100)},
					},
					Nodes: map[nid.ID]util.N{
						101: util.N{},
					},
					Root: 101,
					Size: 10,
				}),
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Execute(c.n, c.p)
			if diff := cmp.Diff(
				c.want, result{
					n: c.n,
					m: got,
				},
				cmp.AllowUnexported(result{}),
				cmp.Comparer(util.Equal),
			); diff != "" {
				t.Errorf("Execute() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
