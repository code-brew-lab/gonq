package dns

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

type (
	Header struct { // 96 bits total
		id                    ID     // [Req + Resp]  Unique request id. Same for the following response.
		flags                 *Flags // [Req + Resp]  See flags.go
		questionCount         uint16 // [Req]         Number of entries inside the question.
		answerCount           uint16 // [Resp]        Number of response entries from DNS server.
		nameServerCount       uint16 // [Req]         Number of name server resource records in the authority records section.
		additionalRecordCount uint16 // [Req]         Number of resource records in the additional records section.
	}

	HeaderBuilder struct {
		*Header
		errors []error
	}
)

const headerSize int = 12

// NewHeaderBuilder creates a new HeaderBuilder instance with default values
func NewHeaderBuilder() *HeaderBuilder {
	return &HeaderBuilder{
		Header: newHeader(),
	}
}

func parseHeader(bytes []byte) (*Header, error) {
	if len(bytes) < headerSize {
		return nil, errors.New("header should be 12 bytes")
	}

	header := new(Header)
	be := binary.BigEndian

	id := bytes[:2]
	flags := bytes[2:4]

	parsedID, err := ParseID(id)
	if err != nil {
		return nil, err
	}

	parsedFlags, err := parseFlags(flags)
	if err != nil {
		return nil, err
	}

	header.id = parsedID
	header.flags = parsedFlags
	header.questionCount = be.Uint16(bytes[4:6])
	header.answerCount = be.Uint16(bytes[6:8])
	header.nameServerCount = be.Uint16(bytes[8:10])
	header.additionalRecordCount = be.Uint16(bytes[10:12])

	return header, nil
}

func newHeader() *Header {
	return &Header{}
}

func (bh *HeaderBuilder) SetID(id ID) *HeaderBuilder {
	bh.id = id
	return bh
}

func (bh *HeaderBuilder) SetFlags(flags *Flags) *HeaderBuilder {
	if flags == nil {
		return bh
	}

	bh.flags = flags
	return bh
}

func (bh *HeaderBuilder) SetQuestionCount(count uint16) *HeaderBuilder {
	bh.questionCount = count
	return bh
}

func (bh *HeaderBuilder) SetAnswerCount(count uint16) *HeaderBuilder {
	bh.answerCount = count
	return bh
}

func (bh *HeaderBuilder) SetNameServerCount(count uint16) *HeaderBuilder {
	bh.nameServerCount = count
	return bh
}

func (bh *HeaderBuilder) SetAdditionalRecordCount(count uint16) *HeaderBuilder {
	bh.additionalRecordCount = count
	return bh
}

func (bh *HeaderBuilder) AddError(err error) {
	bh.errors = append(bh.errors, err)
}

func (bh *HeaderBuilder) Build() *Header {
	return bh.Header
}

func (h *Header) IsTruncated() bool {
	return h.flags.isTruncated
}

func (h *Header) IsAuthoritative() bool {
	return h.flags.isAuthoritative
}

func (h *Header) CanRecursive() bool {
	return h.flags.canRecursive
}

func (h *Header) ResponseCode() ResponseCode {
	return h.flags.responseCode
}

func (h *Header) QuestionCount() uint16 {
	return h.questionCount
}

func (h *Header) AnswerCount() uint16 {
	return h.answerCount
}

func (h *Header) addQuestion() {
	h.questionCount += 1
}

func (h *Header) string(indent int, char string) string {
	i := strings.Repeat(char, indent)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%sHeader:\n", i))
	sb.WriteString(fmt.Sprintf("%sID: %v\n", i, h.id))
	sb.WriteString(fmt.Sprintf("%s%v\n", i, h.flags.string(indent+1, char)))
	sb.WriteString(fmt.Sprintf("%sQuestion Count: %v\n", i, h.questionCount))
	sb.WriteString(fmt.Sprintf("%sAnswer Count: %v\n", i, h.answerCount))
	sb.WriteString(fmt.Sprintf("%sName Server Count: %v\n", i, h.nameServerCount))
	sb.WriteString(fmt.Sprintf("%sAdditional Record Count: %v", i, h.additionalRecordCount))
	return sb.String()
}

func (h *Header) toBytes() []byte {
	var bytes [12]byte
	be := binary.BigEndian

	be.PutUint16(bytes[0:2], h.id.toUint16())
	be.PutUint16(bytes[2:4], h.flags.toUint16())
	be.PutUint16(bytes[4:6], h.questionCount)
	be.PutUint16(bytes[6:8], h.answerCount)
	be.PutUint16(bytes[8:10], h.nameServerCount)
	be.PutUint16(bytes[10:12], h.additionalRecordCount)

	return bytes[:]
}
