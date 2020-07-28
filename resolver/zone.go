package resolver

import (
	"github.com/miekg/dns"
	"log"
	"strings"
)

type Record struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	PointTo []string `yaml:"pointTo"`
}

type Zone struct {
	Root      string   `yaml:"root"`
	ParentDNS []string `yaml:"-"`
	Records   []Record `yaml:"records"`
}

func (z Zone) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		z.parseQuery(m)
	}

	err := w.WriteMsg(m)
	if err != nil {
		log.Printf("Failed to write message:%s\n", err)
	}
}

func (z *Zone) parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		questionType := dns.Type(q.Qtype).String()
		log.Printf("Query %s for %s\n", questionType, q.Name)
		points := []string{}
		for _, record := range z.Records {
			if q.Name != record.Name {
				continue
			}
			if questionType != record.Type {
				continue
			}
			for _, point := range record.PointTo {
				points = append(points, point)
			}
		}
		if len(points) > 0 {
			log.Printf("Points found: %s\n", points)
			for _, point := range points {
				rr, err := dns.NewRR(strings.Join([]string{q.Name, questionType, point}, " "))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		} else {
			log.Println("Points not found")
			m.Answer = z.AskParent(q.Name, q.Qtype)
		}
	}
}
