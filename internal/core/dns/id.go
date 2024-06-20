package dns

import "math/rand/v2"

type (
	ID uint16
)

func NewID() ID {
	id := uint16(rand.Uint32())
	return ID(id)
}

func (i ID) toUint16() uint16 {
	return uint16(i)
}
