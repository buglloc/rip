package defaultip

import (
	"github.com/miekg/dns"

	"github.com/buglloc/rip/v2/pkg/handlers"
	"github.com/buglloc/rip/v2/pkg/handlers/limiter"
)

const ShortName = "d"
const Name = "default"

var _ handlers.Handler = (*Handler)(nil)

type Handler struct {
	handlers.BaseHandler
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

func (h *Handler) Init(_ handlers.Parser) error {
	return nil
}

func (h *Handler) Handle(question dns.Question) ([]dns.RR, bool, error) {
	h.Limiters.Use()
	return handlers.IPsToRR(question, handlers.DefaultIp(question.Qtype)), h.Limiters.MoveOn(), nil
}
