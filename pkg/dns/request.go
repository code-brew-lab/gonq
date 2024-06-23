package dns

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type (
	Request struct {
		addr    *net.UDPAddr
		header  *Header
		queries []query
	}
)

func NewRequestWithHeader(header *Header, distIP string, port int) (*Request, error) {
	if header == nil {
		return nil, errors.New("dns.NewRequestWithHeader: header is nil")
	}

	ip := net.ParseIP(distIP)
	if ip == nil {
		return nil, errors.New("dns.NewRequestWithHeader: invalid destination ip")
	}

	if port < 0 {
		return nil, errors.New("dns.NewRequestWithHeader: invalid destination port")
	}

	addr := &net.UDPAddr{
		IP:   ip,
		Port: port,
	}

	return &Request{
		addr:   addr,
		header: header,
	}, nil
}

func NewRequest(distIP string, port int) (*Request, error) {
	flags := NewFlagsBuilder().SetIsQuery(true).SetIsRecursive(true).Build()
	header := NewHeaderBuilder().SetID(NewID()).SetFlags(flags).Build()

	return NewRequestWithHeader(header, distIP, port)
}

func (r *Request) AddQuery(domain string, rType RecordType, rClass RecordClass) {
	r.queries = append(r.queries, newQuery(domain, rType, rClass))
	r.header.AddQuestion()
}

func (r *Request) Make() (*Response, error) {
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

func (r *Request) Domains() []string {
	var domains []string
	for _, d := range r.queries {
		domains = append(domains, d.domain())
	}

	return domains
}

func (r *Request) ToBytes() []byte {
	var bytes []byte

	header := r.header
	questions := r.queries

	bytes = append(bytes, header.ToBytes()...)
	for _, q := range questions {
		bytes = append(bytes, q.toBytes()...)
	}

	return bytes
}

func (r *Request) String() string {
	var queries []string
	for _, q := range r.queries {
		queries = append(queries, q.string())
	}
	return fmt.Sprintf("Address: %v\nHeader: %v\n Queries: [%v]\n", r.addr, r.header, strings.Join(queries, ", "))
}
