// +build !gcp
// +build debug

package log

import (
	"io"
	"os"
)

func init() {
	DefaultLogger = NewDebugLog(os.Stderr)
}

type debugLog struct {
	dst io.Writer
}

func NewDebugLog(dst io.Writer) *debugLog {
	l := &debugLog{}
	l.dst = dst
	return l
}

func (l *debugLog) Log(s Severity, payload interface{}) {
	writeLog(l.dst, s, payload)
}

func (l *debugLog) Debug(payload interface{})     { l.Log(LDebug, payload) }
func (l *debugLog) Info(payload interface{})      { l.Log(LInfo, payload) }
func (l *debugLog) Notice(payload interface{})    { l.Log(LNotice, payload) }
func (l *debugLog) Warning(payload interface{})   { l.Log(LWarning, payload) }
func (l *debugLog) Error(payload interface{})     { l.Log(LError, payload) }
func (l *debugLog) Critical(payload interface{})  { l.Log(LCritical, payload) }
func (l *debugLog) Alert(payload interface{})     { l.Log(LAlert, payload) }
func (l *debugLog) Emergency(payload interface{}) { l.Log(LEmergency, payload) }
