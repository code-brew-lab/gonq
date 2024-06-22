package dns

type (
	RecordType  uint16
	RecordClass uint16
)

const (
	TypeA     RecordType = 1  // A record type.
	TypeNS    RecordType = 2  // Mail server record type.
	TypeCNAME RecordType = 5  // Canonical name of the domain.
	TypeMX    RecordType = 15 // Name server record type.

	ClassINET RecordClass = 1 // Internet Address class.
)

func (rt RecordType) TypeText() string {
	switch rt {
	case TypeA:
		return "A"
	case TypeNS:
		return "NS"
	case TypeCNAME:
		return "CNAME"
	case TypeMX:
		return "MX"
	default:
		return ""
	}
}

func (rc RecordClass) ClassText() string {
	switch rc {
	case ClassINET:
		return "InternetAddress"
	default:
		return ""
	}
}
