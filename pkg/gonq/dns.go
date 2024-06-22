package gonq

import "github.com/code-brew-lab/gonq.git/internal/dns"

type (
	Response interface {
		String() string
	}
)

func NewRequest(distIP string, port int) (dns.Request, error) {
	return dns.NewRequest(distIP, port)
}
