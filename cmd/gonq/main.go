package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"

	"github.com/code-brew-lab/gonq.git/internal/core/dns"
)

func main() {
	flags := dns.NewFlagsBuilder().SetIsQuery(true).SetIsRecursive(true).Build()
	header := dns.NewHeaderBuilder().SetID(0).SetFlags(flags).Build()

	req, err := dns.NewRequest(header)
	if err != nil {
		log.Fatalln(err)
	}
	req.AddQuestion("google.com", dns.AType, dns.INetClass)

	reqBytes, err := req.BinaryMarshaler()
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < len(reqBytes); i += 2 {
		fmt.Printf("%s\n", hex.EncodeToString(reqBytes[i:i+2]))
	}

	// Send the DNS query
	conn, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(reqBytes)
	if err != nil {
		fmt.Println("Error writing to connection:", err)
		return
	}

	// Read the response
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	// Print the response in hex format
	fmt.Println(hex.EncodeToString(buf[:n]))
}
