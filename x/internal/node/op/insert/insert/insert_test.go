package insert

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/util"
	"github.com/google/go-cmp/cmp"
)

func TestExecute(t *testing.T) {
	type config struct {
		name string
		n    *node.N
		m    *node.N
		want *node.N
	}

	configs := []config{}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Execute(c.n, c.m)
			if diff := cmp.Diff(c.want, got, cmp.Comparer(util.Equal)); diff != "" {
				t.Errorf("Execute() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
