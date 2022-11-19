package util

import (
	"fmt"
	"strings"

	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/heuristic"
)

func PostOrder(n node.N, f func(n node.N)) {
	if !n.IsLeaf() {
		PostOrder(n.Left(), f)
		PostOrder(n.Right(), f)
	}
	f(n)
}

func PreOrder(n node.N, f func(n node.N)) {
	f(n)
	if !n.IsLeaf() {
		PreOrder(n.Left(), f)
		PreOrder(n.Right(), f)
	}
}

// SAH returns the surface area heuristic as defined in Macdonald and Booth
// 1990.
//
// The total heuristic is comprised of three separate components -- the cost of
// the internal nodes, the cost of the leaves, and the cost of testing for
// intersections. We use track these via ci, cl, and co respectively.
//
// Per Aila et al., a "normal" SAH value is around 100.
func SAH(n node.N) float64 {
	var ci, cl, co float64
	PreOrder(n, func(n node.N) {
		if !n.IsLeaf() {
			ci += heuristic.H(n.AABB().R())
		} else {
			cl += heuristic.H(n.AABB().R())
			co += heuristic.H(n.AABB().R()) * float64(len(n.Leaves()))
		}
	})
	return (1.2*ci + 1.0*cl + 0*co) / heuristic.H(n.AABB().R())
}

func S(n node.N) string {
	var s []string

	PreOrder(n, func(n node.N) {
		if n.IsLeaf() {
			leaves := []string{}
			for x := range n.Leaves() {
				leaves = append(leaves, fmt.Sprint(x))
			}

			s = append(s, fmt.Sprintf(
				"ID: %v, AABB: %v, Height: %v, Data: %v",
				n.ID(),
				n.AABB(),
				n.Height(),
				strings.Join(leaves, ","),
			))

		} else {
			s = append(s, fmt.Sprintf(
				"ID: %v, AABB: %v, Height: %v, Left: %v, Right: %v",
				n.ID(),
				n.AABB(),
				n.Height(),
				n.Left().ID(),
				n.Right().ID(),
			))
		}
	})
	return strings.Join(s, "\n")
}
