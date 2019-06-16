package middleware

import (
	"container/list"
)

//Middleware5 .
type Middleware5 struct {
	node list.List
}

//MiddleFunc .
type MiddleFunc func(o interface{}, next func())

//Use .
func (m *Middleware5) Use(f MiddleFunc) {
	m.node.PushBack(f)
}

//Invoke .
func (m *Middleware5) Invoke(o interface{}) {
	curr := m.node.Front()
	var fn func()
	fn = func() {
		if curr != nil {
			v := curr.Value
			switch f := v.(type) {
			case MiddleFunc:
				curr = curr.Next()
				f(o, fn)
			}
		}
	}
	fn()
}
