package node

import (
	"testing"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
)

func interval(min, max float64) hyperrectangle.R {
	return *hyperrectangle.New([]float64{min}, []float64{max})
}

func generate(leaves [][]D) *N {
	if len(leaves) == 0 {
		panic("cannot create node with no data")
	}
	if len(leaves) == 1 {
		return New(O{
			Data: leaves[0],
		})
	}

	return New(O{
		Left:  generate(leaves[:len(leaves)/2]),
		Right: generate(leaves[len(leaves)/2:]),
	})
}

func TestGenerate(t *testing.T) {
	type config struct {
		name   string
		leaves [][]D
		want   *N
	}

	configs := []config{
		{
			name: "Trivial",
			leaves: [][]D{
				{{ID: 1, AABB: interval(0, 100)}},
			},
			want: &N{
				data: []D{
					{ID: 1, AABB: interval(0, 100)},
				},
			},
		},
		func() config {
			leaves := [][]D{
				{{ID: 1, AABB: interval(0, 100)}},
				{{ID: 2, AABB: interval(101, 200)}},
			}

			left := generate([][]D{leaves[0]})
			right := generate([][]D{leaves[1]})
			root := &N{
				left:  left,
				right: right,
			}
			left.parent = root
			right.parent = root

			return config{
				name:   "TwoItems",
				leaves: leaves,
				want:   root,
			}
		}(),
		func() config {
			// We expect the following structure
			//
			//     A
			//    / \
			//   1   B
			//      / \
			//     2   3
			leaves := [][]D{
				{{ID: 1, AABB: interval(0, 100)}},
				{{ID: 2, AABB: interval(101, 200)}},
				{{ID: 3, AABB: interval(201, 300)}},
			}

			n1 := generate([][]D{leaves[0]})
			n2 := generate([][]D{leaves[1]})
			n3 := generate([][]D{leaves[2]})
			nb := &N{
				left:  n2,
				right: n3,
			}
			n2.parent = nb
			n3.parent = nb
			na := &N{
				left:  n1,
				right: nb,
			}
			n1.parent = na
			nb.parent = na

			return config{
				name:   "ThreeItems",
				leaves: leaves,
				want:   na,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := generate(c.leaves)
			if diff := cmp.Diff(
				c.want,
				got,
				cmp.AllowUnexported(
					N{},
					hyperrectangle.R{},
				),
			); diff != "" {
				t.Errorf("generate() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
