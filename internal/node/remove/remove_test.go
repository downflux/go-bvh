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
		{
			name: "RemoveChild",
			data: allocation.C[*node.N]{
				1: node.New(node.O{
					Index: 1,
					Left:  2,
					Right: 3,
					Bound: interval(0, 100),
				}),

				2: node.New(node.O{
					ID:     "foo",
					Index:  2,
					Parent: 1,
					Bound:  interval(0, 1),
				}),
				3: node.New(node.O{
					ID:     "foo",
					Index:  3,
					Parent: 1,
					Bound:  interval(99, 100),
				}),
			},
			nid: 2,
			want: result{
				allocation: allocation.C[*node.N]{
					3: node.New(node.O{
						ID:     "foo",
						Index:  3,
						Parent: 1,
						Bound:  interval(99, 100),
					}),
				},
				root: 3,
			},
		},
		{
			name: "RemoveGrandChildRefit",
			data: allocation.C[*node.N]{
				1: node.New(node.O{
					Index: 1,
					Left:  2,
					Right: 3,
					Bound: interval(0, 100),
				}),

				2: node.New(node.O{
					Index:  2,
					Parent: 1,
					Left:   4,
					Right:  5,
					Bound:  interval(0, 10),
				}),
				3: node.New(node.O{
					Index:  3,
					Parent: 1,
					Left:   6,
					Right:  7,
					Bound:  interval(90, 100),
				}),

				4: node.New(node.O{
					ID:     "foo",
					Index:  4,
					Parent: 2,
					Bound:  interval(0, 1),
				}),
				5: node.New(node.O{
					ID:     "bar",
					Index:  5,
					Parent: 2,
					Bound:  interval(9, 10),
				}),

				6: node.New(node.O{
					ID:     "baz",
					Index:  6,
					Parent: 3,
					Bound:  interval(90, 91),
				}),
				7: node.New(node.O{
					ID:     "qux",
					Index:  7,
					Parent: 3,
					Bound:  interval(99, 100),
				}),
			},
			nid: 7,
			want: result{
				allocation: allocation.C[*node.N]{
					1: node.New(node.O{
						Index: 1,
						Left:  2,
						Right: 6,
						Bound: interval(0, 91),
					}),

					2: node.New(node.O{
						Index:  2,
						Parent: 1,
						Left:   4,
						Right:  5,
						Bound:  interval(0, 10),
					}),

					4: node.New(node.O{
						ID:     "foo",
						Index:  4,
						Parent: 2,
						Bound:  interval(0, 1),
					}),
					5: node.New(node.O{
						ID:     "bar",
						Index:  5,
						Parent: 2,
						Bound:  interval(9, 10),
					}),

					6: node.New(node.O{
						ID:     "baz",
						Index:  6,
						Parent: 1,
						Bound:  interval(90, 91),
					}),
				},
				root: 1,
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
