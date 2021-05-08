package proxy

import (
	"fmt"

	"github.com/miekg/dns"

	"github.com/buglloc/rip/v2/pkg/cfg"
	"github.com/buglloc/rip/v2/pkg/handlers"
	"github.com/buglloc/rip/v2/pkg/handlers/limiter"
	"github.com/buglloc/rip/v2/pkg/resolver"
)

const ShortName = "p"
const Name = "proxy"

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
	if !cfg.AllowProxy {
		return handlers.ErrNotAllowed
	}

	name, _ := p.NextRaw()
	if name == "" {
		return handlers.ErrUnexpectedEOF
	}

	h.TargetFQDN = handlers.PartToFQDN(name)
	return nil
}

func (h *Handler) Handle(question dns.Question) ([]dns.RR, bool, error) {
	ips, err := resolver.ResolveIp(question.Qtype, h.TargetFQDN)
	if err != nil {
		return nil, false, fmt.Errorf("proxy: %w", err)
	}

	h.Limiters.Use()
	return handlers.IPsToRR(question, ips...), h.Limiters.MoveOn(), nil
}
