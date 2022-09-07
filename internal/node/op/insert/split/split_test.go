package split

import (
	"testing"

	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/google/go-cmp/cmp"
)

func TestExecute(t *testing.T) {
	type result struct {
		n *node.N
		m *node.N
	}

	type config struct {
		name string
		p    P
		n    *node.N
		want result
	}

	configs := []config{}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := Execute(c.n, c.p)
			if diff := cmp.Diff(
				c.want, result{
					n: c.n,
					m: got,
				},
				cmp.AllowUnexported(result{}),
				cmp.Comparer(util.Equal),
			); diff != "" {
				t.Errorf("Execute() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
