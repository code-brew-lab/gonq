package main

import (
	"fmt"
	"log"

	"github.com/code-brew-lab/gonq/pkg/dns"
)

func main() {
	args := setupFlags()

	req, err := dns.NewRequest(args.client, 53)
	if err != nil {
		log.Fatalln(err)
	}

	req.AddQuery(args.host, dns.TypeA, dns.ClassINET)

	resp, err := req.Make()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(req.Domains())
	fmt.Println(resp.IPs())
}
