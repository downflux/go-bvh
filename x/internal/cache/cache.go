package cache

type ID int

type C[T any] struct {
	cache []T
	freed map[ID]bool
}

func New[T any]() *C[T] {
	return &C[T]{
		cache: make([]T, 0, 128),
		freed: map[ID]bool{},
	}
}

func (c *C[T]) Len() int { return len(c.cache) - len(c.freed) }

func (c *C[T]) Get(x ID) (T, bool) {
	if int(x) >= len(c.cache) || c.freed[x] {
		var blank T
		return blank, false
	}

	return c.cache[x], true
}

func (c *C[T]) Insert(t T) ID {
	if len(c.freed) > 0 {
		for x := range c.freed {
			delete(c.freed, x)
			c.cache[x] = t
			return x
		}
	}

	c.cache = append(c.cache, t)
	return ID(len(c.cache) - 1)
}

func (c *C[T]) Remove(x ID) bool {
	if int(x) >= len(c.cache) || c.freed[x] {
		return false
	}

	var blank T
	c.cache[x] = blank
	c.freed[x] = true

	return true
}
