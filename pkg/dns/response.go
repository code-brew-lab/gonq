package dns

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strings"
)

type (
	Response struct {
		raw     []byte
		header  *Header
		queries []query
		answers []answer
	}
)

func (r Response) Raw() []byte {
	return r.raw
}

func (r Response) Header() *Header {
	return r.header
}

func (r Response) IPs() []net.IP {
	var ips []net.IP
	for _, a := range r.answers {
		ips = append(ips, a.ip)
	}

	return ips
}

func (r Response) String() string {
	var queries []string
	for _, q := range r.queries {
		queries = append(queries, q.string())
	}

	var answers []string
	for _, a := range r.answers {
		answers = append(answers, a.string())
	}

	return fmt.Sprintf("Raw: %v\nHeader: %v\nQueries: [%v]\nAnswers: [%v]\n",
		hex.EncodeToString(r.raw), r.header, strings.Join(queries, ", "), strings.Join(answers, ", "))
}

func parseResponse(bytes []byte) (*Response, error) {
	raw := copyBytes(bytes)

	header, err := ParseHeader(bytes)
	if err != nil {
		return nil, err
	}

	flags := header.Flags()
	if flags.IsTruncated() {
		return nil, errors.New("response is truncated")
	}
	if respCode := flags.ResponseCode(); respCode != CodeNoError {
		return nil, fmt.Errorf("request responded with status: %s", respCode.CodeText())
	}

	bytes = bytes[HeaderSize:]
	qCount := header.QueryCount()
	queries, qBytes, err := parseQueries(bytes, int(qCount))
	if err != nil {
		return nil, err
	}

	bytes = bytes[qBytes:]
	aCount := header.AnswerCount()
	answers, _, err := parseAnswers(bytes, int(aCount))
	if err != nil {
		return nil, err
	}

	return &Response{
		raw:     raw,
		header:  header,
		queries: queries,
		answers: answers,
	}, nil
}

func parseQueries(bytes []byte, qCount int) ([]query, int, error) {
	queries := make([]query, 0)
	raw := copyBytes(bytes)
	read := 0

	for i := 0; i < int(qCount); i++ {
		q, bytesRead, err := parseQuery(raw)
		if err != nil {
			return nil, 0, err
		}

		queries = append(queries, q)
		if bytesRead > len(bytes) {
			return nil, 0, errors.New("byte slice length exceeded")
		}

		raw = raw[bytesRead:]
		read += bytesRead
	}

	return queries, read, nil
}

func parseAnswers(bytes []byte, aCount int) ([]answer, int, error) {
	answers := make([]answer, 0)
	raw := copyBytes(bytes)
	read := 0

	for i := 0; i < int(aCount); i++ {
		a, bytesRead, err := parseAnswer(raw)
		if err != nil {
			return nil, 0, err
		}

		answers = append(answers, a)
		if bytesRead > len(bytes) {
			return nil, 0, errors.New("byte slice length exceeded")
		}

		raw = raw[bytesRead:]
		read += bytesRead
	}

	return answers, read, nil
}

func copyBytes(bytes []byte) []byte {
	copy := make([]byte, len(bytes))
	for i := 0; i < len(bytes); i++ {
		copy[i] = bytes[i]
	}

	return copy
}
