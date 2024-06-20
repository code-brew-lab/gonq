package dns

import (
	"errors"
)

type (
	Request struct {
		header    *Header
		questions []question
	}
)

func NewRequest() *Request {
	flags := NewFlagsBuilder().SetIsQuery(true).SetIsRecursive(true).Build()
	header := NewHeaderBuilder().SetID(NewID()).SetFlags(flags).Build()

	return &Request{
		header: header,
	}
}

func NewRequestWithHeader(header *Header) (*Request, error) {
	if header == nil {
		return nil, errors.New("dns: header is nil")
	}

	return &Request{
		header:    header,
		questions: make([]question, 1),
	}, nil
}

func (r *Request) AddQuestion(domain string, qType QuestionType, qClass QuestionClass) {
	question := newQuestion(domain, qType, qClass)
	r.questions = append(r.questions, question)
	r.header.IncrementQuestionCount()
}

func (r *Request) BinaryMarshaler() ([]byte, error) {
	var bytes []byte

	header := r.header
	questions := r.questions

	bytes = append(bytes, header.toBytes()...)
	for _, q := range questions {
		bytes = append(bytes, q.toBytes()...)
	}

	return bytes, nil
}
