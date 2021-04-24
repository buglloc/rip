package handlers

import (
	"strings"

	log "github.com/buglloc/simplelog"
	"github.com/miekg/dns"

	"github.com/buglloc/rip/pkg/ip_loop"
)

func LoopHandler(question dns.Question, name string, l log.Logger) (rrs []dns.RR) {
	ips := strings.Split(name, ".")
	if len(ips) < 2 {
		log.Error("failed to parse loop annotation")
		return
	}

	ips = ips[len(ips)-2:]
	var pIp string
	if question.Qtype == dns.TypeA {
		// Move next only for ipv4 request
		// Maybe we can make something better?
		pIp = ip_loop.GetNext(name, ips)
	} else {
		pIp = ip_loop.GetCurrent(name, ips)
	}

	ip := parseIp(question.Qtype, pIp)
	if ip == nil {
		return
	}

	rrs = createIpsRR(question, ip)
	l.Info("cooking response", "mode", "loop", "ip", ip.String())
	return
}
