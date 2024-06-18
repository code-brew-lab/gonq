package dns

import (
	"fmt"

	"github.com/code-brew-lab/gonq.git/internal/pkg/bitwise"
)

type (
	Header struct { // 96 bits total
		ID                    uint16 // [Req + Resp] 16 bit unique request id. Same for the following response.
		Flags                 uint16
		QuestionCount         uint16 // [Req] Number of entries inside the question.
		AnswerCount           uint16 // [Resp] Number of response entries from DNS server.
		NameServerCount       uint16 // [Req] Number of name server resource records in the authority records section.
		AdditionalRecordCount uint16 // [Req] Number of resource records in the additional records section.
	}

	Flags struct { // 16 bits total.
		IsQuery         bitwise.Bit    // [Req + Resp] Is the message query or response. Query = False; Response = True;
		OperationCode   [4]bitwise.Bit // [Req] A four bit field that specifies kind of query in this message.
		IsAuthoritative bitwise.Bit    // [Resp] Is response from authoritative dns server. Authoritative = True;
		IsTruncated     bitwise.Bit    // [Resp] Is response truncated or not.
		IsRecursive     bitwise.Bit    // [Req] Ask DNS server to recursively ask for the domain.
		CanRecursive    bitwise.Bit    // [Resp] Shows if recursion available for DNS server.
		FutureUse       [3]bitwise.Bit // [Req] A three bits future use field.
		ResponseCode    [4]bitwise.Bit // [Resp] A four bits response codes.
	}
)

func (f Flags) BinaryMarshaler() ([]byte, error) {
	byteSet := bitwise.NewSet()

	byteSet.Add(f.IsQuery)
	byteSet.AddSet(f.OperationCode[:])
	byteSet.Add(f.IsAuthoritative)
	byteSet.Add(f.IsTruncated)
	byteSet.Add(f.IsRecursive)
	byteSet.Add(f.CanRecursive)
	byteSet.AddSet(f.FutureUse[:])
	byteSet.AddSet(f.ResponseCode[:])

	bytes, err := byteSet.ToBytes()
	if err != nil {
		return nil, fmt.Errorf("dns: %v", err)
	}

	return bytes, nil
}
