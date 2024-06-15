package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
)

type DNSRequest struct {
	Header    DNSHeader
	Questions []DNSQuestion
}

type DNSHeader struct {
	ID      uint16 // [Req + Resp] 16 bit unique request id. Same for the following response.
	QR      bool   // [Req + Resp] Is the message query or response. Query = False; Response = True;
	OPCODE  uint8  // [Req] A four bit field that specifies kind of query in this message.
	AA      bool   // [Resp] Is response from authoritative dns server. Authoritative = True;
	TC      bool   // [Resp] Is response truncated or not.
	RD      bool   // [Req] Ask DNS server to recursively ask for the domain.
	RA      bool   // [Resp] Shows if recursion available for DNS server.
	Z       byte   // [Req] A three bit future use field.
	RCODE   byte   // [Resp] A four bits response codes.
	QDCOUNT uint16 // [Req] Number of entries inside the question.
	ANCOUNT uint16 // [Resp] Number of response entries from DNS server.
	NSCOUNT uint16 // [Req] Number of name server resource records in the authority records section.
	ARCOUNT uint16 // [Req] Number of resource records in the additional records section.
}

type DNSQuestion struct {
	QNAME  [2]byte
	QTYPE  [2]byte // Specifies the type of the query. Ex: CNAME, A, MX, NS
	QCLASS [2]byte // Specifies the class of the query.
}

func main() {
	dnsIP := net.ParseIP("8.8.8.8")
	if dnsIP == nil {
		log.Fatalln("Invalid IP")
	}

	rootDNS := &net.UDPAddr{
		IP:   dnsIP,
		Port: 53,
	}

	conn, err := net.DialUDP("udp", nil, rootDNS)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	msg, err := hex.DecodeString("064f0120000100000000000106676f6f676c6503636f6d00000100010000291000000000000000")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = conn.Write(msg)
	if err != nil {
		log.Fatalln(err)
	}

	buff := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buff)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Bytes[%d] Received: %x\n", n, buff[:n])
	}

}
