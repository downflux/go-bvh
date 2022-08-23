package move

import (
	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/allocation/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/insert"
	"github.com/downflux/go-bvh/internal/node/remove"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bvh "github.com/downflux/go-bvh/hyperrectangle"
)

// Execute attempts to move a node.
//
// Execute returns the tuple (new, root) allocation IDs.
func Execute(nodes allocation.C[*node.N], i id.ID, r hyperrectangle.R) (id.ID, id.ID) {
	n := nodes[i]
	rid := node.Root(nodes, n).Index()
	if !bvh.Contains(n.Bound(), r) {
		rid = remove.Execute(nodes, i)
		return insert.Execute(nodes, rid, n.ID(), r)
	}
	return n.Index(), rid
}
