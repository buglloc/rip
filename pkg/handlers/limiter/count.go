package limiter

var _ Limiter = (*Count)(nil)

type Count struct {
	Cur int
	Max int
}

func (m *Count) Use() {
	if m.Cur < 0 {
		m.Cur = 0
	}
	m.Cur++
}

func (m *Count) MoveOn() bool {
	expired := m.Cur >= m.Max
	if expired {
		m.Cur = -1
	}
	return expired
}
