package header

type ResponseCode uint8

const (
	CodeNoError        ResponseCode = iota // #0  OK.
	CodeFormatError                        // #1  The name server was unable to interpret the query.
	CodeServerFailure                      // #2  Server unreachable.
	CodeNameError                          // #3  Domain name does not exists.
	CodeNotImplemented                     // #4  Unsupported request query type.
	CodeRefused                            // #5  Refused for policy reasons of the server.
)

func (rc ResponseCode) CodeText() string {
	switch rc {
	case CodeNoError:
		return "NoError"
	case CodeFormatError:
		return "FormatError"
	case CodeServerFailure:
		return "ServerFailure"
	case CodeNameError:
		return "NameError"
	case CodeNotImplemented:
		return "NotImplemented"
	case CodeRefused:
		return "Refused"
	default:
		return ""
	}
}
