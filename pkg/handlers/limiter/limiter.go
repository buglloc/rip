package limiter

import (
	"fmt"
	"strconv"
	"time"
)

type Limiter interface {
	Use()
	MoveOn() bool
}

type Limiters []Limiter

func (m Limiters) Use() {
	for _, mod := range m {
		mod.Use()
	}
}

func (m Limiters) MoveOn() bool {
	if len(m) == 0 {
		return true
	}

	moveOn := false
	for _, mod := range m {
		if mod.MoveOn() {
			moveOn = true
		}
	}

	return moveOn
}

func ParseLimiters(opts map[string]string) ([]Limiter, error) {
	out := make([]Limiter, 0, len(opts))
	for k, v := range opts {
		switch k {
		case "ttl":
			ttl, err := time.ParseDuration(v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse 'ttl' limiter (%q): %w", v, err)
			}
			out = append(out, &TTL{
				TTL: ttl,
			})
		case "cnt":
			cnt, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse 'cnt' limiter (%q): %w", v, err)
			}

			out = append(out, &Count{
				Max: cnt,
			})
		default:
			return nil, fmt.Errorf("unexpected limiter: %s", k)
		}
	}

	return out, nil
}
