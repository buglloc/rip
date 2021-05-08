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
		switch v := rr.(type) {
		case *dns.A:
			ipv4 = append(ipv4, v.A)
		case *dns.AAAA:
			ipv6 = append(ipv6, v.AAAA)
		}
	}

	ttl := time.Duration(res.Answer[0].(dns.RR).Header().Ttl) * time.Second
	if reqType == dns.TypeA {
		dnsCache.Set(dns.TypeA, name, ttl, ipv4)
		return ipv4, nil
	}

	dnsCache.Set(dns.TypeAAAA, name, ttl, ipv6)
	return ipv6, nil
}
