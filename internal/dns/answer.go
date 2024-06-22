package dns

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type (
	answer struct {
		ttl        uint16
		readLength uint16
		data       []byte
	}
)

func parseAnswer(bytes []byte) (answer, int, error) {
	var a answer
	be := binary.BigEndian
	// q, qBytes, err := parseQuestion(bytes)
	// if err != nil {
	// 	return a, 0, err
	// }

	// if len(bytes) < qBytes+4 {
	// 	return a, 0, errors.New("not enough bytes to parse the answer")
	// }

	// a.question = q
	a.ttl = be.Uint16(bytes[:2])
	aBytes := be.Uint16(bytes[2:4])
	a.readLength = aBytes
	a.data = bytes[4 : 4+int(aBytes)]

	return a, 4 + int(a.readLength), nil
}

func (a answer) string(indent int, char string) string {
	i := strings.Repeat(char, indent)

	var sb strings.Builder
	sb.WriteString(i + "Answer:\n")
	// sb.WriteString(fmt.Sprintf("%s\tQuestion: %v\n", i, a.question))
	sb.WriteString(fmt.Sprintf("%s\tTTL: %d\n", i, a.ttl))
	sb.WriteString(fmt.Sprintf("%s\tData Length: %d\n", i, a.readLength))
	sb.WriteString(fmt.Sprintf("%s\tData: %s\n", i, string(a.data)))
	return sb.String()
}
