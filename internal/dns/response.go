package dns

import (
	"errors"
	"fmt"
	"strings"
)

type Response struct {
	Header  *Header
	queries []query
	answers []answer
}

func ParseResponse(response []byte) (*Response, error) {
	if len(response) < 12 {
		return nil, errors.New("dns.ParseResponse: response too short to contain a valid header")
	}

	header, err := parseHeader(response[:12])
	if err != nil {
		return nil, fmt.Errorf("dns.ParseResponse: %v", err)
	}

	if header.IsTruncated() {
		return nil, errors.New("dns.ParseResponse: response is truncated")
	}

	if header.ResponseCode() != CodeNoError {
		return nil, fmt.Errorf("dns.ParseResponse: request responded with status: %s", header.ResponseCode().CodeText())
	}

	response = response[12:]
	qCount := header.QuestionCount()
	var questions []query

	for i := 0; i < int(qCount); i++ {
		q, bytesRead, err := parseQuery(response)
		if err != nil {
			return nil, fmt.Errorf("dns.ParseResponse: %v", err)
		}

		questions = append(questions, q)
		if bytesRead > len(response) {
			return nil, errors.New("dns.ParseResponse: byte slice length exceeded")
		}
		response = response[bytesRead:]
	}

	aCount := header.AnswerCount()
	var answers []answer

	for i := 0; i < int(aCount); i++ {
		a, bytesRead, err := parseAnswer(response)
		if err != nil {
			return nil, fmt.Errorf("dns.ParseResponse: %v", err)
		}

		answers = append(answers, a)
		if bytesRead > len(response) {
			return nil, errors.New("dns.ParseResponse: byte slice length exceeded")
		}
		response = response[bytesRead:]
	}

	return &Response{
		Header:  header,
		queries: questions,
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
