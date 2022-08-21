package rotate

import (
	"testing"

	"github.com/downflux/go-bvh/internal/allocation"
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
		root       allocation.ID
	}

	type config struct {
		name string

		nodes allocation.C[*node.N]
		aid   allocation.ID

		want result
	}

	configs := []config{
		{
			name: "ALeaf",
			aid:  1,
			nodes: allocation.C[*node.N]{
				1: node.New(node.O{
					ID:    "foo",
					Index: 1,
					Bound: interval(0, 100),
				}),
			},
			want: result{
				root: 1,
				allocation: allocation.C[*node.N]{
					1: node.New(node.O{
						ID:    "foo",
						Index: 1,
						Bound: interval(0, 100),
					}),
				},
			},
		},
		{
			name: "Root/Terminate",
			aid:  1,
			nodes: allocation.C[*node.N]{
				1: node.New(node.O{
					Index: 1,
					Left:  2,
					Right: 3,
					Bound: interval(0, 100),
				}),
				2: node.New(node.O{
					Index: 2,
					Bound: interval(0, 100),
				}),
				3: node.New(node.O{
					Index: 3,
					Bound: interval(0, 100),
				}),
			},
			want: result{
				root: 1,
				allocation: allocation.C[*node.N]{
					1: node.New(node.O{
						Index: 1,
						Left:  2,
						Right: 3,
						Bound: interval(0, 100),
					}),
					2: node.New(node.O{
						Index: 2,
						Bound: interval(0, 100),
					}),
					3: node.New(node.O{
						Index: 3,
						Bound: interval(0, 100),
					}),
				},
			},
		},
		{
			name: "Rotate",
			aid:  1,
			nodes: allocation.C[*node.N]{
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
					Bound:  interval(1, 2),
				}),
				3: node.New(node.O{
					Index:  3,
					Parent: 1,
					Left:   4,
					Right:  5,
					Bound:  interval(0, 100),
				}),

				4: node.New(node.O{
					ID:     "bar",
					Index:  4,
					Parent: 3,
					Bound:  interval(99, 100),
				}),
				5: node.New(node.O{
					ID:     "baz",
					Index:  5,
					Parent: 3,
					Bound:  interval(0, 1),
				}),
			},
			want: result{
				root: 1,
				allocation: allocation.C[*node.N]{
					1: node.New(node.O{
						Index: 1,
						Left:  4,
						Right: 3,
						Bound: interval(0, 100),
					}),

					2: node.New(node.O{
						ID:     "foo",
						Index:  2,
						Parent: 3,
						Bound:  interval(1, 2),
					}),
					3: node.New(node.O{
						Index:  3,
						Parent: 1,
						Left:   2,
						Right:  5,
						Bound:  interval(0, 2),
					}),

					4: node.New(node.O{
						ID:     "bar",
						Index:  4,
						Parent: 1,
						Bound:  interval(99, 100),
					}),
					5: node.New(node.O{
						ID:     "baz",
						Index:  5,
						Parent: 3,
						Bound:  interval(0, 1),
					}),
				},
			},
		},
		{
			name: "NoRotate",
			aid:  1,
			nodes: allocation.C[*node.N]{
				1: node.New(node.O{
					Index: 1,
					Left:  4,
					Right: 3,
					Bound: interval(0, 100),
				}),

				2: node.New(node.O{
					ID:     "foo",
					Index:  2,
					Parent: 3,
					Bound:  interval(1, 2),
				}),
				3: node.New(node.O{
					Index:  3,
					Parent: 1,
					Left:   2,
					Right:  5,
					Bound:  interval(0, 2),
				}),

				4: node.New(node.O{
					ID:     "bar",
					Index:  4,
					Parent: 1,
					Bound:  interval(99, 100),
				}),
				5: node.New(node.O{
					ID:     "baz",
					Index:  5,
					Parent: 3,
					Bound:  interval(0, 1),
				}),
			},
			want: result{
				root: 1,
				allocation: allocation.C[*node.N]{
					1: node.New(node.O{
						Index: 1,
						Left:  4,
						Right: 3,
						Bound: interval(0, 100),
					}),

					2: node.New(node.O{
						ID:     "foo",
						Index:  2,
						Parent: 3,
						Bound:  interval(1, 2),
					}),
					3: node.New(node.O{
						Index:  3,
						Parent: 1,
						Left:   2,
						Right:  5,
						Bound:  interval(0, 2),
					}),

					4: node.New(node.O{
						ID:     "bar",
						Index:  4,
						Parent: 1,
						Bound:  interval(99, 100),
					}),
					5: node.New(node.O{
						ID:     "baz",
						Index:  5,
						Parent: 3,
						Bound:  interval(0, 1),
					}),
				},
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			Execute(c.nodes, c.aid)
			var rid allocation.ID
			for rid = c.aid; node.Parent(c.nodes, c.nodes[rid]) != nil; rid = node.Parent(c.nodes, c.nodes[rid]).Index() {
			}

			if diff := cmp.Diff(
				c.want,
				result{
					allocation: c.nodes,
					root:       rid,
				},
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
				t.Errorf("createParent() mismatch (-want +got):\n%v", diff)
			}

		})
	}
}
