package main

import (
	"flag"
)

type cmdArgs struct {
	host string
}

func setupFlags() (*cmdArgs, error) {
	args := new(cmdArgs)

	flag.StringVar(&args.host, "", "google.com", "name of the domain")
	flag.Parse()

	return args, nil
}
