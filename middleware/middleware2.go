package middleware

//Middleware2 中间件
type Middleware2 struct {
	funcs []func(next func())
	index int
}

//Use 使用中间件
func (m *Middleware2) Use(f func(next func())) {
	m.funcs = append(m.funcs, f)
}

//Use 调用中间件
func (m *Middleware2) Invoke() {
	if m.index < len(m.funcs) {
		m.index++
		m.funcs[m.index-1](m.Invoke)
	}
}
