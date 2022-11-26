package util

import (
	"fmt"
	"math"
	"strings"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
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

func ValidateOrDie(c *cache.C, data map[id.ID]hyperrectangle.R, n node.N) {
	if err := Validate(c, data, n); err != nil {
		panic(fmt.Errorf("encountered validation error on node %v: %v", n.ID(), err))
	}
}

func Validate(c *cache.C, data map[id.ID]hyperrectangle.R, n node.N) error {
	var err error
	buf := hyperrectangle.New(
		vector.V(make([]float64, n.AABB().Min().Dimension())),
		vector.V(make([]float64, n.AABB().Min().Dimension())),
	).M()

	PostOrder(n, func(n node.N) {
		if err != nil {
			return
		}

		if n.IsLeaf() {
			l := len(n.Leaves())
			if l == 0 {
				err = fmt.Errorf("leaf node %v has no child objects", n.ID())
				return
			}
			if l > c.LeafSize() {
				err = fmt.Errorf("leaf node %v has too many child objects", n.ID())
				return
			}

			initialized := false
			for x := range n.Leaves() {
				if !initialized {
					initialized = true
					buf.Copy(data[x])
				} else {
					buf.Union(data[x])
				}
			}
		} else {
			if h := int(math.Max(
				float64(n.Left().Height()),
				float64(n.Right().Height()),
			)) + 1; h != n.Height() {
				err = fmt.Errorf("parent node %v height does not match expected", n.ID())
				return
			}

			buf.Copy(n.Left().AABB().R())
			buf.Union(n.Right().AABB().R())

		}
		if !hyperrectangle.Contains(n.AABB().R(), buf.R()) {
			err = fmt.Errorf("parent node %v does not wholly encapsulate its children", n.ID())
			return
		}
	})
	return err
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
