package node_test

import (
	"testing"

	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/google/go-cmp/cmp"
)

type result struct {
	allocation allocation.C[*node.N]
	root       allocation.ID
}

func TestInsert(t *testing.T) {
	type config struct {
		name string

		nodes allocation.C[*node.N]
		root  *node.N
		id    point.ID
		bound hyperrectangle.R
		want  result
	}

	configs := []config{
		{
			name:  "Trival",
			nodes: allocation.New[*node.N](nil),
			root:  nil,
			id:    "foo",
			bound: *hyperrectangle.New(
				*vector.New(0, 0),
				*vector.New(10, 10),
			),
			want: result{
				allocation: allocation.C[*node.N]{
					1: node.New(node.O{
						ID:    "foo",
						Index: 1,
						Bound: *hyperrectangle.New(
							*vector.New(0, 0),
							*vector.New(10, 10),
						),
					}),
				},
				root: 1,
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			q := node.Insert(c.nodes, c.root, c.id, c.bound)
			got := result{
				allocation: c.nodes,
				root:       q.Index(),
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
				t.Errorf("Insert() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
