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

	rr := res.Answer[0]
	var ip net.IP
	ttl := time.Duration(rr.(dns.RR).Header().Ttl) * time.Second
	switch rr.(type) {
	case *dns.A:
		ip = rr.(*dns.A).A
	case *dns.AAAA:
		ip = rr.(*dns.AAAA).AAAA
	case *dns.CNAME:
		ip, err = ResolveIp(reqType, rr.(*dns.CNAME).Target)
	}

	if ip != nil {
		cache.Set(reqType, name, ttl, &ip)
	}
	return ip, nil
}
