package middleware

//Application x
type Application struct {
	components []Component
}

//NewApplication 创建一个新的实例
func NewApplication() *Application {
	return new(Application)
}

//New 新实例
func (a *Application) New(app *Application) *Application {
	return &Application{components: app.components}
}

//Middleware 中间件
type Middleware func(o interface{})

//Component 组件
type Component func(middle Middleware) Middleware

//Use 使用中间件
func (app *Application) Use(middleware Component) {
	app.components = append(app.components, middleware)
}

//Build 创建中间件
func (app *Application) Build() Middleware {
	var middleware Middleware = func(o interface{}) {

	}
	for _, m := range revertComponent(app.components) {
		middleware = m(middleware)
	}
	return middleware
}
func revertComponent(arr []Component) []Component {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
