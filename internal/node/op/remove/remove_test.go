package remove

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
		id   id.ID
		want *node.N
	}

	configs := []config{
		{
			name: "MultiData",
			n: util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {
						1: util.Interval(0, 100),
						2: util.Interval(101, 200),
					},
				},
				Nodes: map[nid.ID]util.N{
					100: {},
				},
				Root: 100,
			}),
			id: 2,
			want: util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {
						1: util.Interval(0, 100),
					},
				},
				Nodes: map[nid.ID]util.N{
					100: {},
				},
				Root: 100,
			}),
		},
		{
			name: "Leaf",
			n: util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: util.Interval(0, 100)},
				},
				Nodes: map[nid.ID]util.N{
					100: {},
				},
				Root: 100,
			}),
			id:   1,
			want: nil,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Execute(c.n, c.id)
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
