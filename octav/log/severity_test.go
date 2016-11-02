package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeverityString(t *testing.T) {
	data := map[Severity]string{
		LDefault:   "Default",
		LDebug:     "Debug",
		LInfo:      "Info",
		LNotice:    "Notice",
		LWarning:   "Warning",
		LError:     "Error",
		LCritical:  "Critical",
		LAlert:     "Alert",
		LEmergency: "Emergency",
	}

	for s, expected := range data {
		t.Logf("%s", s)
		if !assert.Equal(t, expected, s.String()) {
			return
		}
	}
}
