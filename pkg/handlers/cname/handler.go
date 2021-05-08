package cname

import (
	"github.com/miekg/dns"

	"github.com/buglloc/rip/v2/pkg/cfg"
	"github.com/buglloc/rip/v2/pkg/handlers"
	"github.com/buglloc/rip/v2/pkg/handlers/limiter"
)

const ShortName = "c"
const Name = "cname"

var _ handlers.Handler = (*Handler)(nil)

type Handler struct {
	handlers.BaseHandler
	TargetFQDN string
}

func NewHandler(modifiers ...limiter.Limiter) *Handler {
	return &Handler{
		BaseHandler: handlers.BaseHandler{
			Limiters: modifiers,
		},
	}
}

func (h *Handler) Name() string {
	return Name
}

func (h *Handler) Init(p handlers.Parser) error {
	name, _ := p.NextRaw()
	if name == "" {
		return handlers.ErrUnexpectedEOF
	}

	h.TargetFQDN = handlers.PartToFQDN(name)
	return nil
}

func (h *Handler) Handle(question dns.Question) ([]dns.RR, bool, error) {
	rr := []dns.RR{&dns.CNAME{
		Hdr: dns.RR_Header{
			Name:   question.Name,
			Rrtype: dns.TypeCNAME,
			Class:  dns.ClassINET,
			Ttl:    cfg.TTL,
		},
		Target: h.TargetFQDN,
	}}

	h.Limiters.Use()
	return rr, h.Limiters.MoveOn(), nil
}
