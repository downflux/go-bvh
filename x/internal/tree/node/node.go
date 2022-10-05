package node

import (
	"github.com/downflux/go-bvh/x/internal/cache"
)

type N struct {
	cache *cache.C
	id    cache.ID
}

func New(cache *cache.C, x cache.ID) *N {
	return &N{
		cache: cache,
		id:    x,
	}
}

func (n *N) Parent() *N {
	x := n.cache.GetOrDie(n.id).Parent()
	if _, ok := n.cache.Get(x); !ok {
		return nil
	}

	return New(n.cache, x)
}

func (n *N) Right() *N {
	x := n.cache.GetOrDie(n.id).Right()
	if _, ok := n.cache.Get(x); !ok {
		return nil
	}

	return New(n.cache, x)
}

func (n *N) Left() *N {
	x := n.cache.GetOrDie(n.id).Left()
	if _, ok := n.cache.Get(x); !ok {
		return nil
	}

	return New(n.cache, x)
}
