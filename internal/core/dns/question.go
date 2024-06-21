package dns

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

type (
	question struct {
		qNames []questionName
		qType  RecordType  // Specifies the type of the query. Ex: CNAME, A, MX, NS
		qClass RecordClass // Specifies the class of the query.
	}

	questionName struct {
		bytesToRead uint8
		data        []byte
	}
)

func newQuestion(domain string, rType RecordType, rClass RecordClass) question {
	var qNames []questionName

	domainParts := strings.Split(domain, ".")

	for _, domainPart := range domainParts {
		questionName := newQuestionName(domainPart)
		qNames = append(qNames, questionName)
	}

	return question{
		qNames,
		rType,
		rClass,
	}
}

func newQuestionName(domainPart string) questionName {
	length := uint8(len(domainPart))

	return questionName{length, []byte(domainPart)}
}

func parseQuestion(bytes []byte) (question, int, error) {
	var q question
	be := binary.BigEndian

	qNames, bytesRead, err := parseQuestionNames(bytes)
	if err != nil {
		return q, 0, err
	}

	if len(bytes) < bytesRead+4 {
		return q, 0, errors.New("not enough bytes to parse the question")
	}

	q.qNames = qNames
	q.qType = RecordType(be.Uint16(bytes[bytesRead : bytesRead+2]))
	q.qClass = RecordClass(be.Uint16(bytes[bytesRead+2 : bytesRead+4]))

	return q, bytesRead + 4, nil
}

func parseQuestionNames(bytes []byte) ([]questionName, int, error) {
	var qNames []questionName
	totalBytesRead := 0

	for i := 0; i < len(bytes); {
		bytesToRead := bytes[i]
		if bytesToRead == 0 {
			totalBytesRead = i + 1
			break
		}
		if i+int(bytesToRead)+1 > len(bytes) {
			return nil, 0, errors.New("invalid question name length")
		}
		qName := questionName{
			bytesToRead: bytesToRead,
			data:        bytes[i+1 : i+1+int(bytesToRead)],
		}
		qNames = append(qNames, qName)
		i += int(bytesToRead) + 1
	}

	return qNames, totalBytesRead, nil
}

func (q question) string(indent int, char string) string {
	i := strings.Repeat(char, indent)

	var sb strings.Builder
	sb.WriteString("Question {\n")

	for _, qn := range q.qNames {
		sb.WriteString(fmt.Sprintf("%sName: %s\n", i, string(qn.data)))
	}

	sb.WriteString(fmt.Sprintf("%sType: %s\n", i, q.qType.TypeText()))
	sb.WriteString(fmt.Sprintf("%sClass: %s\n", i, q.qClass.ClassText()))
	sb.WriteString(fmt.Sprintf("%s}", i))

	return sb.String()
}

func (q question) toBytes() []byte {
	var (
		bytes []byte
		be    = binary.BigEndian
	)

	for _, qn := range q.qNames {
		bytes = append(bytes, qn.toBytes()...)
	}

	bytes = append(bytes, 0)

	typeAndClass := make([]byte, 4)
	be.PutUint16(typeAndClass[0:2], uint16(q.qType))
	be.PutUint16(typeAndClass[2:4], uint16(q.qClass))
	bytes = append(bytes, typeAndClass...)

	return bytes
}

func (qn questionName) toBytes() []byte {
	var bytes []byte

	bytes = append(bytes, qn.bytesToRead)
	bytes = append(bytes, qn.data...)

	return bytes
}
