package dns

import "math/rand"

type (
	ID uint16
)

func NewID() ID {
	id := uint16(rand.Int31())
	return ID(id)
}

func (i ID) toUint16() uint16 {
	return uint16(i)
}
