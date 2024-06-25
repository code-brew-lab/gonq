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

const queryHeaderSize int = 4

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

	if len(bytes) < bytesRead+queryHeaderSize {
		return q, 0, errors.New("not enough bytes to parse the question")
	}

	q.queryData = qNames
	q.recordType = RecordType(be.Uint16(bytes[bytesRead : bytesRead+2]))
	q.recordClass = RecordClass(be.Uint16(bytes[bytesRead+2 : bytesRead+4]))

	return q, bytesRead + queryHeaderSize, nil
}

func (q query) domain() string {
	var domainParts []string
	for _, qData := range q.queryData {
		domainParts = append(domainParts, string(qData.data))
	}

	return strings.Join(domainParts, ".")
}

func (q query) string() string {
	domain := q.domain()
	return fmt.Sprintf("Domain: %s, RecordType: %s, RecordClass: %s", domain, q.recordType.TypeText(), q.recordClass.ClassText())
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

	typeAndClass := make([]byte, queryHeaderSize)
	be.PutUint16(typeAndClass[0:2], uint16(q.recordType))
	be.PutUint16(typeAndClass[2:4], uint16(q.recordClass))
	bytes = append(bytes, typeAndClass...)

	return bytes
}
