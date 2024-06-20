package dns

type (
	Response struct {
		Header    *Header
		questions []question
		answers   []answer
	}
)

func ParseResponse(response []byte) (*Response, error) {
	return &Response{}, nil
}
