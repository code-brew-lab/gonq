package dns

import (
	"encoding/binary"
	"fmt"
)

type (
	Header struct { // 96 bits total
		ID                    uint16 // [Req + Resp] 16 bit unique request id. Same for the following response.
		Flags                 Flags
		QuestionCount         uint16 // [Req] Number of entries inside the question.
		AnswerCount           uint16 // [Resp] Number of response entries from DNS server.
		NameServerCount       uint16 // [Req] Number of name server resource records in the authority records section.
		AdditionalRecordCount uint16 // [Req] Number of resource records in the additional records section.
	}

	Flags struct { // 16 bits total.
		IsQuery         bool  // [Req + Resp] Is the message query or response. Query = False; Response = True;
		OperationCode   uint8 // [Req] A four-bit field that specifies the kind of query in this message.
		IsAuthoritative bool  // [Resp] Is the response from authoritative DNS server. Authoritative = True;
		IsTruncated     bool  // [Resp] Is the response truncated or not.
		IsRecursive     bool  // [Req] Ask DNS server to recursively ask for the domain.
		CanRecursive    bool  // [Resp] Shows if recursion is available for DNS server.
		FutureUse       uint8 // [Req] A three-bit future use field.
		ResponseCode    uint8 // [Resp] A four-bit response code.
	}
)

// NewHeader creates a new Header instance with default values
func NewHeader(flags Flags, questionCount int) *Header {
	return &Header{
		ID:                    0,
		Flags:                 flags,
		QuestionCount:         uint16(questionCount),
		AnswerCount:           0,
		NameServerCount:       0,
		AdditionalRecordCount: 0,
	}
}

// NewFlags creates a new Flags instance with default values
func NewFlags(isQuery, isRecursive bool) Flags {
	return Flags{
		IsQuery:         isQuery,
		OperationCode:   0,
		IsAuthoritative: false,
		IsTruncated:     false,
		IsRecursive:     isRecursive,
		CanRecursive:    false,
		FutureUse:       0,
		ResponseCode:    0,
	}
}

func (h Header) BinaryMarshaler() ([]byte, error) {
	var bytes [12]byte
	flags, err := h.Flags.BinaryMarshaler()
	if err != nil {
		return nil, fmt.Errorf("dns: %v", err)
	}

	binary.BigEndian.PutUint16(bytes[0:2], h.ID)
	copy(bytes[2:4], flags)
	binary.BigEndian.PutUint16(bytes[4:6], h.QuestionCount)
	binary.BigEndian.PutUint16(bytes[6:8], h.AnswerCount)
	binary.BigEndian.PutUint16(bytes[8:10], h.NameServerCount)
	binary.BigEndian.PutUint16(bytes[10:12], h.AdditionalRecordCount)

	return bytes[:], nil
}

func (f Flags) BinaryMarshaler() ([]byte, error) {
	var flags uint16

	if f.IsQuery {
		flags |= 1 << 15
	}
	flags |= uint16(f.OperationCode&0xF) << 11
	if f.IsAuthoritative {
		flags |= 1 << 10
	}
	if f.IsTruncated {
		flags |= 1 << 9
	}
	if f.IsRecursive {
		flags |= 1 << 8
	}
	if f.CanRecursive {
		flags |= 1 << 7
	}
	flags |= uint16(f.FutureUse&0x7) << 4
	flags |= uint16(f.ResponseCode & 0xF)

	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, flags)

	return bytes, nil
}
