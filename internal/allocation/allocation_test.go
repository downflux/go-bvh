package allocation

import (
	"testing"
)

type m struct{}

func TestAllocate(t *testing.T) {
	const n = 10000
	ids := map[ID]bool{}

	c := New[*m](nil)
	for i := 0; i < n; i++ {
		id := c.Allocate()
		if ids[id] {
			t.Errorf("allocating an already allocated node")
		}
		ids[id] = true
	}
}

func TestInsert(t *testing.T) {
	type config struct {
		name string
		c    C[*m]

		id   ID
		data *m

		succ bool
	}

	configs := []config{
		func() config {
			c := New[*m](nil)
			return config{
				name: "Trivial",
				c:    c,
				id:   c.Allocate(),
				data: &m{},
				succ: true,
			}
		}(),
		func() config {
			c := New[*m](nil)
			id := c.Allocate()
			c.Insert(id, &m{})

			return config{
				name: "Duplicate",
				c:    c,
				id:   id,
				data: &m{},
				succ: false,
			}
		}(),
		{
			name: "Unallocated",
			c:    New[*m](nil),
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
