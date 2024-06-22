package dns

import (
	"errors"
)

type (
	flags struct { // 16 bits total.
		isQuery         bool         // [Req + Resp]  Is the message query or response.
		queryKind       QueryKind    // [Req]         A four-bit field that specifies the kind of query. See querykind.go.
		isAuthoritative bool         // [Resp]        Is the response from authoritative DNS server.
		isTruncated     bool         // [Resp]        Is the response truncated or not.
		isRecursive     bool         // [Req]         Ask DNS server to recursively ask for the domain.
		canRecursive    bool         // [Resp]        Shows if recursion is available for DNS server.
		futureUse       uint8        // [Req]         A three-bit future use field.
		responseCode    ResponseCode // [Resp]        A four-bit response code. See responsecode.go.
	}

	flagsBuilder struct {
		*flags
		errors []error
	}
)

const flagsSize int = 2

func newFlagsBuilder() *flagsBuilder {
	return &flagsBuilder{
		flags: newFlags(),
	}

}

func parseFlags(bytes []byte) (*flags, error) {
	if len(bytes) < flagsSize {
		return nil, errors.New("flags should be 2 bytes")
	}

	byte1 := bytes[0]
	byte2 := bytes[1]
	flags := new(flags)

	isQuery := (byte1 & 0x80) >> 7
	operationCode := (byte1 & 0x78) >> 3
	isAuthoritative := (byte1 & 0x04) >> 2
	isTruncated := (byte1 & 0x02) >> 1
	isRecursive := byte1 & 0x01

	canRecursive := (byte2 & 0x80) >> 7
	futureUse := (byte2 & 0x70) >> 4
	responseCode := byte2 & 0x0F

	flags.isQuery = isQuery == 0
	flags.queryKind = QueryKind(operationCode)
	flags.isAuthoritative = isAuthoritative == 1
	flags.isTruncated = isTruncated == 1
	flags.isRecursive = isRecursive == 1
	flags.canRecursive = canRecursive == 1
	flags.futureUse = uint8(futureUse)
	flags.responseCode = ResponseCode(responseCode)

	return flags, nil
}

func newFlags() *flags {
	return &flags{
		isQuery:      true,
		isRecursive:  true,
		queryKind:    KindStandard,
		responseCode: CodeNoError,
	}
}

func (fb *flagsBuilder) SetIsQuery(isQuery bool) *flagsBuilder {
	fb.isQuery = isQuery
	return fb
}

func (fb *flagsBuilder) SetOperationCode(operationCode QueryKind) *flagsBuilder {
	fb.queryKind = operationCode
	return fb
}

func (fb *flagsBuilder) SetIsAuthoritative(isAuthoritative bool) *flagsBuilder {
	fb.isAuthoritative = isAuthoritative
	return fb
}

func (fb *flagsBuilder) SetIsTruncated(isTruncated bool) *flagsBuilder {
	fb.isTruncated = isTruncated
	return fb
}

func (fb *flagsBuilder) SetIsRecursive(isRecursive bool) *flagsBuilder {
	fb.isRecursive = isRecursive
	return fb
}

func (fb *flagsBuilder) SetCanRecursive(canRecursive bool) *flagsBuilder {
	fb.canRecursive = canRecursive
	return fb
}

func (fb *flagsBuilder) SetFutureUse(futureUse uint8) *flagsBuilder {
	fb.futureUse = futureUse
	return fb
}

func (fb *flagsBuilder) SetResponseCode(responseCode ResponseCode) *flagsBuilder {
	fb.responseCode = responseCode
	return fb
}

func (fb *flagsBuilder) AddError(err error) {
	fb.errors = append(fb.errors, err)
}

func (fb *flagsBuilder) Build() *flags {
	return fb.flags
}

func (f *flags) toUint16() uint16 {
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
