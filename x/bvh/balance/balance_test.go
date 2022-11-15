package balance

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

func TestB(t *testing.T) {
	type config struct {
		name string
		x    node.N
		data map[id.ID]hyperrectangle.R
		want node.N
	}

	configs := []config{}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := B(c.x, c.data, 1); !node.Equal(got, c.want) {
				t.Errorf("B() = %v, want = %v", got, c.want)
			}
		})
	}
}
