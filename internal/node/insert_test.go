package node

import (
	"testing"

	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

func TestInsert(t *testing.T) {
	type config struct {
		name string

		root  *N
		id    point.ID
		bound hyperrectangle.R
		want  *N
	}

	configs := []config{
		{
			name: "Trival",
			root: nil,
			id: "foo",
			bound: *hyperrectangle.New(
				*vector.New(0, 0),
				*vector.New(10, 10),
			),
			want: &N{
				id: "foo",

				nodes: nil,
				index: 0,

				bound: *hyperrectangle.New(
					*vector.New(0, 0),
					*vector.New(10, 10),
				),
			},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := Insert(c.root, c.id, c.bound); got != c.want {
				t.Errorf("Insert() = %v, want = %v", got, c.want)
			}
		})
	}
}
