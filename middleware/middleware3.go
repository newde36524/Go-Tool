package middleware

//Middleware3 中间件
type Middleware3 struct {
	funcs   []func(context interface{}, next func())
	context interface{}
}

//Use 使用中间件
//@f 中间件方法
func (m *Middleware3) Use(f func(context interface{}, next func())) {
	m.funcs = append(m.funcs, f)
}

//Invoke 调用中间件
func (m *Middleware3) Invoke(context interface{}) {
	index := 0
	m.context = context
	var next func()
	next = func() {
		if index < len(m.funcs) {
			index++
			m.funcs[index-1](m.context, next)
		}
	}
	next()
}
