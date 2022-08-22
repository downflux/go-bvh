package allocation

import (
	"testing"

	"github.com/downflux/go-bvh/internal/allocation/id"
)

type m struct{}

func TestAllocate(t *testing.T) {
	const n = 10000
	ids := map[id.ID]bool{}

	c := *New[*m]()
	for i := 0; i < n; i++ {
		i := c.Allocate()
		if ids[i] {
			t.Errorf("allocating an already allocated node")
		}
		ids[i] = true
	}
}

func TestInsert(t *testing.T) {
	type config struct {
		name string
		c    C[*m]

		id   id.ID
		data *m

		succ bool
	}

	configs := []config{
		func() config {
			c := *New[*m]()
			return config{
				name: "Trivial",
				c:    c,
				id:   c.Allocate(),
				data: &m{},
				succ: true,
			}
		}(),
		func() config {
			c := *New[*m]()
			i := c.Allocate()
			c.Insert(i, &m{})

			return config{
				name: "Duplicate",
				c:    c,
				id:   i,
				data: &m{},
				succ: false,
			}
		}(),
		{
			name: "Unallocated",
			c:    *New[*m](),
			id:   0,
			data: &m{},
			succ: false,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			err := c.c.Insert(c.id, c.data)
			if err != nil && c.succ {
				t.Errorf("Insert() = %v, want = nil", err)
			}
			if err == nil && !c.succ {
				t.Errorf("Insert() = nil, want a non-nil error")
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type config struct {
		name string

		c  C[*m]
		id id.ID

		succ bool
	}

	configs := []config{
		{
			name: "RemoveExists",
			c: C[*m]{
				1: &m{},
			},
			id:   1,
			succ: true,
		},
		{
			name: "RemoveDNE",
			c: C[*m]{
				1: &m{},
			},
			id:   2,
			succ: false,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			err := c.c.Remove(c.id)
			if err != nil && c.succ {
				t.Errorf("Remove() = %v, want = nil", err)
			}
			if err == nil && !c.succ {
				t.Errorf("Remove() = nil, want a non-nil error")
			}

			if _, ok := c.c[c.id]; ok {
				t.Errorf("c[i] = _, true, want = false")
			}
		})
	}
}
