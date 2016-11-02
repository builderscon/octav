package log

import bufferpool "github.com/lestrrat/go-bufferpool"

var DefaultLogger Logger
var bufpool = bufferpool.New()

type Severity int

const (
	LDefault Severity = iota + 1
	LDebug
	LInfo
	LNotice
	LWarning
	LError
	LCritical
	LAlert
	LEmergency
)

type Logger interface {
	Log(Severity, interface{})
	Debug(interface{})
	Info(interface{})
	Notice(interface{})
	Warning(interface{})
	Error(interface{})
	Critical(interface{})
	Alert(interface{})
	Emergency(interface{})
}

type nullLog struct{}
