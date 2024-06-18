package dns

type (
	Request struct {
		Header    Header
		Questions []Question
	}

	Question struct {
		QuestionNames []QuestionName
		QuestionType  [2]byte // Specifies the type of the query. Ex: CNAME, A, MX, NS
		QuestionClass [2]byte // Specifies the class of the query.
	}

	QuestionName struct {
		BytesToRead uint8
		Data        []byte
	}
)
