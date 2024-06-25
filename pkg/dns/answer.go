package dns

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

type answer struct {
	compressionType uint16 // 4 bit compression type
	offset          uint16 // 12 bit offset value
	ttl             uint32
	readLength      uint16
	ip              net.IP
}

const answerHeaderSize int = 12

func parseAnswer(bytes []byte) (answer, int, error) {
	var a answer
	be := binary.BigEndian

	if len(bytes) < answerHeaderSize {
		return a, 0, errors.New("not enough bytes to parse the answer")
	}

	firstByte := bytes[0]
	secondByte := bytes[1]

	a.compressionType = uint16(firstByte >> 4)
	a.offset = uint16(firstByte&0x0F)<<8 | uint16(secondByte)
	a.ttl = be.Uint32((bytes[2:10]))
	a.readLength = be.Uint16(bytes[10:12])

	if len(bytes[answerHeaderSize:]) < int(a.readLength) {
		return a, 0, errors.New("not enough bytes to parse the answer")
	}

	a.ip = bytes[answerHeaderSize : answerHeaderSize+int(a.readLength)]

	return a, answerHeaderSize + int(a.readLength), nil
}

func (a answer) string() string {
	return fmt.Sprintf("CompressionType: %d, Offset: %d, TTL: %ds, ReadLength: %d, IP: %v",
		a.compressionType, a.offset, a.ttl, a.readLength, a.ip)
}
