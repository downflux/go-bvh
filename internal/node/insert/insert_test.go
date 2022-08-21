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

func Interval(min float64, max float64) hyperrectangle.R {
	return *hyperrectangle.New(*vector.New(min), *vector.New(max))
}

func TestFindCandidate(t *testing.T) {
	type config struct {
		name string

		nodes allocation.C[*node.N]
		rid   allocation.ID
		bound hyperrectangle.R

		want allocation.ID
	}

	configs := []config{
		{
			name: "SimpleRoot",
			nodes: allocation.C[*node.N]{
				1: node.New(node.O{
					ID:    "foo",
					Index: 1,
					Bound: *hyperrectangle.New(
						*vector.New(0, 0),
						*vector.New(10, 10),
					),
				}),
			},
			rid: 1,
			bound: *hyperrectangle.New(
				*vector.New(100, 100),
				*vector.New(1000, 1000),
			),
			want: 1,
		},
		{
			// Check that we do not travel further down the tree
			// than necessary -- we incur a depth penalty via the
			// inherited cost heuristic.
			name: "Root",
			nodes: allocation.C[*node.N]{
				1: node.New(node.O{
					Index: 1,
					Left:  2,
					Right: 3,
					Bound: Interval(0, 100),
				}),

				2: node.New(node.O{
					ID:    "foo",
					Index: 2,
					Bound: Interval(0, 100),
				}),
				3: node.New(node.O{
					ID:    "bar",
					Index: 3,
					Bound: Interval(0, 100),
				}),
			},
			rid:   1,
			bound: Interval(0, 1),
			want:  1,
		},
		{
			name: "Leaf",
			nodes: allocation.C[*node.N]{
				1: node.New(node.O{
					Index: 1,
					Left:  2,
					Right: 3,
					Bound: Interval(0, 100),
				}),

				2: node.New(node.O{
					ID:    "foo",
					Index: 2,
					Bound: Interval(0, 10),
				}),
				3: node.New(node.O{
					ID:    "bar",
					Index: 3,
					Bound: Interval(50, 100),
				}),
			},
			rid:   1,
			bound: Interval(45, 60),
			want:  3,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := findSibling(c.nodes, c.rid, c.bound)
			if diff := cmp.Diff(c.want, got, cmp.AllowUnexported(node.N{}, hyperrectangle.R{})); diff != "" {
				t.Errorf("findCandidate() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestCreateParent(t *testing.T) {
	type result struct {
		allocation allocation.C[*node.N]
		root       allocation.ID
	}

	type config struct {
		name string

		nodes allocation.C[*node.N]
		rid   allocation.ID
		id    point.ID
		bound hyperrectangle.R

		want result
	}

	configs := []config{
		{
			name: "TrivialRoot",
			nodes: allocation.C[*node.N]{
				1: node.New(node.O{
					ID:    "foo",
					Index: 1,
					Bound: Interval(0, 100),
				}),
			},
			rid:   1,
			id:    "bar",
			bound: Interval(99, 101),
			want: result{
				root: 2,
				allocation: allocation.C[*node.N]{
					1: node.New(node.O{
						ID:     "foo",
						Index:  1,
						Parent: 2,
						Bound:  Interval(0, 100),
					}),
					2: node.New(node.O{
						Left:  3,
						Right: 1,
						Bound: Interval(0, 101),
					}),
					3: node.New(node.O{
						ID:     "bar",
						Index:  3,
						Parent: 2,
						Bound:  Interval(99, 101),
					}),
				},
			},
		},
		{
			name: "Ancestor",
			nodes: allocation.C[*node.N]{
				1: node.New(node.O{
					Index: 1,
					Left:  2,
					Right: 3,
					Bound: Interval(0, 100),
				}),

				2: node.New(node.O{
					ID:     "foo",
					Index:  2,
					Parent: 1,
					Bound:  Interval(0, 10),
				}),
				3: node.New(node.O{
					ID:     "bar",
					Index:  3,
					Parent: 1,
					Bound:  Interval(90, 100),
				}),
			},
			rid:   2,
			id:    "baz",
			bound: Interval(99, 101),
			want: result{
				root: 1,
				allocation: allocation.C[*node.N]{
					1: node.New(node.O{
						Index: 1,
						Left:  4,
						Right: 3,
						Bound: Interval(0, 101),
					}),

					2: node.New(node.O{
						ID:     "foo",
						Index:  2,
						Parent: 4,
						Bound:  Interval(0, 10),
					}),
					3: node.New(node.O{
						ID:     "bar",
						Index:  3,
						Parent: 1,
						Bound:  Interval(90, 100),
					}),

					4: node.New(node.O{
						Index:  4,
						Parent: 1,
						Left:   5,
						Right:  2,
						Bound:  Interval(0, 101),
					}),
					5: node.New(node.O{
						ID:     "baz",
						Index:  5,
						Parent: 4,
						Bound:  Interval(99, 101),
					}),
				},
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			pid := createParent(c.nodes, c.rid, c.id, c.bound)
			for ; node.Parent(c.nodes, c.nodes[pid]) != nil; pid = node.Parent(c.nodes, c.nodes[pid]).Index() {
			}
			if diff := cmp.Diff(
				c.want,
				result{
					allocation: c.nodes,
					root:       pid,
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
				for i, w := range c.want.allocation {
					t.Errorf("want: %v: %v\n", i, w)
				}
				for i, g := range c.nodes {
					t.Errorf("got: %v: %v\n", i, g)
				}
				t.Errorf("createParent() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestRotate(t *testing.T) {
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

	configs := []config{}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			rotate(c.nodes, c.aid)
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
				for i, w := range c.want.allocation {
					t.Errorf("want: %v: %v\n", i, w)
				}
				for i, g := range c.nodes {
					t.Errorf("got: %v: %v\n", i, g)
				}
				t.Errorf("createParent() mismatch (-want +got):\n%v", diff)
			}

		})
	}
}

func TestExecute(t *testing.T) {
	type result struct {
		allocation allocation.C[*node.N]
		root       allocation.ID
	}

	type config struct {
		name string

		nodes allocation.C[*node.N]
		rid   allocation.ID
		id    point.ID
		bound hyperrectangle.R
		want  result
	}

	configs := []config{
		{
			name:  "Trival",
			nodes: *allocation.New[*node.N](),
			rid:   404,
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
			c := *allocation.New[*node.N]()
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
				rid:   n.Index(),
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
			qid := Execute(c.nodes, c.rid, c.id, c.bound)
			got := result{
				allocation: c.nodes,
				root:       qid,
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
