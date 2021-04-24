package handlers

import (
	log "github.com/buglloc/simplelog"
	"github.com/miekg/dns"
)

var _ Handler = DefaultHandler

func DefaultHandler(question dns.Question, name string, l *log.Logger) (rrs []dns.RR) {
	ip := defaultIp(question.Qtype)
	rrs = createIpsRR(question, ip)
	l.Info("cooking response", "mode", "default", "ip", ip.String())
	return
}
