package tcpserver

//Logger 日志接口
type Logger interface {
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Trace(v ...interface{})
	Tracef(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Critical(v ...interface{})
	Criticalf(format string, v ...interface{})
	All(v ...interface{})
	Allf(format string, v ...interface{})
	Fatal(code int, v ...interface{})
	Fatalf(code int, format string, v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}
