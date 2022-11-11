package routtp

type log interface {
	Trace(string)
	Tracef(string, ...any)
	Debug(string)
	Debugf(string, ...any)
	Info(string)
	Infof(string, ...any)
	Warn(string)
	Warnf(string, ...any)
	Error(string)
	Errorf(string, ...any)
}

type nopLogger struct{}

var logger log = nopLogger{}

func SetLogger(l log) {
	logger = l
}

func (l nopLogger) Trace(string)          {}
func (l nopLogger) Tracef(string, ...any) {}
func (l nopLogger) Debug(string)          {}
func (l nopLogger) Debugf(string, ...any) {}
func (l nopLogger) Info(string)           {}
func (l nopLogger) Infof(string, ...any)  {}
func (l nopLogger) Warn(string)           {}
func (l nopLogger) Warnf(string, ...any)  {}
func (l nopLogger) Error(string)          {}
func (l nopLogger) Errorf(string, ...any) {}
