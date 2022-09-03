package insert

import (
	"testing"

	// "github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/util"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
	// nid "github.com/downflux/go-bvh/x/internal/node/id"
)

func TestSibling(t *testing.T) {
	type config struct {
		name string
		root *node.N
		aabb hyperrectangle.R
		want *node.N
	}

	configs := []config{}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := sibling(c.root, c.aabb)
			if diff := cmp.Diff(c.want, got, cmp.Comparer(util.Equal)); diff != "" {
				t.Errorf("sibling() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
