package parser

import (
	"encoding/base32"
	"fmt"
	"net"
	"strings"

	log "github.com/buglloc/simplelog"

	"github.com/buglloc/rip/v2/pkg/handlers"
	"github.com/buglloc/rip/v2/pkg/handlers/cname"
	"github.com/buglloc/rip/v2/pkg/handlers/defaultip"
	"github.com/buglloc/rip/v2/pkg/handlers/ipv4"
	"github.com/buglloc/rip/v2/pkg/handlers/ipv6"
	"github.com/buglloc/rip/v2/pkg/handlers/limiter"
	"github.com/buglloc/rip/v2/pkg/handlers/loop"
	"github.com/buglloc/rip/v2/pkg/handlers/proxy"
	"github.com/buglloc/rip/v2/pkg/handlers/random"
	"github.com/buglloc/rip/v2/pkg/handlers/sticky"
)

var _ handlers.Parser = (*Parser)(nil)

type Parser struct {
	cur     int
	maxPart int
	parts   []string
}

func NewParser(req string) *Parser {
	parts := strings.Split(req, ".")
	reverse(parts)
	return &Parser{
		cur:     0,
		maxPart: len(parts),
		parts:   parts,
	}
}

func (p *Parser) Next() (handlers.Handler, error) {
	if p.cur >= p.maxPart {
		return nil, handlers.ErrEOF
	}

	part := p.parts[p.cur]
	h := parsePart(part)
	if h == nil {
		return nil, handlers.ErrEOF
	}

	p.cur++
	err := h.Init(p)
	if err != nil {
		return nil, fmt.Errorf("can't parse handler %s: %w", h.Name(), err)
	}

	return h, nil
}

func (p *Parser) NextRaw() (string, error) {
	if p.cur >= p.maxPart {
		return "", handlers.ErrEOF
	}

	ret := p.parts[p.cur]
	if strings.HasPrefix(ret, "b32-") {
		decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(ret[4:]))
		if err == nil {
			ret = string(decoded)
		}
	}
	p.cur++
	return ret, nil
}

func (p *Parser) All() ([]handlers.Handler, error) {
	var out []handlers.Handler
	for {
		h, err := p.Next()
		if h == nil {
			return out, nil
		}

		if err != nil {
			return nil, err
		}

		out = append(out, h)
	}
}

func parsePart(part string) handlers.Handler {
	if len(part) == 0 {
		return nil
	}

	parts := strings.Split(part, "-")
	parseModifiers := func() []limiter.Limiter {
		if len(parts) <= 1 {
			return nil
		}

		opts := make(map[string]string, len(parts)-1/2)
		for i := 1; i < len(parts)-1; i += 2 {
			opts[parts[i]] = parts[i+1]
		}

		ret, err := limiter.ParseLimiters(opts)
		if err != nil {
			log.Error("can't parse modifiers", "part", part, "err", err)
			return nil
		}

		return ret
	}

	modifiers := parseModifiers()
	switch parts[0] {
	case defaultip.ShortName, defaultip.Name:
		return defaultip.NewHandler(modifiers...)
	case cname.ShortName, cname.Name:
		return cname.NewHandler(modifiers...)
	case proxy.ShortName, proxy.Name:
		return proxy.NewHandler(modifiers...)
	case random.ShortName, random.Name:
		return random.NewHandler(modifiers...)
	case loop.ShortName, loop.Name:
		return loop.NewHandler(modifiers...)
	case sticky.ShortName, sticky.Name:
		return sticky.NewHandler(modifiers...)
	case ipv4.ShortName, ipv4.Name:
		return ipv4.NewHandler(modifiers...)
	case ipv6.ShortName, ipv6.Name:
		return ipv6.NewHandler(modifiers...)
	default:
		return parseIPHandler(part)
	}
}

func parseIPHandler(part string) handlers.Handler {
	ip := handlers.PartToIP(part)
	switch len(ip) {
	case net.IPv4len:
		return &ipv4.Handler{IP: ip}
	case net.IPv6len:
		return &ipv6.Handler{IP: ip}
	default:
		return nil
	}
}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}
