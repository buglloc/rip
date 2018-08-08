package handlers

import (
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"

	"github.com/buglloc/rip/pkg/cfg"
	"github.com/buglloc/rip/pkg/ip_utils"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func defaultIp(reqType uint16) net.IP {
	if reqType == dns.TypeA {
		return cfg.IPv4
	}
	return cfg.IPv6
}

func parseIp(reqType uint16, name string) net.IP {
	if indx := strings.LastIndex(name, "."); indx != -1 {
		name = name[indx+1:]
	}

	dotCounts := strings.Count(name, "-")

	switch reqType {
	case dns.TypeA:
		if dotCounts == 0 {
			// base-16 form
			ip := ip_utils.HexToIp(name)
			if len(ip) != net.IPv4len {
				return nil
			}
			return ip
		} else if dotCounts != 3 {
			return defaultIp(dns.TypeA)
		}
		return net.ParseIP(strings.Replace(name, "-", ".", -1))
	case dns.TypeAAAA:
		if dotCounts == 0 {
			// base-16 form
			ip := ip_utils.HexToIp(name)
			if len(ip) != net.IPv6len {
				return nil
			}
			return ip
		} else if dotCounts < 2 {
			return defaultIp(dns.TypeAAAA)
		}
		return net.ParseIP(strings.Replace(name, "-", ":", -1))
	default:
		return defaultIp(dns.TypeA)
	}
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func parseSubName(name string) string {
	if indx := strings.LastIndex(name, "."); indx != -1 {
		part := name[indx+1:]
		if strings.Count(part, "-") > 0 {
			// We have request like something.example-com.c.example.com
			return strings.Replace(part, "-", ".", -1)
		}
	}
	return name
}

func createIpRR(question dns.Question, ip net.IP) (rr dns.RR) {
	head := dns.RR_Header{
		Name:   question.Name,
		Rrtype: question.Qtype,
		Class:  dns.ClassINET,
		Ttl:    cfg.TTL,
	}

	if question.Qtype == dns.TypeA {
		rr = &dns.A{
			Hdr: head,
			A:   ip,
		}
	} else {
		rr = &dns.AAAA{
			Hdr:  head,
			AAAA: ip,
		}
	}
	return
}

func createIpsRR(question dns.Question, ips ...net.IP) (result []dns.RR) {
	result = make([]dns.RR, len(ips))
	for i, ip := range ips {
		result[i] = createIpRR(question, ip)
	}

	return
}
