package id

const (
	IDInvalid = ID(-1)
)

type ID int

func (x ID) IsValid() bool { return x > IDInvalid }
