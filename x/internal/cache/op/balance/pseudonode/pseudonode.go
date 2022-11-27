package pseudonode

import (
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type I interface {
	AABB() hyperrectangle.M
	Heuristic() float64
	Height() int
}

type N struct {
	l   node.N
	r   node.N
	buf hyperrectangle.M
}

func New(l node.N, r node.N, buf hyperrectangle.M) *N {
	return &N{
		l:   l,
		r:   r,
		buf: buf,
	}
}

func (n N) AABB() hyperrectangle.M {
	n.buf.Copy(n.l.AABB().R())
	n.buf.Union(n.r.AABB().R())
	return n.buf
}

func (n N) Heuristic() float64 { return heuristic.H(n.AABB().R()) }

func (n N) Height() int {
	if n.l.Height() > n.r.Height() {
		return n.l.Height() + 1
	}
	return n.r.Height() + 1
}
