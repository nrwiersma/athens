package log

import (
	"errors"

	apierrors "github.com/gomods/athens/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Entry is an abstraction to the
// Logger and the logrus.Entry
// so that *Logger always creates
// an Entry copy which ensures no
// Fields are being overwritten.
type Entry interface {
	// Basic Logging Operation
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)

	// Attach contextual information to the logging entry
	WithFields(fields map[string]any) Entry

	// SystemErr is a method that disects the error
	// and logs the appropriate level and fields for it.
	SystemErr(err error)
}

type entry struct {
	*logrus.Entry
}

func (e *entry) WithFields(fields map[string]any) Entry {
	ent := e.Entry.WithFields(fields)
	return &entry{ent}
}

func (e *entry) SystemErr(err error) {
	aErr, ok := errors.AsType[apierrors.Error](err)
	if !ok {
		e.Error(err)
		return
	}

	ent := e.WithFields(errFields(aErr))
	switch apierrors.Severity(err) {
	case logrus.WarnLevel:
		ent.Warnf("%v", err)
	case logrus.InfoLevel:
		ent.Infof("%v", err)
	case logrus.DebugLevel:
		ent.Debugf("%v", err)
	default:
		ent.Errorf("%v", err)
	}
}

func errFields(err apierrors.Error) logrus.Fields {
	f := logrus.Fields{}
	f["operation"] = err.Op
	f["kind"] = apierrors.KindText(err)
	f["module"] = err.Module
	f["version"] = err.Version
	f["ops"] = apierrors.Ops(err)

	return f
}
