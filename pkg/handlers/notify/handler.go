package notify

import (
	"fmt"
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
	init := func() error {
		var err error
		h.channel, err = p.NextRaw()
		if err != nil {
			return err
		}

		h.Nested, err = p.NextHandler()
		if err != nil {
			return err
		}

		return nil
	}

	err := init()
	if err != nil {
		h.reportErr(dns.Question{Name: p.FQDN()}, fmt.Sprintf("can't parse request: %v", err))
	}
	return err
}

func (h *Handler) Handle(question dns.Question) ([]dns.RR, bool, error) {
	rr, moveOn, err := h.Nested.Handle(question)
	if cfg.HubEnabled {
		if err != nil {
			h.reportErr(question, err.Error())
		} else {
			h.reportRR(question, rr)
		}
	}

	return rr, moveOn, err
}

func (h *Handler) reportRR(question dns.Question, rr []dns.RR) {
	if h.channel == "" {
		return
	}

	now := time.Now()
	if len(rr) == 0 {
		hub.Send(h.channel, hub.Message{
			Time:  now,
			Name:  question.Name,
			QType: dns.Type(question.Qtype).String(),
			RR:    "<empty>",
			Ok:    true,
		})
		return
	}

	for _, r := range rr {
		hub.Send(h.channel, hub.Message{
			Time:  now,
			Name:  question.Name,
			QType: dns.Type(question.Qtype).String(),
			RR:    strings.TrimPrefix(r.String(), r.Header().String()),
			Ok:    true,
		})
	}
}

func (h *Handler) reportErr(question dns.Question, err string) {
	if h.channel == "" {
		return
	}

	qType := "n/a"
	if question.Qtype != dns.TypeNone {
		qType = dns.Type(question.Qtype).String()
	}

	hub.Send(h.channel, hub.Message{
		Time:  time.Now(),
		Name:  question.Name,
		QType: qType,
		RR:    err,
		Ok:    false,
	})
}
