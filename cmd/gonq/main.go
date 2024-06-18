package main

import (
	"fmt"
	"log"

	"github.com/code-brew-lab/gonq.git/internal/core/dns"
	"github.com/code-brew-lab/gonq.git/internal/pkg/bitwise"
)

func main() {
	flags := dns.Flags{
		IsQuery:         bitwise.New(false),                                                                            // Query
		OperationCode:   [4]bitwise.Bit{bitwise.New(false), bitwise.New(false), bitwise.New(true), bitwise.New(false)}, // Example: 0010
		IsAuthoritative: bitwise.New(true),                                                                             // Authoritative
		IsTruncated:     bitwise.New(false),                                                                            // Not Truncated
		IsRecursive:     bitwise.New(true),                                                                             // Recursive
		CanRecursive:    bitwise.New(true),                                                                             // Recursion Available
		FutureUse:       [3]bitwise.Bit{bitwise.New(false), bitwise.New(false), bitwise.New(false)},                    // Future use: 000
		ResponseCode:    [4]bitwise.Bit{bitwise.New(false), bitwise.New(true), bitwise.New(false), bitwise.New(true)},  // Example: 0101
	}

	bytes, err := flags.BinaryMarshaler()
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < len(bytes); i++ {
		fmt.Printf("%08b\n", bytes[i])
	}
}

// Wireshark response: 5c6381800001000500000000076773702d73736c086c732d6170706c6503636f6d06616b61646e73036e65740000010001c00c000500010000001800110e6773702d73736c2d67656f6d6170c014c03d0005000100000032001708677370782d73736c026c73056170706c6503636f6d00c05a00050001000004ff0013066765742d62780167076161706c696d67c06cc07d000100010000000d000411fd49cfc07d000100010000000d000411fd49d0

// Gonq response: 5c6381800001000500000000076773702d73736c086c732d6170706c6503636f6d06616b61646e73036e65740000010001c00c000500010000001300110e6773702d73736c2d67656f6d6170c014c03d0005000100000031001708677370782d73736c026c73056170706c6503636f6d00c05a0005000100000e010013066765742d62780167076161706c696d67c06cc07d0001000100000004000411fd49d0c07d0001000100000004000411fd49cf
