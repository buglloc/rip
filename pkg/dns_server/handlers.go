package dns_server

import (
	"github.com/buglloc/rip/pkg/cfg"
	"github.com/buglloc/simplelog"
	"github.com/miekg/dns"
	"net"
	"strings"
)

func NewHandler(zone string) func(w dns.ResponseWriter, req *dns.Msg) {
	return func(w dns.ResponseWriter, req *dns.Msg) {
		defer w.Close()
		msg := handle(zone, req)
		if msg == nil {
			msg = &dns.Msg{}
			msg.SetRcode(req, dns.RcodeServerFailure)
		}
		w.WriteMsg(msg)
	}
}

func handle(zone string, req *dns.Msg) *dns.Msg {
	response := &dns.Msg{}
	response.SetReply(req)
	for _, question := range req.Question {
		if question.Qtype != dns.TypeA && question.Qtype != dns.TypeAAAA {
			log.Debug("skip unknown request", "type", typeToString(question.Qtype))
			return nil
		}

		ip, err := ipFromName(question, zone)
		if err != nil {
			log.Error("failed to parse request", "type", typeToString(question.Qtype), "name", question.Name, "err", err.Error())
			continue
		}

		head := dns.RR_Header{
			Name:   question.Name,
			Rrtype: question.Qtype,
			Class:  dns.ClassINET,
			Ttl:    1,
		}

		var line dns.RR
		if question.Qtype == dns.TypeA {
			line = &dns.A{
				Hdr: head,
				A:   ip,
			}
		} else {
			line = &dns.AAAA{
				Hdr:  head,
				AAAA: ip,
			}
		}
		response.Answer = append(response.Answer, line)
	}

	return response
}

func ipFromName(question dns.Question, zone string) (ip net.IP, err error) {
	name := question.Name
	name = name[:len(name)-len(zone)-1]
	if len(name) <= 2 {
		ip = defaultIp(question.Qtype)
		return
	}

	i := len(name) - 2
	name, suffix := name[:i], name[i:]
	switch {
	case suffix == ".p":
		ip, err = ResolveIp(question.Qtype, name)
	case suffix == ".4" && question.Qtype == dns.TypeA:
		ip = parseIp(dns.TypeA, name)
	case suffix == ".6" && question.Qtype == dns.TypeAAAA:
		ip = parseIp(dns.TypeAAAA, name)
	default:
		ip = defaultIp(question.Qtype)
	}

	return
}

func typeToString(reqType uint16) string {
	if t, ok := dns.TypeToString[reqType]; ok {
		return t
	}
	return "unknown"
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
		if dotCounts != 4 {
			return defaultIp(dns.TypeA)
		}
		return net.ParseIP(strings.Replace(name, "-", ".", -1))
	case dns.TypeAAAA:
		if dotCounts < 2 {
			return defaultIp(dns.TypeAAAA)
		}
		return net.ParseIP(strings.Replace(name, "-", ":", -1))
	default:
		return defaultIp(dns.TypeA)
	}
}
