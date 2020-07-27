package resolver

import (
	"github.com/miekg/dns"
	"log"
	"math/rand"
)

func (z *Zone) AskParent(name string, qtype uint16) []dns.RR {
	dnsserver := z.ParentDNS[rand.Intn(len(z.ParentDNS))]
	log.Printf("Ask %s", dnsserver)
	dnsRequest := new(dns.Msg)
	dnsRequest.SetQuestion(name, qtype)
	dnsRequest.SetEdns0(4096, true)
	dnsClient := new(dns.Client)
	answer, _, err := dnsClient.Exchange(dnsRequest, dnsserver)
	if err != nil {
		log.Println(err)
		return nil
	}
	return answer.Answer
}
