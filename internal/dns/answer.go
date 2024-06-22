package dns

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

type (
	answer struct {
		compressionType compressionType // 4 bit compression type
		offset          uint16          // 12 bit offset value
		ttl             uint32
		readLength      uint16
		data            []byte
	}

	compressionType uint8
)

const (
	pointer compressionType = 0xc
)

func parseAnswer(bytes []byte) (answer, int, error) {
	fmt.Println(hex.EncodeToString(bytes))
	var a answer
	be := binary.BigEndian

	firstByte := bytes[0]
	secondByte := bytes[1]

	a.compressionType = compressionType(firstByte >> 4)
	a.offset = uint16(firstByte&0x0F)<<8 | uint16(secondByte)
	a.ttl = be.Uint32((bytes[2:10]))
	a.readLength = be.Uint16(bytes[10:12])
	a.data = bytes[12 : 12+int(a.readLength)]

	return a, 6 + int(a.readLength), nil
}

func (a answer) string(indent int, char string) string {
	i := strings.Repeat(char, indent)

	var sb strings.Builder
	sb.WriteString(i + "Answer:\n")
	sb.WriteString(fmt.Sprintf("%sCompression: %v\n", i, a.compressionType))
	sb.WriteString(fmt.Sprintf("%sOffset: %v\n", i, a.offset))
	sb.WriteString(fmt.Sprintf("%sTTL: %d\n", i, a.ttl))
	sb.WriteString(fmt.Sprintf("%sData Length: %d\n", i, a.readLength))
	sb.WriteString(fmt.Sprintf("%sData: %d\n", i, a.data))
	return sb.String()
}
