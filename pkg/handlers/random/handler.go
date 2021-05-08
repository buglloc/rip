package random

import (
	"fmt"
	"math/rand"

	"github.com/miekg/dns"

	"github.com/buglloc/rip/v2/pkg/handlers"
	"github.com/buglloc/rip/v2/pkg/handlers/limiter"
)

const ShortName = "r"
const Name = "random"

var _ handlers.Handler = (*Handler)(nil)

type Handler struct {
	handlers.BaseHandler
	Nested [2]handlers.Handler
	Cur    int
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
	var err error
	h.Nested[0], err = p.Next()
	if err != nil {
		return err
	}

	h.Nested[1], err = p.Next()
	if err != nil {
		return err
	}

	h.Cur = randCur()
	return nil
}

func (h *Handler) Handle(question dns.Question) ([]dns.RR, bool, error) {
	rr, moveOn, err := h.Nested[h.Cur].Handle(question)
	if err != nil {
		return nil, false, fmt.Errorf("random: %w", err)
	}

	if moveOn {
		h.Cur = randCur()
	}

	if len(rr) > 0 {
		h.Limiters.Use()
	}

	return rr, h.Limiters.MoveOn(), nil
}

func randCur() int {
	if rand.Intn(50) > 50 {
		return 1
	}

	return 0
}
