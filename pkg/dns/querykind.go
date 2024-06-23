package dns

type queryKind uint8

const (
	kindStandard queryKind = iota
)

func (qk queryKind) kindText() string {
	switch qk {
	case kindStandard:
		return "StandardQuery"
	default:
		return ""
	}
}
