package resolver

import (
	"github.com/miekg/dns"
)

func Start(listen, protocol string, zones []*Zone) error {
	mux := dns.NewServeMux()
	for _, zone := range zones {
		mux.Handle(zone.Root, zone)
	}
	server := dns.Server{
		Addr:    listen,
		Net:     protocol,
		Handler: mux,
	}
	err := server.ListenAndServe()
	defer server.Shutdown()
	return err
}
