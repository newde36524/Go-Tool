package middleware

//Do 以相同的处理逻辑和next 处理切片中每一个元素
func Do(pointers []interface{}, fn func(interface{}, func())) {
	index := 0
	var next func()
	next = func() {
		if index < len(pointers) {
			index++
			fn(pointers[index-1], next)
		}
	}
	next()
}
