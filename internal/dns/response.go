package dns

import (
	"errors"
	"fmt"
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
	}
)

func parseResponse(bytes []byte) (*response, error) {
	raw := copyBytes(bytes)

	header, err := parseHeader(bytes)
	if err != nil {
		return nil, err
	}
	if header.IsTruncated() {
		return nil, errors.New("response is truncated")
	}
	if respCode := header.ResponseCode(); respCode != CodeNoError {
		return nil, fmt.Errorf("request responded with status: %s", respCode.CodeText())
	}

	bytes = bytes[headerSize:]
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
	return r.header.ID()
}

func (r *response) QueryCount() uint16 {
	return r.header.QueryCount()
}

func (r *response) AnswerCount() uint16 {
	return r.header.AnswerCount()
}

func copyBytes(bytes []byte) []byte {
	copy := make([]byte, len(bytes))
	for i := 0; i < len(bytes); i++ {
		copy[i] = bytes[i]
	}

	return copy
}
