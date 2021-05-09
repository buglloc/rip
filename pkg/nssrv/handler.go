package nssrv

import (
	"fmt"
	"sync"

	log "github.com/buglloc/simplelog"
	"github.com/miekg/dns"

	"github.com/buglloc/rip/v2/pkg/cfg"
	"github.com/buglloc/rip/v2/pkg/handlers"
	"github.com/buglloc/rip/v2/pkg/handlers/defaultip"
	"github.com/buglloc/rip/v2/pkg/handlers/parser"
)

var defaultHandler = &defaultip.Handler{}

type cachedHandler struct {
	handlers.Handler
	mu sync.Mutex
}

func (h *cachedHandler) Handle(question dns.Question) ([]dns.RR, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	ret, _, err := h.Handler.Handle(question)
	return ret, err
}

func (s *NSSrv) handleRequest(zone string, req *dns.Msg, logger *log.Logger) *dns.Msg {
	out := &dns.Msg{}
	out.SetReply(req)

	realHandler := func(question dns.Question, zone string) (*cachedHandler, error) {
		if len(question.Name)-len(zone) <= 3 {
			// fast exit
			return &cachedHandler{
				Handler: defaultHandler,
			}, nil
		}

		//if item != nil &&
		item := s.cache.Get(question.Name)
		if item != nil {
			if item.Expired() {
				item.Extend(cfg.CacheTTL)
			}
			return item.Value().(*cachedHandler), nil
		}

		h, err := parser.NewParser(question.Name, zone).Next()
		if err != nil {
			if err != handlers.ErrEOF {
				return nil, err
			}

			if !cfg.UseDefault {
				return nil, fmt.Errorf("no handlers for request %q available", question.Name)
			}

			h = defaultHandler
		}

		if err != nil {
			return nil, err
		}

		ret := &cachedHandler{
			Handler: h,
		}
		s.cache.Set(question.Name, ret, cfg.CacheTTL)
		return ret, nil
	}

	for _, question := range req.Question {
		switch question.Qtype {
		case dns.TypeA, dns.TypeAAAA:
			l := logger.Child("type", dns.Type(question.Qtype), "name", question.Name)
			handler, err := realHandler(question, zone)
			if err != nil {
				l.Error("failed to parse request", "err", err)
				continue
			}

			answers, err := handler.Handle(question)
			if err != nil {
				l.Error("failed to handle request", "err", err)
				continue
			}

			l.Info("cooking response", "answers", fmt.Sprint(answers))
			if len(answers) == 0 {
				continue
			}

			out.Answer = append(out.Answer, answers...)
		}
	}

	return out
}
