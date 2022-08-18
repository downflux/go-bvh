package insert

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
			nodes: allocation.New[*node.N](),
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
		func() config {
			c := allocation.New[*node.N]()
			nid := c.Allocate()
			n := node.New(node.O{
				ID:    "foo",
				Index: nid,
				Bound: *hyperrectangle.New(
					*vector.New(0, 0),
					*vector.New(10, 10),
				),
			})
			c.Insert(nid, n)

			return config{
				name:  "SingleNode",
				nodes: c,
				root:  n,
				id:    "bar",
				bound: *hyperrectangle.New(
					*vector.New(20, 20),
					*vector.New(50, 50),
				),
				want: result{
					allocation: allocation.C[*node.N]{
						1: node.New(node.O{
							Index: 1,
							Left:  2,
							Right: 3,
							Bound: *hyperrectangle.New(
								*vector.New(0, 0),
								*vector.New(50, 50),
							),
						}),
						2: node.New(node.O{
							ID:     "bar",
							Index:  2,
							Parent: 1,
							Bound: *hyperrectangle.New(
								*vector.New(20, 20),
								*vector.New(50, 50),
							),
						}),
						3: node.New(node.O{
							ID:     "foo",
							Index:  3,
							Parent: 1,
							Bound: *hyperrectangle.New(
								*vector.New(0, 0),
								*vector.New(10, 10),
							),
						}),
					},
					root: 1,
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			q := Insert(c.nodes, c.root, c.id, c.bound)
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
