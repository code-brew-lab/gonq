package main

import (
	"fmt"
	"log"

	"github.com/code-brew-lab/gonq.git/internal/dns"
)

func main() {
	req, err := dns.NewRequest("1.1.1.1", 53)
	if err != nil {
		log.Fatalln(err)
	}

	req.AddQuery("gokhanuysal.net", dns.TypeA, dns.ClassINET)

	resp, err := req.Make()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp.ID())
}
