package httperr

import (
	"context"
	"github.com/sirupsen/logrus"
)

// LogrusErrorReporter reports Errors as log entry.
type LogrusErrorReporter struct {
	log *logrus.Entry
}

// NewLogrusErrorReporter creates new instance of LogrusErrorReporter
func NewLogrusErrorReporter(l *logrus.Entry) *LogrusErrorReporter {
	return &LogrusErrorReporter{
		log: l,
	}
}

// Report reports error as log entry
func (r *LogrusErrorReporter) Report(ctx context.Context, err error) {
	if err != nil {
		r.log.Errorf("Handler got application error: %s", err.Error())
	}
}
