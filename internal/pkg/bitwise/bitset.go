package bitwise

import (
	"errors"
)

type BitSet struct {
	bits []Bit
}

// NewSet returns a new, empty BitSet.
func NewSet() *BitSet {
	return &BitSet{make([]Bit, 0)}
}

// Add adds a Bit to the BitSet.
func (bs *BitSet) Add(bit Bit) {
	if bs == nil {
		return
	}

	bs.bits = append(bs.bits, bit)
}

func (bs *BitSet) AddSet(bits []Bit) {
	if bs == nil {
		return
	}

	for i := 0; i < len(bits); i++ {
		bs.bits = append(bs.bits, bits[i])
	}
}

// ToBytes converts the BitSet to a byte slice.
func (bs *BitSet) ToBytes() ([]byte, error) {
	if bs == nil {
		return nil, errors.New("nil bitset object")
	}

	bitSet := bs.bits
	byteCount := (len(bitSet) + 7) / 8 // Calculate the number of bytes needed to represent the bits
	bytes := make([]byte, byteCount)

	for i, bit := range bitSet {
		if bit == One {
			byteIndex := uint8(i / 8)      // Determine which byte this bit belongs to
			bitIndex := uint8(7 - (i % 8)) // Determine the position of the bit within the byte (from left)
			byte := uint8(1 << bitIndex)   // Create a byte with the bit set at the correct position
			bytes[byteIndex] |= byte       // Set the bit in the corresponding byte
		}
	}

	return bytes, nil
}
