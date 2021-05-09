package cname

import (
	"strings"

	"github.com/miekg/dns"

	"github.com/buglloc/rip/v2/pkg/cfg"
	"github.com/buglloc/rip/v2/pkg/handlers"
	"github.com/buglloc/rip/v2/pkg/handlers/limiter"
	"github.com/buglloc/rip/v2/pkg/handlers/slices"
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
	parts, _ := p.RestValues()
	if len(parts) == 0 {
		return handlers.ErrUnexpectedEOF
	}

	slices.StringsReverse(parts)
	h.TargetFQDN = dns.Fqdn(strings.Join(parts, "."))
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
