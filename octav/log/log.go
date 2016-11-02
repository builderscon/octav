package log

import (
	"encoding/json"
	"io"
	"strings"
)

func Log(s Severity, payload interface{}) { DefaultLogger.Log(s, payload) }
func Debug(payload interface{})           { DefaultLogger.Debug(payload) }
func Info(payload interface{})            { DefaultLogger.Info(payload) }
func Notice(payload interface{})          { DefaultLogger.Notice(payload) }
func Warning(payload interface{})         { DefaultLogger.Warning(payload) }
func Error(payload interface{})           { DefaultLogger.Error(payload) }
func Critical(payload interface{})        { DefaultLogger.Critical(payload) }
func Alert(payload interface{})           { DefaultLogger.Alert(payload) }
func Emergency(payload interface{})       { DefaultLogger.Emergency(payload) }

func writeLog(w io.Writer, s Severity, payload interface{}) {
	out := bufpool.Get()
	defer bufpool.Release(out)
	s.WritePrefix(out)

	switch payload.(type) {
	case string:
		s := payload.(string)
		if len(s) == 0 {
			return
		}

		if !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		out.WriteString(s)
	case []byte:
		b := payload.([]byte)
		if len(b) == 0 {
			return
		}

		if b[len(b)-1] != '\n' {
			b = append(b, '\n')
		}
		out.Write(b)
	default:
		json.NewEncoder(out).Encode(payload)
		out.WriteByte('\n')
	}

	out.WriteTo(w)
}
