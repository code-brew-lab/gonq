package dns

import (
	"errors"
	"fmt"
	"strings"
)

type (
	Flags struct { // 16 bits total.
		isQuery         bool          // [Req + Resp]  Is the message query or response.
		operationCode   OperationCode // [Req]         A four-bit field that specifies the kind of query.
		isAuthoritative bool          // [Resp]        Is the response from authoritative DNS server.
		isTruncated     bool          // [Resp]        Is the response truncated or not.
		isRecursive     bool          // [Req]         Ask DNS server to recursively ask for the domain.
		canRecursive    bool          // [Resp]        Shows if recursion is available for DNS server.
		futureUse       uint8         // [Req]         A three-bit future use field.
		responseCode    ResponseCode  // [Resp]        A four-bit response code.
	}

	FlagsBuilder struct {
		*Flags
		errors []error
	}

	OperationCode uint8
	ResponseCode  uint8
)

const (
	StandardQuery OperationCode = iota
)

const (
	NoError        ResponseCode = iota // #0  OK.
	FormatError                        // #1  The name server was unable to interpret the query.
	ServerFailure                      // #2  Server unreachable.
	NameError                          // #3  Domain name does not exists.
	NotImplemented                     // #4  Unsupported request query type.
	Refused                            // #5  Refused for policy reasons of the server.
)

// NewFlagsBuilder creates a new FlagsBuilder instance with default values
func NewFlagsBuilder() *FlagsBuilder {
	return &FlagsBuilder{
		Flags: newFlags(),
	}
}

func parseFlags(bytes []byte) (*Flags, error) {
	if len(bytes) != 2 {
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
	flags.operationCode = OperationCode(operationCode)
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
		isQuery:       true,
		isRecursive:   true,
		operationCode: StandardQuery,
		responseCode:  NoError,
	}
}

func (fb *FlagsBuilder) SetIsQuery(isQuery bool) *FlagsBuilder {
	fb.isQuery = isQuery
	return fb
}

func (fb *FlagsBuilder) SetOperationCode(operationCode OperationCode) *FlagsBuilder {
	fb.operationCode = operationCode
	return fb
}

func (fb *FlagsBuilder) SetIsAuthoritative(isAuthoritative bool) *FlagsBuilder {
	fb.isAuthoritative = isAuthoritative
	return fb
}

func (fb *FlagsBuilder) SetIsTruncated(isTruncated bool) *FlagsBuilder {
	fb.isTruncated = isTruncated
	return fb
}

func (fb *FlagsBuilder) SetIsRecursive(isRecursive bool) *FlagsBuilder {
	fb.isRecursive = isRecursive
	return fb
}

func (fb *FlagsBuilder) SetCanRecursive(canRecursive bool) *FlagsBuilder {
	fb.canRecursive = canRecursive
	return fb
}

func (fb *FlagsBuilder) SetFutureUse(futureUse uint8) *FlagsBuilder {
	fb.futureUse = futureUse
	return fb
}

func (fb *FlagsBuilder) SetResponseCode(responseCode ResponseCode) *FlagsBuilder {
	fb.responseCode = responseCode
	return fb
}

func (fb *FlagsBuilder) AddError(err error) {
	fb.errors = append(fb.errors, err)
}

func (fb *FlagsBuilder) Build() *Flags {
	return fb.Flags
}

func (f *Flags) string(indent int, char string) string {
	i := strings.Repeat(char, indent)

	var sb strings.Builder
	sb.WriteString("Flags {\n")
	sb.WriteString(fmt.Sprintf("%sisQuery: %v\n", i, f.isQuery))
	sb.WriteString(fmt.Sprintf("%soperationCode: %v\n", i, f.operationCode))
	sb.WriteString(fmt.Sprintf("%sisAuthoritative: %v\n", i, f.isAuthoritative))
	sb.WriteString(fmt.Sprintf("%sisTruncated: %v\n", i, f.isTruncated))
	sb.WriteString(fmt.Sprintf("%sisRecursive: %v\n", i, f.isRecursive))
	sb.WriteString(fmt.Sprintf("%scanRecursive: %v\n", i, f.canRecursive))
	sb.WriteString(fmt.Sprintf("%sfutureUse: %v\n", i, f.futureUse))
	sb.WriteString(fmt.Sprintf("%sresponseCode: %v\n", i, f.responseCode))
	sb.WriteString(fmt.Sprintf("%s}", i))
	return sb.String()
}

func (f *Flags) toUint16() uint16 {
	var result uint16

	if !f.isQuery {
		result |= 1 << 15
	}

	result |= uint16(f.operationCode&0x0F) << 11

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
