package dns

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

type query struct {
	queryData   []queryData
	recordType  RecordType  // Specifies the type of the query. Ex: CNAME, A, MX, NS
	recordClass RecordClass // Specifies the class of the query.
}

func newQuery(domain string, rType RecordType, rClass RecordClass) query {
	var qNames []queryData

	domainParts := strings.Split(domain, ".")

	for _, domainPart := range domainParts {
		questionName := newQueryData(domainPart)
		qNames = append(qNames, questionName)
	}

	return query{
		qNames,
		rType,
		rClass,
	}
}

func parseQuery(bytes []byte) (query, int, error) {
	var q query
	be := binary.BigEndian

	qNames, bytesRead, err := parseQueryData(bytes)
	if err != nil {
		return q, 0, err
	}

	if len(bytes) < bytesRead+4 {
		return q, 0, errors.New("not enough bytes to parse the question")
	}

	q.queryData = qNames
	q.recordType = RecordType(be.Uint16(bytes[bytesRead : bytesRead+2]))
	q.recordClass = RecordClass(be.Uint16(bytes[bytesRead+2 : bytesRead+4]))

	return q, bytesRead + 4, nil
}

func (q query) toBytes() []byte {
	var (
		bytes []byte
		be    = binary.BigEndian
	)

	for _, qn := range q.queryData {
		bytes = append(bytes, qn.toBytes()...)
	}

	bytes = append(bytes, 0)

	typeAndClass := make([]byte, 4)
	be.PutUint16(typeAndClass[0:2], uint16(q.recordType))
	be.PutUint16(typeAndClass[2:4], uint16(q.recordClass))
	bytes = append(bytes, typeAndClass...)

	return bytes
}

func (q query) string(indent int, char string) string {
	i := strings.Repeat(char, indent)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%sQuestion:\n", i))
	sb.WriteString(fmt.Sprintf("%sName: %s\n", i, q.domainToString()))
	sb.WriteString(fmt.Sprintf("%sType: %s\n", i, q.recordType.TypeText()))
	sb.WriteString(fmt.Sprintf("%sClass: %s\n", i, q.recordClass.ClassText()))

	return sb.String()
}

func (q query) domainToString() string {
	var sb strings.Builder

	for i, qn := range q.queryData {
		sb.WriteString(qn.string())
		if i == len(q.queryData)-1 {
			break
		}

		sb.WriteString(".")
	}

	return sb.String()
}
