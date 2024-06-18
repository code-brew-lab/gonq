package bitwise

// Bit represents a single bit as a uint8.
type Bit uint8

const (
	Zero Bit = 0
	One  Bit = 1
)

// New returns a new bit initialized to the specified value.
func New(value bool) Bit {
	if value {
		return One
	}
	return Zero
}

// Invert flips the value of the bit.
func (b *Bit) Invert() {
	if b == nil {
		return
	}
	*b = 1 - *b
}

// Value returns the current value of the bit.
func (b Bit) Value() uint8 {
	return uint8(b)
}
