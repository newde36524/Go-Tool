package caller

import (
	"fmt"
	"runtime"
)

//GetCallerName 获取当天调用者方法名
func GetCallerName() string {
	pc, _, _, _ := runtime.Caller(1)
	fun := runtime.FuncForPC(pc)
	return fun.Name()
}

//GetCallerNames 获取多个调用者的方法名,获取调用链
func GetCallerNames() {
	var ps []uintptr = make([]uintptr, 100)
	i := runtime.Callers(0, ps)
	ps = ps[:i]
	for _, p := range ps {
		fmt.Println(runtime.FuncForPC(p).Name())
	}
}
