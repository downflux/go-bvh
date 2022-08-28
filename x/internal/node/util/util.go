package util

import (
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Equal(a *node.N, b *node.N) bool {
	if (a.IsLeaf() && !b.IsLeaf()) || (!a.IsLeaf() && b.IsLeaf()) {
		return false
	}
	if a.IsLeaf() && b.IsLeaf() {
		return cmp.Equal(
			a,
			b,
			cmpopts.IgnoreFields(
				node.N{},
				"left",
				"right",
				"parent",
				"aabbCacheIsValid",
				"aabbCache",
			),
			cmp.AllowUnexported(node.N{}, hyperrectangle.R{}),
		)
	}
	return (Equal(a.Left(), b.Left()) && Equal(a.Right(), b.Right())) || (Equal(a.Left(), b.Right()) && Equal(a.Right(), b.Left()))
}
