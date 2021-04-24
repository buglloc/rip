package handlers

import (
	log "github.com/buglloc/simplelog"
	"github.com/miekg/dns"

	"github.com/buglloc/rip/pkg/cfg"
)

func CnameHandler(question dns.Question, name string, l log.Logger) (rrs []dns.RR) {
	subName := parseSubName(name)
	rrs = []dns.RR{&dns.CNAME{
		Hdr: dns.RR_Header{
			Name:   question.Name,
			Rrtype: dns.TypeCNAME,
			Class:  dns.ClassINET,
			Ttl:    cfg.TTL,
		},
		Target: dns.Fqdn(subName),
	}}
	l.Info("cooking response", "mode", "cname", "target", subName)
	return
}
