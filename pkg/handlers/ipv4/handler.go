package ipv4

import (
	"fmt"
	"net"

	"github.com/miekg/dns"

	"github.com/buglloc/rip/v2/pkg/handlers"
	"github.com/buglloc/rip/v2/pkg/handlers/limiter"
)

const ShortName = "4"
const Name = "v4"

var _ handlers.Handler = (*Handler)(nil)

type Handler struct {
	handlers.BaseHandler
	IP net.IP
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
	if len(h.IP) > 0 {
		return nil
	}

	ip, _ := p.NextRaw()
	if ip == "" {
		return handlers.ErrUnexpectedEOF
	}

	targetIP := handlers.PartToIP(ip)
	if len(targetIP) != net.IPv4len {
		return fmt.Errorf("not IPv4 address: %s", ip)
	}

	h.IP = targetIP
	return nil
}

func (h *Handler) Handle(question dns.Question) ([]dns.RR, bool, error) {
	if question.Qtype != dns.TypeA {
		return nil, false, nil
	}

	h.Limiters.Use()
	return handlers.IPsToRR(question, h.IP), h.Limiters.MoveOn(), nil
}
