package id

import (
	"math/rand"
)

type ID uint64

func New() ID { return ID(rand.Uint64()) }
