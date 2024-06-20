package dns

import (
	"encoding/binary"
	"strings"
)

type (
	question struct {
		questionNames []questionName
		questionType  RecordType  // Specifies the type of the query. Ex: CNAME, A, MX, NS
		questionClass RecordClass // Specifies the class of the query.
	}

	questionName struct {
		bytesToRead uint8
		data        []byte
	}

	RecordType  uint16
	RecordClass uint16
)

const (
	A     RecordType = 1  // A record type.
	NS    RecordType = 2  // Mail server record type.
	CNAME RecordType = 5  // Canonical name of the domain.
	MX    RecordType = 15 // Name server record type.

	INET RecordClass = 1 // Internet Address class.
)

func newQuestion(domain string, rType RecordType, rClass RecordClass) question {
	var questionNames []questionName

	domainParts := strings.Split(domain, ".")

	for _, domainPart := range domainParts {
		questionName := newQuestionName(domainPart)
		questionNames = append(questionNames, questionName)
	}

	return question{
		questionNames,
		rType,
		rClass,
	}
}

func newQuestionName(domainPart string) questionName {
	length := uint8(len(domainPart))

	return questionName{length, []byte(domainPart)}
}

func (q question) toBytes() []byte {
	var (
		bytes []byte
		be    = binary.BigEndian
	)

	for _, qn := range q.questionNames {
		bytes = append(bytes, qn.toBytes()...)
	}

	bytes = append(bytes, 0)

	typeAndClass := make([]byte, 4)
	be.PutUint16(typeAndClass[0:2], uint16(q.questionType))
	be.PutUint16(typeAndClass[2:4], uint16(q.questionClass))
	bytes = append(bytes, typeAndClass...)

	return bytes
}

func (qn questionName) toBytes() []byte {
	var bytes []byte

	bytes = append(bytes, qn.bytesToRead)
	bytes = append(bytes, qn.data...)

	return bytes
}
