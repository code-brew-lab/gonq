package dns

import (
	"encoding/binary"
	"errors"
)

type (
	header struct { // 96 bits total
		id                    ID     // [Req + Resp]  Unique request id. Same for the following response.
		flags                 *flags // [Req + Resp]  See flags.go
		queryCount            uint16 // [Req]         Number of entries inside the query.
		answerCount           uint16 // [Resp]        Number of response entries from DNS server.
		nameServerCount       uint16 // [Req]         Number of name server resource records in the authority records section.
		additionalRecordCount uint16 // [Req]         Number of resource records in the additional records section.
	}

	headerBuilder struct {
		*header
		errors []error
	}
)

const headerSize int = 12

func newHeaderBuilder() *headerBuilder {
	return &headerBuilder{
		header: newHeader(),
	}
}

func parseHeader(bytes []byte) (*header, error) {
	if len(bytes) < headerSize {
		return nil, errors.New("header should be 12 bytes")
	}

	header := new(header)
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
	header.queryCount = be.Uint16(bytes[4:6])
	header.answerCount = be.Uint16(bytes[6:8])
	header.nameServerCount = be.Uint16(bytes[8:10])
	header.additionalRecordCount = be.Uint16(bytes[10:12])

	return header, nil
}

func newHeader() *header {
	return &header{}
}

func (bh *headerBuilder) SetID(id ID) *headerBuilder {
	bh.id = id
	return bh
}

func (bh *headerBuilder) SetFlags(flags *flags) *headerBuilder {
	if flags == nil {
		return bh
	}

	bh.flags = flags
	return bh
}

func (bh *headerBuilder) SetQuestionCount(count uint16) *headerBuilder {
	bh.queryCount = count
	return bh
}

func (bh *headerBuilder) SetAnswerCount(count uint16) *headerBuilder {
	bh.answerCount = count
	return bh
}

func (bh *headerBuilder) SetNameServerCount(count uint16) *headerBuilder {
	bh.nameServerCount = count
	return bh
}

func (bh *headerBuilder) SetAdditionalRecordCount(count uint16) *headerBuilder {
	bh.additionalRecordCount = count
	return bh
}

func (bh *headerBuilder) AddError(err error) {
	bh.errors = append(bh.errors, err)
}

func (bh *headerBuilder) Build() *header {
	return bh.header
}

func (h *header) ID() ID {
	return h.id
}

func (h *header) IsTruncated() bool {
	return h.flags.isTruncated
}

func (h *header) IsAuthoritative() bool {
	return h.flags.isAuthoritative
}

func (h *header) CanRecursive() bool {
	return h.flags.canRecursive
}

func (h *header) ResponseCode() ResponseCode {
	return h.flags.responseCode
}

func (h *header) QueryCount() uint16 {
	return h.queryCount
}

func (h *header) AnswerCount() uint16 {
	return h.answerCount
}

func (h *header) addQuestion() {
	h.queryCount += 1
}

func (h *header) toBytes() []byte {
	var bytes [12]byte
	be := binary.BigEndian

	be.PutUint16(bytes[0:2], h.id.toUint16())
	be.PutUint16(bytes[2:4], h.flags.toUint16())
	be.PutUint16(bytes[4:6], h.queryCount)
	be.PutUint16(bytes[6:8], h.answerCount)
	be.PutUint16(bytes[8:10], h.nameServerCount)
	be.PutUint16(bytes[10:12], h.additionalRecordCount)

	return bytes[:]
}
