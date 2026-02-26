package utils

import (
	"github.com/nrednav/cuid2"
)

func GenerateID() string {
	gen, _ := cuid2.Init(cuid2.WithLength(32))
	return gen()
}

func IsValidID(id string) bool {
	return cuid2.IsCuid(id)
}
