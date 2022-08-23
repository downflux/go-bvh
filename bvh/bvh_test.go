package bvh

import (
	"testing"

	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/google/go-cmp/cmp"
)

func interval(min float64, max float64) hyperrectangle.R {
	return *hyperrectangle.New(*vector.New(min), *vector.New(max))
}

type n struct {
	id    point.ID
	bound hyperrectangle.R
}

func (n n) Bound() hyperrectangle.R { return n.bound }
func (n n) ID() point.ID            { return n.id }

func Equal(a, b *BVH[*n]) bool {
	if !util.Equal(a.allocation, a.root, b.allocation, b.root) {
		return false
	}

	for pid, t := range a.data {
		if !cmp.Equal(
			t,
			b.data[pid],
			cmp.AllowUnexported(
				n{},
				hyperrectangle.R{},
			),
		) {
			return false
		}
	}

	for pid, aid := range a.lookup {
		if !util.Equal(a.allocation, aid, b.allocation, b.lookup[pid]) {
			return false
		}
	}
	return true
}

func TestRemove(t *testing.T) {
	type config struct {
		name string
		bvh  *BVH[*n]
		i    point.ID
		want *BVH[*n]
	}

	configs := []config{
		{
			name: "RemoveRoot",
			bvh: New([]*n{
				{id: "foo", bound: interval(0, 100)},
			}),
			i:    "foo",
			want: New([]*n{}),
		},
		{
			name: "RemoveNode",
			bvh: New([]*n{
				{id: "foo", bound: interval(99, 100)},
				{id: "bar", bound: interval(0, 1)},
			}),
			i: "foo",
			want: New([]*n{
				{id: "bar", bound: interval(0, 1)},
			}),
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			c.bvh.Remove(c.i)
			if diff := cmp.Diff(
				c.want,
				c.bvh,
				cmp.Comparer(func(a, b *BVH[*n]) bool {
					return Equal(a, b)
				})); diff != "" {
				t.Errorf("Remove() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
