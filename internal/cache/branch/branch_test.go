package branch

import (
	"testing"
)

func TestSibling(t *testing.T) {
	type config struct {
		name string
		b    B
		want B
	}
	configs := []config{
		{
			name: "BLeft",
			b:    BLeft,
			want: BRight,
		},
		{
			name: "BRight",
			b:    BRight,
			want: BLeft,
		},
	}
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := c.b.Sibling(); got != c.want {
				t.Errorf("Sibling() = %v, want = %v", got, c.want)
			}
		})
	}
}
