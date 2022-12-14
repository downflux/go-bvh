package impl

import (
	"testing"

	"github.com/downflux/go-bvh/internal/cache/node"
	"github.com/downflux/go-bvh/internal/cache/node/util/cmp"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/internal/cache/id"
)

var (
	_ node.N = &N{}
	_ A      = &MockCache{}
)

type MockCache struct {
	data map[cid.ID]*N
}

func (m *MockCache) K() vector.D   { return vector.D(2) }
func (m *MockCache) LeafSize() int { return 1 }

func (m *MockCache) Get(x cid.ID) (node.N, bool) {
	n, ok := m.data[x]
	return n, ok
}

func TestSetHeight(t *testing.T) {
	type config struct {
		name string
		n    *N
		want int
	}

	configs := []config{
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}

			n := New(c, 0)
			n.Allocate(-1, -1, -1)

			c.data[n.ID()] = n

			return config{
				name: "Leaf",
				n:    n,
				want: 0,
			}
		}(),
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}

			n := New(c, 0)
			n.Allocate(-1, -1, -1)

			l := New(c, 1)
			l.Allocate(n.ID(), -1, -1)

			r := New(c, 2)
			r.Allocate(n.ID(), -1, -1)

			n.SetLeft(l.ID())
			n.SetRight(r.ID())

			c.data[n.ID()] = n
			c.data[l.ID()] = l
			c.data[r.ID()] = r

			return config{
				name: "NonLeaf",
				n:    n,
				want: 1,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			node.SetHeight(c.n)
			if got := c.n.Height(); got != c.want {
				t.Errorf("Height() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	type config struct {
		name string
		n    *N
		m    *N
		want bool
	}

	configs := []config{
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}
			n := New(c, 0)
			n.Allocate(-1, -1, -1)

			m := New(c, 0)
			m.Allocate(-1, -1, -1)

			return config{
				name: "ID",
				n:    n,
				m:    m,
				want: true,
			}
		}(),
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}
			n := New(c, 0)
			n.Allocate(-1, -1, -1)

			m := New(c, 1)
			m.Allocate(-1, -1, -1)

			return config{
				name: "ID/NotEqual",
				n:    n,
				m:    m,
				want: false,
			}
		}(),
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}
			n := New(c, 0)
			n.Allocate(-1, -1, -1)
			n.SetHeight(10)

			m := New(c, 0)
			m.Allocate(-1, -1, -1)
			m.SetHeight(10)

			return config{
				name: "Height",
				n:    n,
				m:    m,
				want: true,
			}
		}(),
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}
			n := New(c, 0)
			n.Allocate(-1, -1, -1)
			n.SetHeight(10)

			m := New(c, 0)
			m.Allocate(-1, -1, -1)
			m.SetHeight(11)

			return config{
				name: "Height/NotEqual",
				n:    n,
				m:    m,
				want: false,
			}
		}(),
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}

			c.data[0] = New(c, 0)
			c.data[0].Allocate(-1, -1, -1)

			n := New(c, 1)
			n.Allocate(0, -1, -1)

			m := New(c, 1)
			m.Allocate(0, -1, -1)
			return config{
				name: "Parent",
				n:    n,
				m:    m,
				want: true,
			}
		}(),
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}

			c.data[0] = New(c, 0)
			c.data[0].Allocate(-1, -1, -1)

			c.data[1] = New(c, 1)
			c.data[1].Allocate(-1, -1, -1)

			n := New(c, 2)
			n.Allocate(0, -1, -1)

			m := New(c, 2)
			m.Allocate(1, -1, -1)

			return config{
				name: "Parent/NotEqual",
				n:    n,
				m:    m,
				want: false,
			}
		}(),
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}

			c.data[1] = New(c, 1)
			c.data[1].Allocate(0, -1, -1)

			c.data[2] = New(c, 2)
			c.data[2].Allocate(0, -1, -1)

			n := New(c, 0)
			n.Allocate(-1, 1, 2)
			n.SetHeight(1)

			m := New(c, 0)
			m.Allocate(-1, 1, 2)
			m.SetHeight(1)

			return config{
				name: "Child",
				n:    n,
				m:    m,
				want: true,
			}
		}(),
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}

			c.data[1] = New(c, 1)
			c.data[1].Allocate(0, -1, -1)

			c.data[2] = New(c, 2)
			c.data[2].Allocate(0, -1, -1)

			c.data[3] = New(c, 3)
			c.data[3].Allocate(0, -1, -1)

			n := New(c, 0)
			n.Allocate(-1, 1, 2)
			n.SetHeight(1)

			m := New(c, 0)
			m.Allocate(-1, 1, 3)
			m.SetHeight(1)

			return config{
				name: "Child/NotEqual",
				n:    n,
				m:    m,
				want: false,
			}
		}(),
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}
			n := New(c, 0)
			n.Allocate(-1, -1, -1)
			n.Leaves()[0] = struct{}{}

			m := New(c, 0)
			m.Allocate(-1, -1, -1)
			m.Leaves()[0] = struct{}{}

			return config{
				name: "Leaves",
				n:    n,
				m:    m,
				want: true,
			}
		}(),
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}
			n := New(c, 0)
			n.Allocate(-1, -1, -1)
			n.Leaves()[0] = struct{}{}

			m := New(c, 0)
			m.Allocate(-1, -1, -1)
			m.Leaves()[1] = struct{}{}

			return config{
				name: "Leaves/NotEqual",
				n:    n,
				m:    m,
				want: false,
			}
		}(),
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}
			n := New(c, 0)
			n.Allocate(-1, -1, -1)
			n.Leaves()[0] = struct{}{}

			m := New(c, 0)
			m.Allocate(-1, -1, -1)

			return config{
				name: "Leaves/Len/NotEqual",
				n:    n,
				m:    m,
				want: false,
			}
		}(),
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}
			n := New(c, 0)
			n.Allocate(-1, -1, -1)
			n.AABB().Copy(*hyperrectangle.New(
				vector.V([]float64{1, 1}),
				vector.V([]float64{2, 2}),
			))

			m := New(c, 0)
			m.Allocate(-1, -1, -1)
			m.AABB().Copy(*hyperrectangle.New(
				vector.V([]float64{1, 1}),
				vector.V([]float64{2, 2}),
			))

			return config{
				name: "AABB",
				n:    n,
				m:    m,
				want: true,
			}
		}(),
		func() config {
			c := &MockCache{
				data: map[cid.ID]*N{},
			}
			n := New(c, 0)
			n.Allocate(-1, -1, -1)
			n.AABB().Copy(*hyperrectangle.New(
				vector.V([]float64{1, 1}),
				vector.V([]float64{2, 2}),
			))

			m := New(c, 0)
			m.Allocate(-1, -1, -1)
			m.AABB().Copy(*hyperrectangle.New(
				vector.V([]float64{2, 2}),
				vector.V([]float64{3, 3}),
			))

			return config{
				name: "AABB/NotEqual",
				n:    n,
				m:    m,
				want: false,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := cmp.Equal(c.n, c.m); got != c.want {
				t.Errorf("Equal() = %v, want = %v", got, c.want)
			}
		})
	}
}
