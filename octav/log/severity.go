package log

import (
	"fmt"
	"io"
	"strings"
)

const _Severity_name = "DefaultDebugInfoNoticeWarningErrorCriticalAlertEmergency"

var _Severity_index = [...]uint8{0, 7, 12, 16, 22, 29, 34, 42, 47, 56}

func (i Severity) String() string {
	i -= 1
	if i < 0 || i >= Severity(len(_Severity_index)-1) {
		return fmt.Sprintf("Severity(%d)", i+1)
	}
	return _Severity_name[_Severity_index[i]:_Severity_index[i+1]]
}

func (i Severity) WritePrefix(w io.Writer) {
	out := bufpool.Get()
	defer bufpool.Release(out)

	out.WriteByte('[')
	out.WriteString(strings.ToUpper(i.String()))
	out.WriteByte(']')
	out.WriteByte(' ')

	out.WriteTo(w)
}
