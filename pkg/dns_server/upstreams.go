package dns_server

import (
	"net"
	"time"

	"github.com/miekg/dns"

	"github.com/buglloc/rip/pkg/cfg"
	"github.com/buglloc/rip/pkg/dns_cache"
)

var (
	dnsClient *dns.Client
	cache     *dns_cache.Cache
)

func init() {
	cache = dns_cache.NewCache()
	dnsClient = &dns.Client{
		Net:          "tcp",
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
	}
}

func ResolveIp(reqType uint16, name string) (net.IP, error) {
	if ip := cache.Get(reqType, name); ip != nil {
		return *ip, nil
	}

	msg := &dns.Msg{}
	msg.SetQuestion(dns.Fqdn(name), reqType)
	res, _, err := dnsClient.Exchange(msg, cfg.Upstream)
	if err != nil || len(res.Answer) == 0 {
		return nil, err
	}

	var ip net.IP
	for _, rr := range res.Answer {
		ttl := time.Duration(rr.(dns.RR).Header().Ttl) * time.Second
		switch rr.(type) {
		case *dns.A:
			rip := rr.(*dns.A).A
			if reqType == dns.TypeA {
				ip = rip
			}
			cache.Set(dns.TypeA, name, ttl, &rip)
		case *dns.AAAA:
			rip := rr.(*dns.AAAA).AAAA
			if reqType == dns.TypeAAAA {
				ip = rip
			}
			cache.Set(dns.TypeAAAA, name, ttl, &rip)
		}
	}

	return ip, nil
}
