package rotate

import (
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/internal/node/id"
)

func TestRebalance(t *testing.T) {
	type config struct {
		name string
		z    *node.N
		want *node.N
	}

	configs := []config{
		{
			name: "Root",
			z: util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: util.Interval(0, 100)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{},
				},
				Root: 100,
				Size: 1,
			}),
			want: nil,
		},
		func() config {
			root := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					101: {1: util.Interval(0, 100)},
					102: {2: util.Interval(101, 200)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{Left: 101, Right: 102},
					101: util.N{Parent: 100},
					102: util.N{Parent: 100},
				},
				Root: 100,
				Size: 1,
			})
			return config{
				name: "Leaf",
				z:    root.Left(),
				want: root,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Rebalance(c.z)
			if diff := cmp.Diff(c.want, got, cmp.Comparer(util.Equal)); diff != "" {
				t.Errorf("Rebalance() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
