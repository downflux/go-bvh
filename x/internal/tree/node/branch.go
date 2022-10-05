package node

type B int

const (
	BLeft B = iota
	BRight
	BInvalid
)

func (b B) IsValid() bool { return b == BLeft || b == BRight }
func (b B) Sibling() B    { return b ^ 1 }
