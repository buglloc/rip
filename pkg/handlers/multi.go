package handlers

import (
	"net"
	"strings"

	log "github.com/buglloc/simplelog"
	"github.com/miekg/dns"
)

func MultiHandler(question dns.Question, name string, l log.Logger) (rrs []dns.RR) {
	ips := strings.Split(name, ".")
	parsedIps := make([]net.IP, 0, len(ips))
	for i := 0; i < len(ips)-1; i++ {
		var ip net.IP
		if question.Qtype == dns.TypeAAAA && ips[i+1] == "6" {
			ip = parseIp(dns.TypeAAAA, ips[i])
		} else if question.Qtype == dns.TypeA && ips[i+1] == "4" {
			ip = parseIp(dns.TypeA, ips[i])
		}

		if ip != nil {
			parsedIps = append(parsedIps, ip)
		}
	}

	if len(parsedIps) == 0 {
		return
	}

	rrs = createIpsRR(question, parsedIps...)
	l.Info("cooking response", "mode", "multi", "ip", parsedIps)
	return
}
