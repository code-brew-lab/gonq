package dns

import (
	"encoding/binary"
	"strings"
)

type (
	Question struct {
		questionNames []QuestionName
		questionType  QuestionType  // Specifies the type of the query. Ex: CNAME, A, MX, NS
		questionClass QuestionClass // Specifies the class of the query.
	}

	QuestionName struct {
		bytesToRead uint8
		data        []byte
	}

	QuestionType  uint16
	QuestionClass uint16
)

const (
	AType  QuestionType = 1  // A record type.
	NSType QuestionType = 2  // Mail server record type.
	MXType QuestionType = 15 // Name server record type.

	INetClass QuestionClass = 1 // Internet Address class.
)

func newQuestion(domain string, qType QuestionType, qClass QuestionClass) Question {
	var questionNames []QuestionName

	domainParts := strings.Split(domain, ".")

	for _, domainPart := range domainParts {
		questionName := newQuestionName(domainPart)
		questionNames = append(questionNames, questionName)
	}

	return Question{
		questionNames,
		qType,
		qClass,
	}
}

func newQuestionName(domainPart string) QuestionName {
	length := uint8(len(domainPart))

	return QuestionName{length, []byte(domainPart)}
}

func (q Question) toBytes() []byte {
	var (
		bytes []byte
		be    = binary.BigEndian
		skip  = 0
	)

	for _, qn := range q.questionNames {
		bytes = append(bytes, qn.toBytes()...)
		skip += (1 + len(qn.data))
	}

	bytes = append(bytes, 0)
	skip += 1

	typeAndClass := make([]byte, 4)
	be.PutUint16(typeAndClass[0:2], uint16(q.questionType))
	be.PutUint16(typeAndClass[2:4], uint16(q.questionClass))
	bytes = append(bytes, typeAndClass...)

	return bytes
}

func (qn QuestionName) toBytes() []byte {
	var bytes []byte

	bytes = append(bytes, qn.bytesToRead)
	bytes = append(bytes, qn.data...)

	return bytes
}
