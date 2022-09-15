package node

import (
	"fmt"
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSet(t *testing.T) {
	type result struct {
		height uint
		aabb   hyperrectangle.R
	}

	type config struct {
		name string
		n    *N
		want result
	}

	configs := []config{
		func() config {
			c := Cache()
			n := New(O{
				Nodes: c,
				ID:    100,
				Right: 101,
				// Left is unset.
				Size: 1,
				K:    1,
			})
			New(O{
				Nodes:  c,
				ID:     101,
				Parent: 100,
				Data: map[id.ID]hyperrectangle.R{
					1: *hyperrectangle.New([]float64{0}, []float64{25}),
				},
				Size: 1,
				K:    1,
			})
			m := New(O{
				Nodes:  c,
				ID:     102,
				Parent: 100,
				Left:   103,
				Right:  104,
				Size:   1,
				K:      1,
			})
			New(O{
				Nodes:  c,
				ID:     103,
				Parent: 102,
				Data: map[id.ID]hyperrectangle.R{
					2: *hyperrectangle.New([]float64{26}, []float64{50}),
				},
				Size: 1,
				K:    1,
			})
			New(O{
				Nodes:  c,
				ID:     104,
				Parent: 102,
				Data: map[id.ID]hyperrectangle.R{
					3: *hyperrectangle.New([]float64{51}, []float64{75}),
				},
				Size: 1,
				K:    1,
			})
			n.SetLeft(m)
			return config{
				name: "SetLeft",
				n:    n,
				want: result{
					height: 3,
					aabb:   *hyperrectangle.New([]float64{0}, []float64{75}),
				},
			}
		}(),
		func() config {
			c := Cache()
			n := New(O{
				Nodes: c,
				ID:    100,
				// Right is unset.
				Left: 101,
				Size: 1,
				K:    1,
			})
			New(O{
				Nodes:  c,
				ID:     101,
				Parent: 100,
				Data: map[id.ID]hyperrectangle.R{
					1: *hyperrectangle.New([]float64{0}, []float64{25}),
				},
				Size: 1,
				K:    1,
			})
			m := New(O{
				Nodes:  c,
				ID:     102,
				Parent: 100,
				Left:   103,
				Right:  104,
				Size:   1,
				K:      1,
			})
			New(O{
				Nodes:  c,
				ID:     103,
				Parent: 102,
				Data: map[id.ID]hyperrectangle.R{
					2: *hyperrectangle.New([]float64{26}, []float64{50}),
				},
				Size: 1,
				K:    1,
			})
			New(O{
				Nodes:  c,
				ID:     104,
				Parent: 102,
				Data: map[id.ID]hyperrectangle.R{
					3: *hyperrectangle.New([]float64{51}, []float64{75}),
				},
				Size: 1,
				K:    1,
			})
			n.SetRight(m)
			return config{
				name: "SetRight",
				n:    n,
				want: result{
					height: 3,
					aabb:   *hyperrectangle.New([]float64{0}, []float64{75}),
				},
			}
		}(),
		func() config {
			c := Cache()
			m := New(O{
				Nodes: c,
				ID:    100,
				Left:  101,
				Right: 102,
				Size:  1,
				K:     1,
			})
			New(O{
				Nodes:  c,
				ID:     101,
				Parent: 100,
				Data: map[id.ID]hyperrectangle.R{
					1: *hyperrectangle.New([]float64{0}, []float64{25}),
				},
				Size: 1,
				K:    1,
			})
			n := New(O{
				Nodes: c,
				ID:    102,
				// Parent is unset.
				Left:  103,
				Right: 104,
				Size:  1,
				K:     1,
			})
			New(O{
				Nodes:  c,
				ID:     103,
				Parent: 102,
				Data: map[id.ID]hyperrectangle.R{
					2: *hyperrectangle.New([]float64{26}, []float64{50}),
				},
				Size: 1,
				K:    1,
			})
			New(O{
				Nodes:  c,
				ID:     104,
				Parent: 102,
				Data: map[id.ID]hyperrectangle.R{
					3: *hyperrectangle.New([]float64{51}, []float64{75}),
				},
				Size: 1,
				K:    1,
			})
			n.SetParent(m)
			return config{
				name: "SetParent",
				n:    m,
				want: result{
					height: 3,
					aabb:   *hyperrectangle.New([]float64{0}, []float64{75}),
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(fmt.Sprintf("%v/AABB", c.name), func(t *testing.T) {
			if diff := cmp.Diff(
				c.want.aabb,
				c.n.AABB(),
				cmp.AllowUnexported(hyperrectangle.R{}),
			); diff != "" {
				t.Errorf("AABB() mismatch (-want +got):\n%v", diff)
			}

		})
		t.Run(fmt.Sprintf("%v/Height", c.name), func(t *testing.T) {
			if c.want.height != c.n.Height() {
				t.Errorf("Height() = %v, want = %v", c.n.Height(), c.want.height)
			}

		})
	}
}

func TestBroadPhase(t *testing.T) {
	type config struct {
		name string
		n    *N
		q    hyperrectangle.R
		want []id.ID
	}

	configs := []config{
		func() config {
			c := Cache()
			root := New(O{
				Nodes: c,
				ID:    100,
				Left:  101,
				Right: 102,
				Size:  2,
				K:     1,
			})
			New(O{
				Nodes:  c,
				ID:     101,
				Parent: 100,
				Data: map[id.ID]hyperrectangle.R{
					1: *hyperrectangle.New([]float64{0}, []float64{24}),
					2: *hyperrectangle.New([]float64{25}, []float64{49}),
				},
				Size: 2,
				K:    1,
			})
			New(O{
				Nodes:  c,
				ID:     102,
				Parent: 100,
				Data: map[id.ID]hyperrectangle.R{
					3: *hyperrectangle.New([]float64{50}, []float64{74}),
					4: *hyperrectangle.New([]float64{75}, []float64{99}),
				},
				Size: 2,
				K:    1,
			})

			return config{
				name: "Internal",
				n:    root,
				q:    *hyperrectangle.New([]float64{26}, []float64{73}),
				want: []id.ID{2, 3},
			}
		}(),
		{
			name: "Leaf/Overlaps",
			n: New(O{
				Nodes: Cache(),
				Data: map[id.ID]hyperrectangle.R{
					1: *hyperrectangle.New([]float64{51}, []float64{100}),
					2: *hyperrectangle.New([]float64{0}, []float64{50}),
				},
				Size: 2,
				K:    1,
			}),
			q:    *hyperrectangle.New([]float64{1}, []float64{99}),
			want: []id.ID{1, 2},
		},
		{
			name: "Leaf/Disjoint",
			n: New(O{
				Nodes: Cache(),
				Data: map[id.ID]hyperrectangle.R{
					1: *hyperrectangle.New([]float64{51}, []float64{100}),
					2: *hyperrectangle.New([]float64{0}, []float64{50}),
				},
				Size: 2,
				K:    1,
			}),
			q:    *hyperrectangle.New([]float64{100.1}, []float64{100.2}),
			want: []id.ID{},
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := c.n.BroadPhase(c.q)
			if diff := cmp.Diff(
				c.want,
				got,
				cmpopts.SortSlices(
					func(a, b id.ID) bool { return a < b },
				),
			); diff != "" {
				t.Errorf("BroadPhase() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
