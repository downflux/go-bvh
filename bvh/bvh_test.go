package bvh

import (
	"log"
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	nid "github.com/downflux/go-bvh/internal/node/id"
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
				Size: 1,
			})
			bvh := &BVH{
				lookup: map[id.ID]*node.N{
					1: root,
				},
				root: root,
				size: 1,
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
				Size: 1,
			})
			bvh := &BVH{
				lookup: map[id.ID]*node.N{
					1: root,
				},
				root: root,
				size: 1,
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
				Size: 1,
			})
			want := util.New(util.T{
				Data: map[nid.ID]map[id.ID]hyperrectangle.R{
					100: {1: util.Interval(100, 201)},
				},
				Nodes: map[nid.ID]util.N{
					100: util.N{},
				},
				Root: 100,
				Size: 1,
			})

			return config{
				name: "Simple",
				bvh: &BVH{
					lookup: map[id.ID]*node.N{
						1: root,
					},
					root: root,
					size: 1,
				},
				id:   1,
				q:    util.Interval(101, 200),
				aabb: util.Interval(100, 201),
				want: &BVH{
					lookup: map[id.ID]*node.N{
						1: want,
					},
					root: want,
					size: 1,
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
				Size: 1,
			})
			return config{
				name: "PartialMatch",
				bvh: &BVH{
					lookup: map[id.ID]*node.N{
						1: root,
					},
					root: root,
					size: 1,
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
				Size: 1,
			})
			return config{
				name: "NoMatch",
				bvh: &BVH{
					lookup: map[id.ID]*node.N{
						1: root,
					},
					root: root,
					size: 1,
				},
				q:    util.Interval(101, 102),
				want: []id.ID{},
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
				Size: 1,
			})
			bvh := &BVH{
				lookup: map[id.ID]*node.N{
					1: root,
				},
				root: root,
				size: 1,
			}
			return config{
				name: "Trivial",
				bvh:  bvh,
				id:   1,
				want: &BVH{
					lookup: map[id.ID]*node.N{},
					root:   nil,
					size:   1,
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
				Size: 1,
			})
			bvh := &BVH{
				lookup: map[id.ID]*node.N{
					1: root,
				},
				root: root,
				size: 1,
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
				Size: 1,
			})
			want := &BVH{
				lookup: map[id.ID]*node.N{1: root},
				root:   root,
				size:   1,
			}
			return config{
				name: "NilRoot",
				bvh:  New(O{Size: 1}),
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
				Size: 1,
			})
			want := &BVH{
				lookup: map[id.ID]*node.N{1: root},
				root:   root,
				size:   1,
			}
			return config{
				name: "DuplicateID",
				bvh:  want,
				id:   1,
				aabb: util.Interval(0, 100),
				want: want,
			}
		}(),
		// Based on experimental results, we want to validate the
		// branching algorithm is acting in an intuitive manner.
		func() config {
			data := map[nid.ID]map[id.ID]hyperrectangle.R{
				101: {1: *hyperrectangle.New([]float64{346, 0}, []float64{347, 1})},
				102: {2: *hyperrectangle.New([]float64{239, 0}, []float64{240, 1})},
				103: {3: *hyperrectangle.New([]float64{896, 0}, []float64{897, 1})},
				104: {4: *hyperrectangle.New([]float64{826, 0}, []float64{827, 1})},
			}

			bvh := New(O{Size: 1, Logger: log.Default()})
			bvh.Insert(1, data[101][1])
			bvh.Insert(2, data[102][2])
			bvh.Insert(3, data[103][3])

			root := util.New(util.T{
				Data: data,
				Nodes: map[nid.ID]util.N{
					100: {Left: 105, Right: 106},
					105: {Left: 104, Right: 103, Parent: 100},
					106: {Left: 102, Right: 101, Parent: 100},
					101: {Parent: 106},
					102: {Parent: 106},
					103: {Parent: 105},
					104: {Parent: 105},
				},
				Root: 100,
				Size: 1,
			})
			want := &BVH{
				lookup: map[id.ID]*node.N{
					1: root.Right().Right(),
					2: root.Right().Left(),
					3: root.Left().Right(),
					4: root.Left().Left(),
				},
				root: root,
				size: 1,
			}

			return config{
				name: "Experimental",
				bvh:  bvh,
				id:   4,
				aabb: data[104][4],
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
				cmpopts.IgnoreFields(BVH{}, "logger"),
			); diff != "" {
				t.Errorf("Insert() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
