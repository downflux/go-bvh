package bvh

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

func TestAVL(t *testing.T) {
	type config struct {
		name    string
		x       node.N
		data    map[id.ID]hyperrectangle.R
		epsilon float64
		want    node.N
	}

	configs := []config{}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := avl(c.x, c.data, c.epsilon); !node.Equal(got, c.want) {
				t.Errorf("avl() = %v, want = %v", got, c.want)
			}
		})
	}
}
