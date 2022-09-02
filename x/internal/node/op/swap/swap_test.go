package swap

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	nid "github.com/downflux/go-bvh/x/internal/node/id"
)

func TestIsAncestor(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		m    *node.N
		want bool
	}

	configs := []config{
		func() config {
			root := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: util.Interval(0, 100)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{},
				},
				Root: 100,
			})
			return config{
				name: "TrivialRoot",
				n:    root,
				m:    root,
				want: true,
			}
		}(),
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
			})
			return config{
				name: "TrivialSibling",
				n:    root.Left(),
				m:    root.Right(),
				want: false,
			}
		}(),
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
			})
			return config{
				name: "ImmediateChild",
				n:    root.Left(),
				m:    root,
				want: false,
			}
		}(),
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
			})
			return config{
				name: "ImmediateParent",
				n:    root,
				m:    root.Left(),
				want: true,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := isAncestor(c.n, c.m); got != c.want {
				t.Errorf("isAncestor() = %v, want = %v", got, c.want)
			}
		})
	}
}
