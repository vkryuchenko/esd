package resolver

import (
	"github.com/miekg/dns"
	"log"
	"time"
)

func (z *Zone) AskParent(name string, qtype uint16) []dns.RR {
	for _, dnsserver := range z.Parent.Servers {
		log.Printf("Ask %s", dnsserver)
		dnsRequest := new(dns.Msg)
		dnsRequest.SetQuestion(name, qtype)
		dnsRequest.SetEdns0(4096, true)
		dnsClient := new(dns.Client)
		dnsClient.Timeout = time.Duration(z.Parent.Timeout) * time.Second
		answer, _, err := dnsClient.Exchange(dnsRequest, dnsserver)
		if err != nil {
			log.Println(err)
			continue
		}
		return answer.Answer
	}
	return nil
}
