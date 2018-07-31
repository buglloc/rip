package handlers

import (
	"strings"

	"github.com/buglloc/simplelog"
	"github.com/miekg/dns"
)

func RandomHandler(question dns.Question, name string, l log.Logger) (rrs []dns.RR) {
	ips := strings.Split(name, ".")
	if len(ips) < 2 {
		l.Error("failed to parse random annotation")
		return
	}

	var currentIp string
	if random(0, 100) > 50 {
		currentIp = ips[len(ips)-2]
	} else {
		currentIp = ips[len(ips)-1]
	}

	ip := parseIp(question.Qtype, currentIp)
	if ip == nil {
		return
	}

	rrs = createIpsRR(question, ip)
	l.Info("cooking response", "mode", "random", "ip", ip.String())
	return
}
