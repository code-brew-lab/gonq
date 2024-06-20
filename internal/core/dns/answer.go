package dns

type (
	answer struct {
		question    question
		answerType  RecordType
		answerClass RecordClass
		ttl         uint16
		readLength  uint16
	}
)
