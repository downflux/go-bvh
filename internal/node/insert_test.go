package node

import (
	"testing"

	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/google/go-cmp/cmp"
)

func TestInsert(t *testing.T) {
	type config struct {
		name string

		nodes allocation.C[*N]
		root  *N
		id    point.ID
		bound hyperrectangle.R
		want  allocation.C[*N]
	}

	configs := []config{
		{
			name:  "Trival",
			nodes: allocation.New[*N](nil),
			root:  nil,
			id:    "foo",
			bound: *hyperrectangle.New(
				*vector.New(0, 0),
				*vector.New(10, 10),
			),
			want: allocation.C[*N]{
				5577006791947779410: New(O{
					ID:    "foo",
					Index: 5577006791947779410,
					Bound: *hyperrectangle.New(
						*vector.New(0, 0),
						*vector.New(10, 10),
					),
				}),
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			Insert(c.nodes, c.root, c.id, c.bound)
			if diff := cmp.Diff(
				c.want,
				c.nodes,
				cmp.AllowUnexported(
					N{},
					hyperrectangle.R{},
				),
			); diff != "" {
				t.Errorf("Insert() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
