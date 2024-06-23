package main

import (
	"flag"
)

type cmdArgs struct {
	host   string
	client string
}

func setupFlags() *cmdArgs {
	args := new(cmdArgs)

	flag.StringVar(&args.host, "d", "google.com", "name of the domain")
	flag.StringVar(&args.client, "c", "1.1.1.1", "ip of the dns client")

	flag.Parse()

	return args
}
