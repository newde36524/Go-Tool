package httptool

type HttpMethod struct {
	method string
}

func NewHttpMethod(method string) *HttpMethod {
	return &HttpMethod{method}
}
func Delete() *HttpMethod {
	return NewHttpMethod("Delete")
}
func Get() *HttpMethod {
	return NewHttpMethod("GET")
}
func Head() *HttpMethod {
	return NewHttpMethod("Head")
}
func Options() *HttpMethod {
	return NewHttpMethod("Options")
}
func Patch() *HttpMethod {
	return NewHttpMethod("Patch")
}
func Post() *HttpMethod {
	return NewHttpMethod("Post")
}
func Put() *HttpMethod {
	return NewHttpMethod("Put")
}
func Trace() *HttpMethod {
	return NewHttpMethod("Trace")
}
