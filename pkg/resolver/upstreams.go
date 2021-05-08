package resolver

import (
	"net"
	"time"

	"github.com/miekg/dns"

	"github.com/buglloc/rip/v2/pkg/cfg"
)

var (
	dnsClient = &dns.Client{
		Net:          "tcp",
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
	}
	dnsCache = NewCache()
)

func ResolveIp(reqType uint16, name string) ([]net.IP, error) {
	if ips := dnsCache.Get(reqType, name); ips != nil {
		return ips, nil
	}

	msg := &dns.Msg{}
	msg.SetQuestion(dns.Fqdn(name), reqType)
	res, _, err := dnsClient.Exchange(msg, cfg.Upstream)
	if err != nil || len(res.Answer) == 0 {
		return nil, err
	}

	var ipv4 []net.IP
	var ipv6 []net.IP
	for _, rr := range res.Answer {
		switch rr.(type) {
		case *dns.A:
			ipv4 = append(ipv4, rr.(*dns.A).A)
		case *dns.AAAA:
			ipv6 = append(ipv6, rr.(*dns.AAAA).AAAA)
		}
	}

	ttl := time.Duration(res.Answer[0].(dns.RR).Header().Ttl) * time.Second
	dnsCache.Set(dns.TypeA, name, ttl, ipv4)
	dnsCache.Set(dns.TypeAAAA, name, ttl, ipv6)

	if reqType == dns.TypeA {
		return ipv4, nil
	}

	return ipv6, nil
}
