package tcpserver3

import (
	"github.com/issue9/logs"
)

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

//DefaultLogger .
type DefaultLogger struct {
	Logger
}

//NewDefaultLogger .
func NewDefaultLogger() (result DefaultLogger, err error) {
	result = DefaultLogger{}
	return
}

//Info .
func (DefaultLogger) Info(v ...interface{}) {
	logs.Info(v...)
}

//Infof .
func (DefaultLogger) Infof(format string, v ...interface{}) {
	logs.Infof(format, v...)
}

//Debug .
func (DefaultLogger) Debug(v ...interface{}) {
	logs.Debug(v...)
}

//Debugf .
func (DefaultLogger) Debugf(format string, v ...interface{}) {
	logs.Debugf(format, v...)
}

//Trace .
func (DefaultLogger) Trace(v ...interface{}) {
	logs.Trace(v...)
}

//Tracef .
func (DefaultLogger) Tracef(format string, v ...interface{}) {
	logs.Tracef(format, v...)
}

//Warn .
func (DefaultLogger) Warn(v ...interface{}) {
	logs.Warn(v...)
}

//Warnf .
func (DefaultLogger) Warnf(format string, v ...interface{}) {
	logs.Warnf(format, v...)
}

//Error .
func (DefaultLogger) Error(v ...interface{}) {
	logs.Error(v...)
}

//Errorf .
func (DefaultLogger) Errorf(format string, v ...interface{}) {
	logs.Errorf(format, v...)
}

//Critical .
func (DefaultLogger) Critical(v ...interface{}) {
	logs.Critical(v...)
}

//Criticalf .
func (DefaultLogger) Criticalf(format string, v ...interface{}) {
	logs.Criticalf(format, v...)
}

//All .
func (DefaultLogger) All(v ...interface{}) {
	logs.All(v...)
}

//Allf .
func (DefaultLogger) Allf(format string, v ...interface{}) {
	logs.Allf(format, v...)
}

//Fatal .
func (DefaultLogger) Fatal(code int, v ...interface{}) {
	logs.Fatal(code, v...)
}

//Fatalf .
func (DefaultLogger) Fatalf(code int, format string, v ...interface{}) {
	logs.Fatalf(code, format, v...)
}

//Panic .
func (DefaultLogger) Panic(v ...interface{}) {
	logs.Panic(v...)
}

//Panicf .
func (DefaultLogger) Panicf(format string, v ...interface{}) {
	logs.Panicf(format, v...)
}
