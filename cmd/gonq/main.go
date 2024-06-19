package main

import (
	"fmt"
	"log"

	"github.com/code-brew-lab/gonq.git/internal/core/dns"
)

func main() {
	flags := dns.NewFlagsBuilder().SetIsTruncated(true).Build()
	header := dns.NewHeaderBuilder().SetID(0).SetFlags(flags).Build()

	bytes, err := header.BinaryMarshaler()
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < len(bytes); i += 2 {
		fmt.Printf("%08b %08b\n", bytes[i], bytes[i+1])
	}
}
