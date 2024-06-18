package bitwise

// bit represents a single bit as a boolean.
type Bit bool

// New returns a new bit initialized to the specified value.
func New(value bool) Bit {
	return Bit(value)
}

// Invert flips the value of the bit.
func (b *Bit) Invert() {
	if b == nil {
		return
	}
	*b = !*b
}

// Value returns the current value of the bit.
func (b Bit) Value() uint8 {
	if bool(b) {
		return 1
	}

	return 0
}
