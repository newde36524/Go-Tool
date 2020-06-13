package middleware

import "context"

type (
	//Pipe .
	Pipe interface {
		Regist(h interface{}) Pipe
		schedule(fn func(h interface{}, ctx interface{}), ctx interface{})
	}

	//pipeLine .
	pipeLine struct {
		ctx     context.Context
		handles []interface{}
	}

	//canNext .
	canNext interface {
		setNext(next func())
	}
)

//Regist .
func (p *pipeLine) Regist(h interface{}) Pipe {
	p.handles = append(p.handles, h)
	return p
}

//schedule pipeline provider
func (p *pipeLine) schedule(fn func(h interface{}, ctx interface{}), ctx interface{}) {
	index := 0
	next := func() {
		if index < len(p.handles) {
			index++
			fn(p.handles[index-1], ctx)
		}
	}
	if v, ok := ctx.(canNext); ok {
		v.setNext(next)
	}
	next()
}
