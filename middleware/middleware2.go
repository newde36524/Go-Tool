package middleware

type Middleware2 struct {
	funcs []func(next func())
	index int
}

func (m *Middleware2) Use(f func(next func())) {
	m.funcs = append(m.funcs, f)
}

func (m *Middleware2) Invoke() {
	if m.index < len(m.funcs) {
		m.index++
		m.funcs[m.index-1](m.Invoke)
	}
}
