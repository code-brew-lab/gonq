package dns

import (
	"encoding/binary"
)

type (
	Header struct { // 96 bits total
		id                    uint16 // [Req + Resp]  Unique request id. Same for the following response.
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

// NewHeaderBuilder creates a new HeaderBuilder instance with default values
func NewHeaderBuilder() *HeaderBuilder {
	return &HeaderBuilder{
		Header: newHeader(),
	}
}

func newHeader() *Header {
	return &Header{}
}

func (bh *HeaderBuilder) SetID(id uint16) *HeaderBuilder {
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

func (h *Header) IncrementQuestionCount() {
	h.questionCount += 1
}

func (h *Header) toBytes() []byte {
	var bytes [12]byte
	be := binary.BigEndian

	be.PutUint16(bytes[0:2], h.id)
	be.PutUint16(bytes[2:4], h.flags.toUint16())
	be.PutUint16(bytes[4:6], h.questionCount)
	be.PutUint16(bytes[6:8], h.answerCount)
	be.PutUint16(bytes[8:10], h.nameServerCount)
	be.PutUint16(bytes[10:12], h.additionalRecordCount)

	return bytes[:]
}
