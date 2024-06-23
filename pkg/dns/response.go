package dns

import (
	"errors"
	"fmt"
	"net"
)

type (
	response struct {
		raw     []byte
		header  *header
		queries []query
		answers []answer
	}

	Response interface {
		ID() ID
		QueryCount() int
		AnswerCount() int
		IPs() []net.IP
	}
)

func parseResponse(bytes []byte) (*response, error) {
	raw := copyBytes(bytes)

	header, err := parseHeader(bytes)
	if err != nil {
		return nil, err
	}

	flags := header.flags
	if flags.isTruncated {
		return nil, errors.New("response is truncated")
	}
	if respCode := flags.responseCode; respCode != CodeNoError {
		return nil, fmt.Errorf("request responded with status: %s", respCode.CodeText())
	}

	bytes = bytes[headerSize:]
	qCount := header.queryCount
	queries, qBytes, err := parseQueries(bytes, int(qCount))
	if err != nil {
		return nil, err
	}

	bytes = bytes[qBytes:]
	aCount := header.answerCount
	answers, _, err := parseAnswers(bytes, int(aCount))
	if err != nil {
		return nil, err
	}

	return &response{
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

func (r *response) ID() ID {
	return r.header.id
}

func (r *response) QueryCount() int {
	return int(r.header.queryCount)
}

func (r *response) AnswerCount() int {
	return int(r.header.answerCount)
}

func (r *response) IPs() []net.IP {
	var ips []net.IP
	for _, a := range r.answers {
		ips = append(ips, a.IP())
	}

	return ips
}

func copyBytes(bytes []byte) []byte {
	copy := make([]byte, len(bytes))
	for i := 0; i < len(bytes); i++ {
		copy[i] = bytes[i]
	}

	return copy
}
