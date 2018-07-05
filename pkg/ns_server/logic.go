package ns_server

import (
	"net"
	"strings"

	"github.com/buglloc/simplelog"
	"github.com/miekg/dns"

	"github.com/buglloc/rip/pkg/cfg"
	"github.com/buglloc/rip/pkg/ip_loop"
)

func parseName(question dns.Question, zone string, l log.Logger) (msg dns.RR, ip net.IP, err error) {
	if len(question.Name)-len(zone) <= 3 {
		ip = defaultIp(question.Qtype)
		return
	}

	l = l.Child("qtype", typeToString(question.Qtype), "name", question.Name)
	name := question.Name[:len(question.Name)-len(zone)-1]
	i := len(name) - 2
	name, suffix := name[:i], name[i:]
	switch {
	case suffix == ".p" && cfg.AllowProxy:
		subName := parseSubName(name)
		ip, err = ResolveIp(question.Qtype, subName)
		if err != nil {
			l.Error("failed to resolve proxied name", "target", subName, "err", err.Error())
			break
		}
		l.Info("cooking response",
			"mode", "proxy", "target", subName, "ip", ip.String())
	case suffix == ".l":
		ips := strings.Split(name, ".")
		if len(ips) < 2 {
			log.Error("failed to parse loop annotation")
			break
		}

		ips = ips[len(ips)-2:]
		var pIp string
		if question.Qtype == dns.TypeA {
			// Move next only for ipv4 request
			pIp = ip_loop.GetNext(name, ips)
		} else {
			pIp = ip_loop.GetCurrent(name, ips)
		}
		ip = parseIp(question.Qtype, pIp)
		l.Info("cooking response", "mode", "loop", "ip", ip.String())
	case suffix == ".r":
		ips := strings.Split(name, ".")
		if len(ips) < 2 {
			l.Error("failed to parse random annotation")
			break
		}

		var pIp string
		if random(0, 100) > 50 {
			pIp = ips[len(ips)-2]
		} else {
			pIp = ips[len(ips)-1]
		}
		ip = parseIp(question.Qtype, pIp)
		l.Info("cooking response", "mode", "random", "ip", ip.String())
	case suffix == ".c":
		subName := parseSubName(name)
		msg = &dns.CNAME{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: dns.TypeCNAME,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			Target: dns.Fqdn(subName),
		}
		l.Info("cooking response", "mode", "cname", "target", subName)
	case suffix == ".4":
		if question.Qtype == dns.TypeA {
			ip = parseIp(dns.TypeA, name)
		} else if !cfg.StrictMode {
			ip = defaultIp(question.Qtype)
		}
		l.Info("cooking response",
			"mode", "ipv6", "ip", ip.String())
	case suffix == ".6":
		if question.Qtype == dns.TypeAAAA {
			ip = parseIp(dns.TypeAAAA, name)
		} else if !cfg.StrictMode {
			ip = defaultIp(question.Qtype)
		}
		l.Info("cooking response",
			"mode", "ipv4", "ip", ip.String())
	default:
		ip = defaultIp(question.Qtype)
		l.Info("cooking response",
			"mode", "default", "ip", ip.String())
	}

	return
}
