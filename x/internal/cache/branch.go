package cache

type B int

const (
	BInvalid B = iota
	BLeft
	BRight
)

func (b B) IsValid() bool { return b == BLeft || b == BRight }
func (b B) Sibling() B    { return ((b - 1) ^ 1) + 1 }
