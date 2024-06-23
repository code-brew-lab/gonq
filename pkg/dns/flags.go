package dns

import (
	"errors"
	"fmt"
)

type (
	Flags struct { // 16 bits total.
		isQuery         bool         // [Req + Resp]  Is the message query or response.
		queryKind       queryKind    // [Req]         A four-bit field that specifies the kind of query. See querykind.go.
		isAuthoritative bool         // [Resp]        Is the response from authoritative DNS server.
		isTruncated     bool         // [Resp]        Is the response truncated or not.
		isRecursive     bool         // [Req]         Ask DNS server to recursively ask for the domain.
		canRecursive    bool         // [Resp]        Shows if recursion is available for DNS server.
		futureUse       uint8        // [Req]         A three-bit future use field.
		responseCode    ResponseCode // [Resp]        A four-bit response code. See responsecode.go.
	}

	FlagsBuilder struct {
		*Flags
	}
)

const FlagsSize int = 2

func NewFlagsBuilder() *FlagsBuilder {
	return &FlagsBuilder{
		Flags: newFlags(),
	}
}

func parseFlags(bytes []byte) (*Flags, error) {
	if len(bytes) < FlagsSize {
		return nil, errors.New("flags should be 2 bytes")
	}

	byte1 := bytes[0]
	byte2 := bytes[1]
	flags := new(Flags)

	isQuery := (byte1 & 0x80) >> 7
	operationCode := (byte1 & 0x78) >> 3
	isAuthoritative := (byte1 & 0x04) >> 2
	isTruncated := (byte1 & 0x02) >> 1
	isRecursive := byte1 & 0x01

	canRecursive := (byte2 & 0x80) >> 7
	futureUse := (byte2 & 0x70) >> 4
	responseCode := byte2 & 0x0F

	flags.isQuery = isQuery == 0
	flags.queryKind = queryKind(operationCode)
	flags.isAuthoritative = isAuthoritative == 1
	flags.isTruncated = isTruncated == 1
	flags.isRecursive = isRecursive == 1
	flags.canRecursive = canRecursive == 1
	flags.futureUse = uint8(futureUse)
	flags.responseCode = ResponseCode(responseCode)

	return flags, nil
}

func newFlags() *Flags {
	return &Flags{
		isQuery:      true,
		isRecursive:  true,
		queryKind:    kindStandard,
		responseCode: CodeNoError,
	}
}

func (fb *FlagsBuilder) SetIsQuery(isQuery bool) *FlagsBuilder {
	fb.isQuery = isQuery
	return fb
}

func (fb *FlagsBuilder) SetQueryKind(queryKind queryKind) *FlagsBuilder {
	fb.queryKind = queryKind
	return fb
}

func (fb *FlagsBuilder) SetIsRecursive(isRecursive bool) *FlagsBuilder {
	fb.isRecursive = isRecursive
	return fb
}

func (fb *FlagsBuilder) SetFutureUse(futureUse uint8) *FlagsBuilder {
	fb.futureUse = futureUse
	return fb
}

func (fb *FlagsBuilder) Build() *Flags {
	return fb.Flags
}

func (f *Flags) IsQuery() bool {
	return f.isQuery
}

func (f *Flags) QueryKind() queryKind {
	return f.queryKind
}

func (f *Flags) IsAuthoritative() bool {
	return f.isAuthoritative
}

func (f *Flags) IsTruncated() bool {
	return f.isTruncated
}

func (f *Flags) IsRecursive() bool {
	return f.isRecursive
}

func (f *Flags) CanRecursive() bool {
	return f.canRecursive
}

func (f *Flags) FutureUse() uint8 {
	return f.futureUse
}

func (f *Flags) ResponseCode() ResponseCode {
	return f.responseCode
}

func (f *Flags) String() string {
	return fmt.Sprintf("IsQuery: %v, QueryKind: %v, IsAuthoritative: %v, IsTruncated: %v, IsRecursive: %v, CanRecursive: %v, ResponseCode: %v",
		f.isQuery, f.queryKind.kindText(), f.isAuthoritative, f.isTruncated, f.isRecursive, f.canRecursive, f.responseCode.CodeText())
}

func (f *Flags) toUint16() uint16 {
	var result uint16

	if !f.isQuery {
		result |= 1 << 15
	}

	result |= uint16(f.queryKind&0x0F) << 11

	if f.isAuthoritative {
		result |= 1 << 10
	}

	if f.isTruncated {
		result |= 1 << 9
	}

	if f.isRecursive {
		result |= 1 << 8
	}

	if f.canRecursive {
		result |= 1 << 7
	}

	result |= (uint16(f.futureUse) & 0x07) << 4
	result |= uint16(f.responseCode) & 0x0F

	return result
}
