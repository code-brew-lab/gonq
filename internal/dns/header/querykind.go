package header

type QueryKind uint8

const (
	KindStandard QueryKind = iota
)

func (qk QueryKind) KindText() string {
	switch qk {
	case KindStandard:
		return "StandardQuery"
	default:
		return ""
	}
}
