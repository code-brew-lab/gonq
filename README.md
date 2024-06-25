# Gonq
Gonq (Go Name Query) is a simple command-line application written in Go that mimics the functionality of the dig tool. It takes a domain name as input and returns the corresponding IP address. Gonq is designed to be lightweight and efficient, providing a quick and easy way to query domain names and retrieve their IP addresses.

### Installation

#### CLI Tool

You can install the CLI tool using `go install`:

```sh
go install github.com/code-brew-lab/gonq/cmd/gonq@v1.0.0
```

or you can compile it manually by using the Makefile. See Makefile for more details
```sh
make build
```
#### Go Module
To use the Go module in your project, import it as follows:
```go
import "github.com/code-brew-lab/gonq/pkg/dns"
```
### Usage
#### CLI Tool
Perform a DNS query by:
```sh
gonq -d google.com
```
Response:
```sh
Raw: e0808180000100010000000006676f6f676c6503636f6d0000010001c00c0001000100000e100004d8ef2678
Header: ID: 57472
        Flags: IsQuery: false, QueryKind: StandardQuery, IsAuthoritative: false, IsTruncated: false, IsRecursive: true, CanRecursive: true, ResponseCode: NoError
        QDCount: 1
        ANCount: 1
        NSCount: 0
        ARCount: 0
Queries: [Domain: google.com, RecordType: A, RecordClass: InternetAddress]
Answers: [CompressionType: 12, Offset: 12, TTL: 65537s, ReadLength: 4, IP: 216.239.38.120]
```
#### Go Module
Example usage of go module:
```go
package main

import (
    "log"
    "github.com/code-brew-lab/gonq/pkg/dns"
)

func main() {
    req, err := dns.NewRequest("1.1.1.1", 53)
    if err != nil {
        log.Fatalln(err)
    }

    req.AddQuery("google.com", dns.TypeA, dns.ClassINET)

    resp, err := req.Make()
    if err != nil {
        log.Fatalln(err)
    }
}
```
### Documentation
Comprehensive documentation for the Go module is available at [pkg.go.dev](https://pkg.go.dev/github.com/code-brew-lab/gonq@v1.0.0/pkg/dns).

### Feedback
Your feedback is important to us. If you encounter any issues or have suggestions for improvement, please open an issue on our [GitHub repository](https://github.com/code-brew-lab/gonq/issues).

Thank you for using `gonq`!
