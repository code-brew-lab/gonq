package dns

import (
	"errors"
	"fmt"
	"strings"
)

type Response struct {
	Raw     []byte
	Header  *Header
	queries []query
	answers []answer
}

func ParseResponse(response []byte) (*Response, error) {
	raw := copyBytes(response)

	header, err := parseHeader(response)
	if err != nil {
		return nil, fmt.Errorf("dns.ParseResponse: %v", err)
	}
	if header.IsTruncated() {
		return nil, errors.New("dns.ParseResponse: response is truncated")
	}
	if respCode := header.ResponseCode(); respCode != CodeNoError {
		return nil, fmt.Errorf("dns.ParseResponse: request responded with status: %s", respCode.CodeText())
	}

	response = response[headerSize:]
	qCount := header.QuestionCount()
	queries, qBytes, err := parseQueries(response, int(qCount))
	if err != nil {
		return nil, fmt.Errorf("dns.ParseResponse: %v", err)
	}

	response = response[qBytes:]
	aCount := header.AnswerCount()
	answers, _, err := parseAnswers(response, int(aCount))
	if err != nil {
		return nil, fmt.Errorf("dns.ParseResponse: %v", err)
	}

	return &Response{
		Raw:     raw,
		Header:  header,
		queries: queries,
		answers: answers,
	}, nil
}

func (r *Response) String() string {
	var sb strings.Builder
	sb.WriteString("Response:\n")
	if r.Header != nil {
		sb.WriteString(fmt.Sprintf("%s\n", r.Header.string(1, "\t")))
	} else {
		sb.WriteString("\tHeader: nil\n")
	}

	if len(r.queries) > 0 {
		sb.WriteString("\tQuestions: [\n")
		for _, q := range r.queries {
			sb.WriteString(fmt.Sprintf("%s\n", q.string(2, "\t")))
		}
		sb.WriteString("\t]\n")
	} else {
		sb.WriteString("\tQuestions: []\n")
	}

	if len(r.answers) > 0 {
		sb.WriteString("\tAnswers: [\n")
		for _, a := range r.answers {
			sb.WriteString(fmt.Sprintf("%s\n", a.string(2, "\t")))
		}
		sb.WriteString("\t]\n")
	} else {
		sb.WriteString("\tAnswers: []\n")
	}
	return sb.String()
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
