package resolver

import (
	"time"

	"github.com/miekg/dns"
)

func (z *Zone) AskParent(name string, qtype uint16) []dns.RR {
	for _, dnsserver := range z.Parent.Servers {
		z.Logger.Debugf("ask %s about %s", dnsserver, name)
		dnsRequest := new(dns.Msg)
		dnsRequest.SetQuestion(name, qtype)
		dnsRequest.SetEdns0(4096, true)
		dnsClient := new(dns.Client)
		dnsClient.Timeout = time.Duration(z.Parent.Timeout) * time.Second
		answer, _, err := dnsClient.Exchange(dnsRequest, dnsserver)
		if err != nil {
			z.Logger.Error(err)
			continue
		}
		return answer.Answer
	}
	return nil
}
