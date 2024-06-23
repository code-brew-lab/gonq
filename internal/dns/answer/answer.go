package answer

import (
	"encoding/binary"
	"errors"
)

type Answer struct {
	compressionType uint16 // 4 bit compression type
	offset          uint16 // 12 bit offset value
	ttl             uint32
	readLength      uint16
	data            []byte
}

const answerHeaderSize int = 12

func Parse(bytes []byte) (Answer, int, error) {
	var a Answer
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

	a.data = bytes[answerHeaderSize : answerHeaderSize+int(a.readLength)]

	return a, answerHeaderSize + int(a.readLength), nil
}
