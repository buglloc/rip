package notify

import (
	"strings"
	"time"

	"github.com/miekg/dns"

	"github.com/buglloc/rip/v2/pkg/cfg"
	"github.com/buglloc/rip/v2/pkg/handlers"
	"github.com/buglloc/rip/v2/pkg/hub"
)

const ShortName = "n"
const Name = "notify"

var _ handlers.Handler = (*Handler)(nil)

type Handler struct {
	handlers.BaseHandler
	channel string
	Nested  handlers.Handler
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Name() string {
	return Name
}

func (h *Handler) Init(p handlers.Parser) error {
	var err error
	h.channel, err = p.NextRaw()
	if err != nil {
		return err
	}

	h.Nested, err = p.Next()
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Handle(question dns.Question) ([]dns.RR, bool, error) {
	rr, moveOn, err := h.Nested.Handle(question)
	if cfg.HubEnabled {
		now := time.Now()
		for _, r := range rr {
			hub.Send(h.channel, hub.Message{
				Time:  now,
				Name:  question.Name,
				QType: dns.Type(question.Qtype).String(),
				RR:    strings.TrimPrefix(r.String(), r.Header().String()),
			})
		}
	}

	return rr, moveOn, err
}
