package dns

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type (
	request struct {
		addr      *net.UDPAddr
		header    *header
		questions []query
	}

	Request interface {
		AddQuery(domain string, rType RecordType, rClass RecordClass)
		Make() (Response, error)
	}
)

func NewRequest(distIP string, port int) (Request, error) {
	flags := newFlagsBuilder().SetIsQuery(true).SetIsRecursive(true).Build()
	header := newHeaderBuilder().SetID(NewID()).SetFlags(flags).Build()

	ip := net.ParseIP(distIP)
	if ip == nil {
		return nil, errors.New("dns.NewRequest: invalid destination ip")
	}

	if port < 0 {
		return nil, errors.New("dns.NewRequest: invalid destination port")
	}

	addr := &net.UDPAddr{
		IP:   ip,
		Port: port,
	}

	return &request{
		addr:   addr,
		header: header,
	}, nil
}

func (r *request) AddQuery(domain string, rType RecordType, rClass RecordClass) {
	r.questions = append(r.questions, newQuery(domain, rType, rClass))
	r.header.addQuestion()
}

func (r *request) Make() (Response, error) {
	if r == nil {
		return nil, errors.New("dns.Make: request is nil")
	}

	reqBytes := r.ToBytes()

	conn, err := net.DialUDP("udp", nil, r.addr)
	if err != nil {
		return nil, fmt.Errorf("dns.Make: %v", err)
	}

	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))

	_, err = conn.Write(reqBytes)
	if err != nil {
		return nil, fmt.Errorf("dns.Make: %v", err)
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("dns.Make error reading from connection: %v", err)
	}

	resp, err := parseResponse(buf[:n])
	if err != nil {
		return nil, fmt.Errorf("dns.Make error parsing response: %v", err)
	}
	defer conn.Close()

	return resp, nil
}

func (r *request) ToBytes() []byte {
	var bytes []byte

	header := r.header
	questions := r.questions

	bytes = append(bytes, header.toBytes()...)
	for _, q := range questions {
		bytes = append(bytes, q.toBytes()...)
	}

	return bytes
}
