package dns

import "errors"

type queryData struct {
	bytesToRead uint8
	data        []byte
}

func newQueryData(domainPart string) queryData {
	length := uint8(len(domainPart))

	return queryData{length, []byte(domainPart)}
}

func parseQueryData(bytes []byte) ([]queryData, int, error) {
	var qNames []queryData
	totalBytesRead := 0

	for i := 0; i < len(bytes); {
		bytesToRead := bytes[i]
		if bytesToRead == 0 {
			totalBytesRead = i + 1
			break
		}
		if i+int(bytesToRead)+1 > len(bytes) {
			return nil, 0, errors.New("invalid question name length")
		}
		qName := queryData{
			bytesToRead: bytesToRead,
			data:        bytes[i+1 : i+1+int(bytesToRead)],
		}
		qNames = append(qNames, qName)
		i += int(bytesToRead) + 1
	}

	return qNames, totalBytesRead, nil
}

func (qn queryData) toBytes() []byte {
	var bytes []byte

	bytes = append(bytes, qn.bytesToRead)
	bytes = append(bytes, qn.data...)

	return bytes
}

func (qn queryData) string() string {
	return string(qn.data)
}
