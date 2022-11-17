package util

import (
	"fmt"
	"strings"

	"github.com/downflux/go-bvh/x/internal/cache/node"
)

func PreOrder(n node.N, f func(n node.N)) {
	f(n)
	if !n.IsLeaf() {
		PreOrder(n.Left(), f)
		PreOrder(n.Right(), f)
	}
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
