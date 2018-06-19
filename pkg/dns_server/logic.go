package dns_server

import (
	"net"
	"strings"

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
		if err != nil {
			log.Error("failed to resolve proxied name",
				"qtype", typeToString(question.Qtype), "name", question.Name, "err", err.Error())
			break
		}
		log.Info("cooking response",
			"mode", "proxy", "qtype", typeToString(question.Qtype), "name", question.Name, "ip", ip.String())
	case suffix == ".r":
		ips := strings.Split(name, ".")
		if len(ips) < 2 {
			log.Error("failed to parse random annotation",
				"qtype", typeToString(question.Qtype), "name", question.Name)
			break
		}

		var pIp string
		if random(0, 100) > 50 {
			pIp = ips[len(ips)-2]
		} else {
			pIp = ips[len(ips)-1]
		}
		ip = parseIp(question.Qtype, pIp)
		log.Info("cooking response",
			"mode", "random", "qtype", typeToString(question.Qtype), "name", question.Name, "ip", ip.String())
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
		log.Info("cooking response",
			"mode", "cname", "qtype", typeToString(question.Qtype), "name", question.Name, "target", name)
	case suffix == ".4":
		if question.Qtype == dns.TypeA {
			ip = parseIp(dns.TypeA, name)
		} else if !cfg.StrictMode {
			ip = defaultIp(question.Qtype)
		}
		log.Info("cooking response",
			"mode", "ipv6", "qtype", typeToString(question.Qtype), "name", question.Name, "ip", ip.String())
	case suffix == ".6":
		if question.Qtype == dns.TypeAAAA {
			ip = parseIp(dns.TypeAAAA, name)
		} else if !cfg.StrictMode {
			ip = defaultIp(question.Qtype)
		}
		log.Info("cooking response",
			"mode", "ipv4", "qtype", typeToString(question.Qtype), "name", question.Name, "ip", ip.String())
	default:
		ip = defaultIp(question.Qtype)
		log.Info("cooking response",
			"mode", "default", "qtype", typeToString(question.Qtype), "name", question.Name, "ip", ip.String())
	}

	return
}
