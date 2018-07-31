package ns_server

import (
	"net"
	"strings"

	"github.com/buglloc/simplelog"
	"github.com/miekg/dns"

	"github.com/buglloc/rip/pkg/cfg"
	"github.com/buglloc/rip/pkg/ip_loop"
)

func parseName(question dns.Question, zone string, l log.Logger) (rrs []dns.RR, err error) {
	//TODO(buglloc): refactor this holly shit!!!

	if len(question.Name)-len(zone) <= 3 {
		ip := defaultIp(question.Qtype)
		rrs = []dns.RR{createIpMsg(question, ip)}
		l.Info("cooking response",
			"mode", "default", "target", question.Name, "ip", ip.String())
		return
	}

	l = l.Child("qtype", typeToString(question.Qtype), "name", question.Name)
	name := question.Name[:len(question.Name)-len(zone)-1]
	i := strings.LastIndex(name, ".")
	name, suffix := name[:i], name[i+1:]
	switch {
	case suffix == "p" && cfg.AllowProxy:
		subName := parseSubName(name)
		ip, rerr := ResolveIp(question.Qtype, subName)
		if rerr != nil {
			err = rerr
			l.Error("failed to resolve proxied name", "target", subName, "err", err.Error())
			break
		}
		rrs = []dns.RR{createIpMsg(question, ip)}
		l.Info("cooking response",
			"mode", "proxy", "target", subName, "ip", ip.String())
	case suffix == "l":
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
		ip := parseIp(question.Qtype, pIp)
		rrs = []dns.RR{createIpMsg(question, ip)}
		l.Info("cooking response", "mode", "loop", "ip", ip.String())
	case suffix == "r":
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
		ip := parseIp(question.Qtype, pIp)
		rrs = []dns.RR{createIpMsg(question, ip)}
		l.Info("cooking response", "mode", "random", "ip", ip.String())
	case suffix == "c":
		subName := parseSubName(name)
		rrs = []dns.RR{&dns.CNAME{
			Hdr: dns.RR_Header{
				Name:   question.Name,
				Rrtype: dns.TypeCNAME,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			Target: dns.Fqdn(subName),
		}}
		l.Info("cooking response", "mode", "cname", "target", subName)
	case suffix == "4":
		var ip net.IP
		if question.Qtype == dns.TypeA {
			ip = parseIp(dns.TypeA, name)
		} else if cfg.UseDefault {
			ip = defaultIp(question.Qtype)
		}

		if !ip.IsUnspecified() {
			rrs = []dns.RR{createIpMsg(question, ip)}
			l.Info("cooking response",
				"mode", "ipv4", "ip", ip.String())
		}
	case suffix == "6":
		var ip net.IP
		if question.Qtype == dns.TypeAAAA {
			ip = parseIp(dns.TypeAAAA, name)
		} else if cfg.UseDefault {
			ip = defaultIp(question.Qtype)
		}

		if !ip.IsUnspecified() {
			rrs = []dns.RR{createIpMsg(question, ip)}
			l.Info("cooking response",
				"mode", "ipv6", "ip", ip.String())
		}
	case suffix == "m":
		ips := strings.Split(name, ".")
		parsedIps := make([]net.IP, 0, len(ips))
		rrs = make([]dns.RR, 0)

		for i := 0; i < len(ips)-1; i++ {
			var ip net.IP
			if question.Qtype == dns.TypeAAAA && ips[i+1] == "6" {
				ip = parseIp(dns.TypeAAAA, ips[i])
			} else if question.Qtype == dns.TypeA && ips[i+1] == "4" {
				ip = parseIp(dns.TypeA, ips[i])
			}

			if ip != nil {
				rrs = append(rrs, createIpMsg(question, ip))
				parsedIps = append(parsedIps, ip)
			}
		}

		if len(parsedIps) > 0 {
			l.Info("cooking response",
				"mode", "multi", "ip", parsedIps)
		}
	default:
		ip := defaultIp(question.Qtype)
		rrs = []dns.RR{createIpMsg(question, ip)}
		l.Info("cooking response",
			"mode", "default", "ip", ip.String())
	}

	return
}

func createIpMsg(question dns.Question, ip net.IP) (rr dns.RR) {
	head := dns.RR_Header{
		Name:   question.Name,
		Rrtype: question.Qtype,
		Class:  dns.ClassINET,
		Ttl:    0,
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
