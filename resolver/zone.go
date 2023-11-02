package resolver

import (
	"strings"

	"github.com/miekg/dns"
	"go.uber.org/zap"
)

type Parent struct {
	Timeout uint64   `yaml:"timeout"`
	Servers []string `yaml:"servers"`
}

type Record struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	PointTo []string `yaml:"pointTo"`
}

type Zone struct {
	Logger  *zap.SugaredLogger `yaml:"-"`
	Root    string             `yaml:"root"`
	Parent  Parent             `yaml:"-"`
	Records []Record           `yaml:"records"`
}

func (z *Zone) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		z.parseQuery(m, w.RemoteAddr().String())
	}

	err := w.WriteMsg(m)
	if err != nil {
		z.Logger.Errorf("failed to write message for %s: %s", w.RemoteAddr().String(), err)
	}
}

func (z *Zone) parseQuery(m *dns.Msg, client string) {
	for _, q := range m.Question {
		questionType := dns.Type(q.Qtype).String()
		z.Logger.Debugf("%s query for %s from %s", questionType, q.Name, client)
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
			z.Logger.Debugf("%s points for %s found: %s", questionType, q.Name, points)
			for _, point := range points {
				rr, err := dns.NewRR(strings.Join([]string{q.Name, questionType, point}, " "))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		} else {
			z.Logger.Warnf("%s points not found for %s requested by %s", questionType, q.Name, client)
			m.Answer = z.AskParent(q.Name, q.Qtype)
		}
	}
}
