package dns

type (
	answer struct {
		question   question
		ttl        uint16
		readLength uint16
	}
)
