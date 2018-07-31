package handlers

import (
	"github.com/buglloc/simplelog"
	"github.com/miekg/dns"

	"github.com/buglloc/rip/pkg/cfg"
	"github.com/buglloc/rip/pkg/resolver"
)

func ProxyHandler(question dns.Question, name string, l log.Logger) (rrs []dns.RR) {
	if !cfg.AllowProxy {
		return
	}

	subName := parseSubName(name)
	ip, err := resolver.ResolveIp(question.Qtype, subName)
	if err != nil {
		l.Error("failed to resolve proxied name", "target", subName, "err", err.Error())
		return
	}

	rrs = createIpsRR(question, ip)
	l.Info("cooking response", "mode", "proxy", "target", subName, "ip", ip.String())
	return
}
