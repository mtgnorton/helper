package mlogger

type Logger interface {
	Init(opts ...Option) error
	Options() Options
	Fields(fields map[string]interface{}) Logger
	Log(level Level, v ...interface{})
	Logf(level Level, format string, v ...interface{})
	String() string
}
