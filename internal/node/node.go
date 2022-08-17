package node

import (
	"fmt"
	"math/rand"

	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type Nodes map[Index]*N

func (ns Nodes) Allocate() Index {
	var i Index
	found := true
	for ; found; i = Index(rand.Int()) {
		_, found = ns[i]
	}
	ns[i] = nil

	return i
}

func (ns Nodes) Insert(i Index, n *N) {
	m, ok := ns[i]

	// Index must be allocated first.
	if !ok {
		panic(fmt.Sprintf("inserting an unallocated node %v", i))
	}

	if m != nil {
		panic(fmt.Sprintf("duplicate node found with same index %v", i))
	}

	ns[i] = n
}

type Index int

type O struct {
	ID point.ID

	Nodes  Nodes
	Index  Index
	Parent Index
	Left   Index
	Right  Index

	Bound hyperrectangle.R
}

type N struct {
	id point.ID

	nodes  Nodes
	index  Index
	parent Index
	left   Index
	right  Index

	bound hyperrectangle.R
}

func New(o O) *N {
	return &N{
		id: o.ID,

		nodes:  o.Nodes,
		index:  o.Index,
		parent: o.Parent,
		left:   o.Left,
		right:  o.Right,

		bound: o.Bound,
	}
}

func (n *N) Bound() hyperrectangle.R { return n.bound }

func (n *N) Leaf() bool   { return n.Left() == nil && n.Right() == nil }
func (n *N) Index() Index { return n.index }
func (n *N) Left() *N     { return n.nodes[n.left] }
func (n *N) Right() *N    { return n.nodes[n.right] }
func (n *N) Parent() *N   { return n.nodes[n.parent] }
func (n *N) ID() point.ID { return n.id }

func (n *N) Move(id point.ID, offset vector.V) bool { return false }
func (n *N) Remove(id point.ID) bool                { return false }
