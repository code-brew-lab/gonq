package dns

import (
	"fmt"
	"strings"
)

type (
	Response struct {
		Header    *Header
		questions []question
		answers   []answer
	}
)

func ParseResponse(response []byte) (*Response, error) {
	header, err := parseHeader(response[0:12])
	if err != nil {
		return nil, fmt.Errorf("dns.ParseResponse: %v", err)
	}

	return &Response{
		Header: header,
	}, nil
}

func (r *Response) String() string {
	var sb strings.Builder
	sb.WriteString("Response {\n")
	if r.Header != nil {
		sb.WriteString(fmt.Sprintf("%s\n", r.Header.string(1, "\t")))
	} else {
		sb.WriteString("  Header: nil\n")
	}
	sb.WriteString("}")
	return sb.String()
}
