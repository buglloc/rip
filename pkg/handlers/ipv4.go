package handlers

import (
	"github.com/buglloc/simplelog"
	"github.com/miekg/dns"
)

func Ipv4Handler(question dns.Question, name string, l log.Logger) (rrs []dns.RR) {
	if question.Qtype != dns.TypeA {
		return
	}

	ip := parseIp(dns.TypeA, name)
	if ip == nil {
		return
	}

	rrs = createIpsRR(question, ip)
	l.Info("cooking response", "mode", "ipv4", "ip", ip.String())
	return
}
