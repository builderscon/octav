// +build !debug
// +build gcp

package log

import (
	"context"
	"os"

	"github.com/pkg/errors"

	"cloud.google.com/go/logging"
)

type gcpLog struct {
	client *logging.Client
	logger *logging.Logger
}

func init() {
	l, err := NewGCPLog(context.Background(),
		os.Getenv("GOOGLE_PROJECT_ID"),
		os.Getenv("GOOGLE_LOG_ID"),
	)
	if err != nil {
		panic(err.Error())
	}
	DefaultLogger = l
}

func NewGCPLog(ctx context.Context, projID, logID string) (*gcpLog, error) {
	l := &gcpLog{}
	cl, err := logging.NewClient(ctx, projID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stackdriver log client")
	}

	lg := cl.Logger(logID)
	l.client = cl
	l.logger = lg
	return l, nil
}

func (l *gcpLog) Log(s Severity, payload interface{}) {
	var e logging.Entry
	var ok bool
	e, ok = payload.(logging.Entry)
	if !ok {
		e.Payload = payload
	}

	l.logger.Log(e)
}

func toGoogleLoggingLevel(s Severity) logging.Severity {
	switch s {
	case LDebug:
		return logging.Debug
	case LInfo:
		return logging.Info
	case LNotice:
		return logging.Notice
	case LWarning:
		return logging.Warning
	case LError:
		return logging.Error
	case LCritical:
		return logging.Critical
	case LAlert:
		return logging.Alert
	case LEmergency:
		return logging.Emergency
	default:
		return logging.Default
	}
}

func (l *gcpLog) Debug(payload interface{})     { l.Log(LDebug, payload) }
func (l *gcpLog) Info(payload interface{})      { l.Log(LInfo, payload) }
func (l *gcpLog) Notice(payload interface{})    { l.Log(LNotice, payload) }
func (l *gcpLog) Warning(payload interface{})   { l.Log(LWarning, payload) }
func (l *gcpLog) Error(payload interface{})     { l.Log(LError, payload) }
func (l *gcpLog) Critical(payload interface{})  { l.Log(LCritical, payload) }
func (l *gcpLog) Alert(payload interface{})     { l.Log(LAlert, payload) }
func (l *gcpLog) Emergency(payload interface{}) { l.Log(LEmergency, payload) }
