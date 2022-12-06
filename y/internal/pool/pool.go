package pool

import (
	"fmt"

	"github.com/downflux/go-bvh/y/internal/node"
	"github.com/downflux/go-geometry/nd/vector"
)

type P struct {
	counter int

	k         vector.D
	n         int
	tolerance float64

	nodes map[node.ID]*node.N
	free  []*node.N
}

func (p *P) Insert() *node.N {
	if len(p.free) > 0 {
		var n *node.N
		n, p.free = p.free[len(p.free)-1], p.free[:len(p.free)-1]
		p.nodes[n.ID()] = n
		return n
	}

	n := node.New(p.k, p.n, p.tolerance, node.ID(p.counter))
	p.nodes[n.ID()] = n
	p.counter++

	return n
}

func (p *P) Remove(x node.ID) {
	n, ok := p.nodes[x]
	if !ok {
		panic(fmt.Sprintf("cannot find node %v", x))
	}

	node.Free(n)
	delete(p.nodes, x)
	p.free = append(p.free, n)
}

func (p *P) Get(x node.ID) *node.N {
	n, ok := p.nodes[x]
	if !ok {
		panic(fmt.Sprintf("cannot find node %v", x))
	}

	return n
}
