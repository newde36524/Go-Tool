package tcpserver

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

type DefaultLogger struct {
	Logger
}

func NewDefaultLogger() (result DefaultLogger, err error) {
	result = DefaultLogger{}
	return
}

func (DefaultLogger) Info(v ...interface{}) {
	logs.Info(v)
}
func (DefaultLogger) Infof(format string, v ...interface{}) {
	logs.Infof(format, v)
}
func (DefaultLogger) Debug(v ...interface{}) {
	logs.Debug(v)
}
func (DefaultLogger) Debugf(format string, v ...interface{}) {
	logs.Debugf(format, v)
}
func (DefaultLogger) Trace(v ...interface{}) {
	logs.Trace(v)
}
func (DefaultLogger) Tracef(format string, v ...interface{}) {
	logs.Tracef(format, v)
}
func (DefaultLogger) Warn(v ...interface{}) {
	logs.Warn(v)
}
func (DefaultLogger) Warnf(format string, v ...interface{}) {
	logs.Warnf(format, v)
}
func (DefaultLogger) Error(v ...interface{}) {
	logs.Error(v)
}
func (DefaultLogger) Errorf(format string, v ...interface{}) {
	logs.Errorf(format, v)
}
func (DefaultLogger) Critical(v ...interface{}) {
	logs.Critical(v)
}
func (DefaultLogger) Criticalf(format string, v ...interface{}) {
	logs.Criticalf(format, v)
}
func (DefaultLogger) All(v ...interface{}) {
	logs.All(v)
}
func (DefaultLogger) Allf(format string, v ...interface{}) {
	logs.Allf(format, v)
}
func (DefaultLogger) Fatal(code int, v ...interface{}) {
	logs.Fatal(code, v)
}
func (DefaultLogger) Fatalf(code int, format string, v ...interface{}) {
	logs.Fatalf(code, format, v)
}
func (DefaultLogger) Panic(v ...interface{}) {
	logs.Panic(v)
}
func (DefaultLogger) Panicf(format string, v ...interface{}) {
	logs.Panicf(format, v)
}
