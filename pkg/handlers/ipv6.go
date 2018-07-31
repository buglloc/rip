package handlers

import (
	"github.com/buglloc/simplelog"
	"github.com/miekg/dns"
)

func Ipv6Handler(question dns.Question, name string, l log.Logger) (rrs []dns.RR) {
	if question.Qtype != dns.TypeAAAA {
		return
	}

	ip := parseIp(dns.TypeAAAA, name)
	if ip == nil {
		return
	}

	rrs = createIpsRR(question, ip)
	l.Info("cooking response", "mode", "ipv6", "ip", ip.String())
	return
}
