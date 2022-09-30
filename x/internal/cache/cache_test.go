package cache

import (
	"testing"
)

func TestInsert(t *testing.T) {
	for _, c := range []struct {
		name string
		c    *C[int]
		t    int
		want ID
	}{
		{
			name: "Empty",
			c: &C[int]{
				freed: map[ID]bool{},
			},
			t:    100,
			want: 0,
		},
		{
			name: "Freed",
			c: &C[int]{
				cache: make([]int, 1),
				freed: map[ID]bool{0: true},
			},
			t:    100,
			want: 0,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := c.c.Insert(c.t); got != c.want {
				t.Errorf("Insert() = %v, want = %v", got, c.want)
			}
			if got, _ := c.c.Get(c.want); got != c.t {
				t.Errorf("Get() = %v, want = %v", got, c.t)
			}
		})
	}
}

func BenchmarkInsert(b *testing.B) {
	c := New[int]()
	for i := 0; i < b.N; i++ {
		c.Insert(100)
	}
}
