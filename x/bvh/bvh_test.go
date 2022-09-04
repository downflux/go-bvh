package bvh

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	nid "github.com/downflux/go-bvh/x/internal/node/id"
)

func TestUpdate(t *testing.T) {
	type config struct {
		name string
		bvh  *BVH
		id   id.ID
		q    hyperrectangle.R
		aabb hyperrectangle.R
		want *BVH
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
			bvh := &BVH{
				lookup: map[id.ID]*node.N{
					1: root,
				},
				root: root,
			}
			return config{
				name: "NoUpdate",
				bvh:  bvh,
				id:   1,
				q:    util.Interval(1, 99),
				aabb: util.Interval(101, 200),
				want: bvh,
			}
		}(),
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
			bvh := &BVH{
				lookup: map[id.ID]*node.N{
					1: root,
				},
				root: root,
			}
			return config{
				name: "DNE",
				bvh:  bvh,
				id:   2,
				q:    util.Interval(1, 99),
				aabb: util.Interval(101, 200),
				want: bvh,
			}
		}(),
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
			want := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: util.Interval(100, 201)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{},
				},
				Root: 100,
			})

			return config{
				name: "Simple",
				bvh: &BVH{
					lookup: map[id.ID]*node.N{
						1: root,
					},
					root: root,
				},
				id:   1,
				q:    util.Interval(101, 200),
				aabb: util.Interval(100, 201),
				want: &BVH{
					lookup: map[id.ID]*node.N{
						1: want,
					},
					root: want,
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			c.bvh.Update(c.id, c.q, c.aabb)
			if diff := cmp.Diff(
				c.want,
				c.bvh,
				cmp.Comparer(util.Equal),
				cmp.AllowUnexported(BVH{}),
			); diff != "" {
				t.Errorf("Update() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestBroadPhase(t *testing.T) {
	type config struct {
		name string
		bvh  *BVH
		q    hyperrectangle.R
		want []id.ID
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
				name: "PartialMatch",
				bvh: &BVH{
					lookup: map[id.ID]*node.N{
						1: root,
					},
					root: root,
				},
				q:    util.Interval(1, 2),
				want: []id.ID{1},
			}
		}(),
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
				name: "NoMatch",
				bvh: &BVH{
					lookup: map[id.ID]*node.N{
						1: root,
					},
					root: root,
				},
				q:    util.Interval(101, 102),
				want: nil,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := BroadPhase(c.bvh, c.q)
			if diff := cmp.Diff(
				c.want,
				got,
				cmpopts.SortSlices(func(a, b id.ID) bool { return a < b }),
			); diff != "" {
				t.Errorf("BroadPhase() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	type config struct {
		name string
		bvh  *BVH
		id   id.ID
		want *BVH
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
			bvh := &BVH{
				lookup: map[id.ID]*node.N{
					1: root,
				},
				root: root,
			}
			return config{
				name: "Trivial",
				bvh:  bvh,
				id:   1,
				want: &BVH{
					lookup: map[id.ID]*node.N{},
					root:   nil,
				},
			}
		}(),
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
			bvh := &BVH{
				lookup: map[id.ID]*node.N{
					1: root,
				},
				root: root,
			}
			return config{
				name: "DNE",
				bvh:  bvh,
				id:   2,
				want: bvh,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			c.bvh.Remove(c.id)
			if diff := cmp.Diff(
				c.want,
				c.bvh,
				cmp.Comparer(util.Equal),
				cmp.AllowUnexported(BVH{}),
			); diff != "" {
				t.Errorf("Remove() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	type config struct {
		name string
		bvh  *BVH
		id   id.ID
		aabb hyperrectangle.R
		want *BVH
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
			want := &BVH{
				lookup: map[id.ID]*node.N{1: root},
				root:   root,
			}
			return config{
				name: "NilRoot",
				bvh:  New(),
				id:   1,
				aabb: util.Interval(0, 100),
				want: want,
			}
		}(),
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
			want := &BVH{
				lookup: map[id.ID]*node.N{1: root},
				root:   root,
			}
			return config{
				name: "DuplicateID",
				bvh:  want,
				id:   1,
				aabb: util.Interval(0, 100),
				want: want,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			c.bvh.Insert(c.id, c.aabb)
			if diff := cmp.Diff(
				c.want,
				c.bvh,
				cmp.Comparer(util.Equal),
				cmp.AllowUnexported(BVH{}),
			); diff != "" {
				t.Errorf("Insert() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
