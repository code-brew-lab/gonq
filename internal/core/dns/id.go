package dns

import (
	"encoding/binary"
	"errors"
	"math/rand/v2"
)

type (
	ID uint16
)

func NewID() ID {
	id := uint16(rand.Uint32())
	return ID(id)
}

func ParseID(bytes []byte) (ID, error) {
	if len(bytes) != 2 {
		return ID(0), errors.New("id should be 2 bytes")
	}

	be := binary.BigEndian
	id := be.Uint16(bytes[:])

	return ID(id), nil
}

func (i ID) toUint16() uint16 {
	return uint16(i)
}
