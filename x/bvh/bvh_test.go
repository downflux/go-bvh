package bvh

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	nid "github.com/downflux/go-bvh/x/internal/node/id"
)

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
