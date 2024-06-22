package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/code-brew-lab/gonq.git/internal/dns"
)

func main() {
	req := dns.NewRequest()
	req.AddQuery("google.com", dns.TypeA, dns.ClassINET)

	reqBytes, err := req.ToBytes()
	if err != nil {
		log.Fatalln("Error marshaling request:", err)
	}

	fmt.Println("Request in hex:", hex.EncodeToString(reqBytes))

	// Send the DNS query
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatalln("Error dialing:", err)
	}
	defer conn.Close()

	// Set a read deadline
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))

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

	resp, err := dns.ParseResponse(buf[:n])
	if err != nil {
		log.Fatalln("Error parsing response:", err)
	}

	fmt.Println(resp.String())

}
