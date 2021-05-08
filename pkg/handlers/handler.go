package handlers

import (
	"github.com/miekg/dns"

	"github.com/buglloc/rip/v2/pkg/handlers/limiter"
)

type Handler interface {
	Name() string
	Init(p Parser) error
	SetDefaultLimiters(modifiers ...limiter.Limiter)
	Handle(question dns.Question) (rrs []dns.RR, moveOn bool, err error)
}

type BaseHandler struct {
	Limiters limiter.Limiters
}

func (h *BaseHandler) SetDefaultLimiters(modifiers ...limiter.Limiter) {
	if len(h.Limiters) == 0 {
		h.Limiters = modifiers
	}
}
