package handlers

import (
	"strings"

	"github.com/buglloc/simplelog"
	"github.com/miekg/dns"

	"github.com/buglloc/rip/pkg/cfg"
)

type Handler func(question dns.Question, zone string, l log.Logger) (rrs []dns.RR)

var handlers = map[string]Handler{
	"c": CnameHandler,
	"p": ProxyHandler,
	"l": LoopHandler,
	"s": StickyHandler,
	"r": RandomHandler,
	"4": Ipv4Handler,
	"6": Ipv6Handler,
	"m": MultiHandler,
}

func Handle(question dns.Question, zone string, l log.Logger) (rrs []dns.RR, err error) {
	if len(question.Name)-len(zone) <= 3 {
		// Fast exit
		rrs = DefaultHandler(question, question.Name, l)
		return
	}

	var name, mode string
	name = question.Name[:len(question.Name)-len(zone)-1]
	i := strings.LastIndex(name, ".")
	if i > 0 && i < len(name) {
		name, mode = name[:i], name[i+1:]
	}

	h, ok := handlers[mode]
	if !ok {
		h = DefaultHandler
	}

	rrs = h(question, name, l)
	if rrs == nil || len(rrs) == 0 {
		if cfg.UseDefault {
			rrs = DefaultHandler(question, name, l)
		} else {
			l.Info("failed to handle request")
		}
	}
	return
}
