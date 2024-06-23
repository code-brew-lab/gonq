package dns

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type (
	Header struct { // 96 bits total
		id                    ID     // [Req + Resp]  Unique request id. Same for the following response.
		flags                 *Flags // [Req + Resp]  See flags.go
		queryCount            uint16 // [Req]         Number of entries inside the query.
		answerCount           uint16 // [Resp]        Number of response entries from DNS server.
		nameServerCount       uint16 // [Req]         Number of name server resource records in the authority records section.
		additionalRecordCount uint16 // [Req]         Number of resource records in the additional records section.
	}

	HeaderBuilder struct {
		*Header
	}
)

const HeaderSize int = 12

func NewHeaderBuilder() *HeaderBuilder {
	return &HeaderBuilder{
		Header: newHeader(),
	}
}

func ParseHeader(bytes []byte) (*Header, error) {
	if len(bytes) < HeaderSize {
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
	header.queryCount = be.Uint16(bytes[4:6])
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

func (bh *HeaderBuilder) Build() *Header {
	return bh.Header
}

func (h *Header) ID() ID {
	return h.id
}

func (h *Header) Flags() *Flags {
	return h.flags
}

func (h *Header) QueryCount() uint16 {
	return h.queryCount
}

func (h *Header) AnswerCount() uint16 {
	return h.answerCount
}

func (h *Header) NameServerCount() uint16 {
	return h.nameServerCount
}

func (h *Header) AdditionalRecordCount() uint16 {
	return h.additionalRecordCount
}

func (h *Header) AddQuestion() {
	h.queryCount += 1
}

func (h *Header) ToBytes() []byte {
	var bytes [12]byte
	be := binary.BigEndian

	be.PutUint16(bytes[0:2], h.id.ToUint16())
	be.PutUint16(bytes[2:4], h.flags.toUint16())
	be.PutUint16(bytes[4:6], h.queryCount)
	be.PutUint16(bytes[6:8], h.answerCount)
	be.PutUint16(bytes[8:10], h.nameServerCount)
	be.PutUint16(bytes[10:12], h.additionalRecordCount)

	return bytes[:]
}

func (h *Header) String() string {
	return fmt.Sprintf("ID: %v\n\tFlags: %v\n\tQDCount: %d\n\tANCount: %d\n\tNSCount: %d\n\tARCount: %d",
		h.id, h.flags, h.queryCount, h.answerCount, h.nameServerCount, h.additionalRecordCount)
}
