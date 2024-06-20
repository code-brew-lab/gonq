package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/code-brew-lab/gonq.git/internal/core/dns"
)

func main() {
	req := dns.NewRequest()
	req.AddQuestion("google.com", dns.AType, dns.INetClass)

	reqBytes, err := req.BinaryMarshaler()
	if err != nil {
		log.Fatalln("Error marshaling request:", err)
	}

	fmt.Println("Request in hex:", hex.EncodeToString(reqBytes))

	// Send the DNS query
	conn, err := net.Dial("udp", "1.1.1.1:53")
	if err != nil {
		log.Fatalln("Error dialing:", err)
	}
	defer conn.Close()

	// Set a read deadline
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	_, err = conn.Write(reqBytes)
	if err != nil {
		log.Fatalln("Error writing to connection:", err)
	}

	// Read the response
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln("Error reading from connection:", err)
	}

	// Print the response in hex format
	fmt.Println("Response in hex:", hex.EncodeToString(buf[:n]))
}
