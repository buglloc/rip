package limiter

import "time"

var _ Limiter = (*Count)(nil)

type TTL struct {
	TTL   time.Duration
	Until time.Time
}

func (m *TTL) Use() {
	if m.Until.Second() == 0 {
		m.Until = time.Now().Add(m.TTL)
	}
}

func (m *TTL) MoveOn() bool {
	if time.Now().After(m.Until) {
		m.Until = time.Unix(0, 0)
		return true
	}

	return false
}
