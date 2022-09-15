package bvh

import (
	"fmt"
	"sync"
	"testing"

	"github.com/downflux/go-bvh/container/bruteforce"
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-bvh/perf"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	bhru "github.com/downflux/go-bvh/hyperrectangle/util"
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
				cmp.AllowUnexported(BVH{}, sync.RWMutex{}, sync.Mutex{}),
			); diff != "" {
				t.Errorf("Update() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestBroadPhaseConformance(t *testing.T) {
	type config struct {
		name string
		n    int
		k    vector.D
		size uint
	}

	var configs []config
	for _, k := range perf.PerfTestSize(perf.SizeUnit).K() {
		for _, n := range perf.PerfTestSize(perf.SizeUnit).N() {
			for _, size := range perf.PerfTestSize(perf.SizeUnit).LeafSize() {
				configs = append(configs, config{
					name: fmt.Sprintf("K=%v/N=%v/LeafSize=%v", k, n, size),
					k:    k,
					n:    n,
					size: size,
				})
			}
		}
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			data := map[id.ID]hyperrectangle.R{}
			for i := 0; i < c.n; i++ {
				data[id.ID(i+1)] = bhru.RR(0, 500, c.k)
			}
			bvh := New(O{Size: c.size})
			bf := bruteforce.New()

			for x, aabb := range data {
				bvh.Insert(x, aabb)
				bf.Insert(x, aabb)
			}

			q := bhru.RR(0, 500, c.k)
			got := bvh.BroadPhase(q)
			want := bf.BroadPhase(q)
			if diff := cmp.Diff(
				want, got,
				cmpopts.SortSlices(func(a, b id.ID) bool { return a < b }),
			); diff != "" {
				t.Errorf("BroadPhase() mismatch (-want +got):\n%v", diff)
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
			got := c.bvh.BroadPhase(c.q)
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
				cmp.AllowUnexported(BVH{}, sync.RWMutex{}, sync.Mutex{}),
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
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			c.bvh.Insert(c.id, c.aabb)
			if diff := cmp.Diff(
				c.want,
				c.bvh,
				cmp.Comparer(util.Equal),
				cmp.AllowUnexported(BVH{}, sync.RWMutex{}, sync.Mutex{}),
				cmpopts.IgnoreFields(BVH{}, "logger"),
			); diff != "" {
				t.Errorf("Insert() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
