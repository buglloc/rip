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
	"github.com/buglloc/rip/v2/pkg/handlers/notify"
	"github.com/buglloc/rip/v2/pkg/handlers/proxy"
	"github.com/buglloc/rip/v2/pkg/handlers/random"
	"github.com/buglloc/rip/v2/pkg/handlers/slices"
	"github.com/buglloc/rip/v2/pkg/handlers/sticky"
)

var _ handlers.Parser = (*Parser)(nil)

type Parser struct {
	cur      int
	maxLabel int
	labels   []string
	fqdn     string
}

func NewParser(fqdn, zone string) *Parser {
	ripReq := fqdn
	if len(zone) > 0 {
		ripReq = fqdn[:len(fqdn)-len(zone)-1]
	}

	labels := strings.Split(ripReq, ".")
	slices.StringsReverse(labels)
	return &Parser{
		cur:      0,
		maxLabel: len(labels),
		labels:   labels,
		fqdn:     fqdn,
	}
}

func (p *Parser) FQDN() string {
	return p.fqdn
}

func (p *Parser) NextHandler() (handlers.Handler, error) {
	if p.cur >= p.maxLabel {
		return nil, handlers.ErrEOF
	}

	part := p.labels[p.cur]
	h := parseHandler(part)
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
	if p.cur >= p.maxLabel {
		return "", handlers.ErrEOF
	}

	ret := p.labels[p.cur]
	if strings.HasPrefix(ret, "b32-") {
		decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(ret[4:]))
		if err == nil {
			ret = string(decoded)
		}
	}
	p.cur++
	return ret, nil
}

func (p *Parser) RestValues() ([]string, error) {
	var out []string
	for {
		v, err := p.NextValue()
		if v == "" {
			return out, nil
		}

		if err != nil {
			return nil, err
		}

		out = append(out, v)
	}
}

func (p *Parser) RestHandlers() ([]handlers.Handler, error) {
	var out []handlers.Handler
	for {
		h, err := p.NextHandler()
		if h == nil {
			return out, nil
		}

		if err != nil {
			return nil, err
		}

		out = append(out, h)
	}
}

func (p *Parser) NextValue() (string, error) {
	if p.cur >= p.maxLabel {
		return "", handlers.ErrEOF
	}

	label := p.labels[p.cur]
	handlerName := label
	if indx := strings.IndexByte(handlerName, '-'); indx > 0 {
		handlerName = handlerName[:indx]
	}

	switch handlerName {
	case ipv4.ShortName, ipv4.Name:
		fallthrough
	case ipv6.ShortName, ipv6.Name:
		fallthrough
	case cname.ShortName, cname.Name:
		fallthrough
	case proxy.ShortName, proxy.Name:
		fallthrough
	case random.ShortName, random.Name:
		fallthrough
	case loop.ShortName, loop.Name:
		fallthrough
	case sticky.ShortName, sticky.Name:
		fallthrough
	case notify.ShortName, notify.Name:
		fallthrough
	case defaultip.ShortName, defaultip.Name:
		return "", handlers.ErrEOF
	}

	if strings.HasPrefix(label, "b32-") {
		decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(label[4:]))
		if err == nil {
			label = string(decoded)
		}
	}
	p.cur++
	return label, nil
}

func parseHandler(label string) handlers.Handler {
	if len(label) == 0 {
		return nil
	}

	parts := strings.Split(label, "-")
	parseLimiters := func() []limiter.Limiter {
		if len(parts) <= 1 {
			return nil
		}

		opts := make(map[string]string, len(parts)/2)
		for i := 1; i < len(parts)-1; i += 2 {
			opts[parts[i]] = parts[i+1]
		}

		ret, err := limiter.ParseLimiters(opts)
		if err != nil {
			log.Error("can't parse limiter", "label", label, "err", err)
			return nil
		}

		return ret
	}

	switch parts[0] {
	case ipv4.ShortName, ipv4.Name:
		return ipv4.NewHandler(parseLimiters()...)
	case ipv6.ShortName, ipv6.Name:
		return ipv6.NewHandler(parseLimiters()...)
	case cname.ShortName, cname.Name:
		return cname.NewHandler(parseLimiters()...)
	case proxy.ShortName, proxy.Name:
		return proxy.NewHandler(parseLimiters()...)
	case random.ShortName, random.Name:
		return random.NewHandler(parseLimiters()...)
	case loop.ShortName, loop.Name:
		return loop.NewHandler(parseLimiters()...)
	case sticky.ShortName, sticky.Name:
		return sticky.NewHandler(parseLimiters()...)
	case notify.ShortName, notify.Name:
		return notify.NewHandler()
	case defaultip.ShortName, defaultip.Name:
		return defaultip.NewHandler(parseLimiters()...)
	default:
		return parseIPHandler(label)
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
