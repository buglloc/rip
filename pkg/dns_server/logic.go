package dns_server

import (
	"net"

	"github.com/buglloc/simplelog"
	"github.com/miekg/dns"

	"github.com/buglloc/rip/pkg/cfg"
)

func parseName(question dns.Question, zone string) (msg dns.RR, ip net.IP, err error) {
	if len(question.Name)-len(zone) <= 3 {
		ip = defaultIp(question.Qtype)
		return
	}

	name := question.Name[:len(question.Name)-len(zone)-1]
	i := len(name) - 2
	name, suffix := name[:i], name[i:]
	switch {
	case suffix == ".p" && cfg.AllowProxy:
		ip, err = ResolveIp(question.Qtype, name)
		log.Info("cooking response", "type", "proxy", "name", question.Name, "ip", ip.String())
	case suffix == ".c":
		msg = &dns.CNAME{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: dns.TypeCNAME,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			Target: dns.Fqdn(name),
		}
		log.Info("cooking response", "type", "cname", "name", question.Name, "target", name)
	case suffix == ".4":
		if question.Qtype == dns.TypeA {
			ip = parseIp(dns.TypeA, name)
		} else if !cfg.StrictMode {
			ip = defaultIp(question.Qtype)
		}
		log.Info("cooking response", "type", "A", "name", question.Name, "ip", ip.String())
	case suffix == ".6":
		if question.Qtype == dns.TypeAAAA {
			ip = parseIp(dns.TypeAAAA, name)
		} else if !cfg.StrictMode {
			ip = defaultIp(question.Qtype)
		}
		log.Info("cooking response", "type", "AAAA", "name", question.Name, "ip", ip.String())
	default:
		ip = defaultIp(question.Qtype)
		log.Info("cooking response", "type", typeToString(question.Qtype), "name", question.Name, "ip", ip.String())
	}

	return
}
