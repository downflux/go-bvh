package remove

import (
	"testing"

	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/allocation/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/google/go-cmp/cmp"
)

func interval(min float64, max float64) hyperrectangle.R {
	return *hyperrectangle.New(*vector.New(min), *vector.New(max))
}

func TestExecute(t *testing.T) {
	type result struct {
		allocation allocation.C[*node.N]
		root       id.ID
	}

	type config struct {
		name string

		data allocation.C[*node.N]
		nid  id.ID

		want result
	}

	configs := []config{
		{
			name: "RemoveRoot",
			data: allocation.C[*node.N]{
				1: node.New(node.O{
					ID:    "foo",
					Index: 1,
					Bound: interval(0, 100),
				}),
			},
			nid: 1,
			want: result{
				allocation: allocation.C[*node.N]{},
				root:       0,
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			rid := Execute(c.data, c.nid)
			got := result{
				allocation: c.data,
				root:       rid,
			}

			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(
					node.N{},
					hyperrectangle.R{},
				),
				cmp.Comparer(
					func(r result, s result) bool {
						return util.Equal(r.allocation, r.root, s.allocation, s.root)
					},
				),
			); diff != "" {
				t.Errorf("Execute() mismatch (-want +got):\n%v", diff)
			}

		})
	}
}
