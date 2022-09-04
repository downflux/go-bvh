package id

import (
	"math/rand"
)

type ID uint64

func Generate() ID         { return ID(rand.Uint64()) }
func Increment(id ID) ID   { return ID(uint64(id) + 1) }
func (id ID) IsZero() bool { return id == ID(0) }
