package handlers

import (
	"strings"

	"github.com/buglloc/simplelog"
	"github.com/miekg/dns"

	"github.com/buglloc/rip/pkg/ip_stick"
)

func StickyHandler(question dns.Question, name string, l log.Logger) (rrs []dns.RR) {
	ips := strings.Split(name, ".")
	if len(ips) < 2 {
		log.Error("failed to parse loop annotation")
		return
	}

	var key string
	if question.Qtype == dns.TypeA {
		key = name + "A"
	} else {
		key = name + "AAAA"
	}

	ips = ips[len(ips)-2:]
	targetIp := ip_stick.GetCurrent(key, ips)
	ip := parseIp(question.Qtype, targetIp)
	if ip == nil {
		return
	}

	rrs = createIpsRR(question, ip)
	l.Info("cooking response", "mode", "sticky", "ip", ip.String())
	return
}
