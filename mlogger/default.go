package mlogger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	DefaultLogger Logger = NewLogger()
)

func init() {
	lvl, err := GetLevel(os.Getenv("M_LOG_LEVEL"))
	if err != nil {
		lvl = DebugLevel
	}
	DefaultLogger = NewLogger(WithLevel(lvl))
}

type defaultLogger struct {
	opts Options
	sync.RWMutex
}

func NewLogger(opts ...Option) Logger {
	options := Options{
		Out:             os.Stderr,
		Context:         context.Background(),
		Fields:          make(map[string]interface{}),
		Level:           InfoLevel,
		CallerSkipCount: 2,
	}
	l := &defaultLogger{opts: options}
	if err := l.Init(opts...); err != nil {
		l.Log(FatalLevel, err)
	}
	return l
}

func (l *defaultLogger) Options() Options {
	l.RLock()
	opts := l.opts
	opts.Fields = copyFields(l.opts.Fields)
	l.RUnlock()
	return opts
}

func (l *defaultLogger) Fields(fields map[string]interface{}) Logger {
	l.Lock()
	nfields := make(map[string]interface{}, len(l.opts.Fields))
	for k, v := range l.opts.Fields {
		nfields[k] = v
	}
	l.Unlock()

	for k, v := range fields {
		nfields[k] = v
	}
	return &defaultLogger{
		opts: Options{
			Out:             l.opts.Out,
			Context:         l.opts.Context,
			Fields:          nfields,
			Level:           l.opts.Level,
			CallerSkipCount: l.opts.CallerSkipCount,
		},
	}
}

func (l *defaultLogger) Log(level Level, v ...interface{}) {
	if !l.opts.Level.Enabled(level) {
		return
	}
	l.RLock()
	fields := copyFields(l.opts.Fields)
	l.RUnlock()
	fields["level"] = level.String()

	if _, file, line, ok := runtime.Caller(l.opts.CallerSkipCount); ok {
		fields["file"] = fmt.Sprintf("%s:%d", logCallerfilePath(file), line)
	}
	keys := make([]string, 0, len(fields))
	for k, _ := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	metadata := ""
	for _, k := range keys {
		metadata += fmt.Sprintf(" %s=%v", k, fields[k])
	}

	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s %s %v\n", t, metadata, fmt.Sprint(v...))
}

func (l *defaultLogger) Logf(level Level, format string, v ...interface{}) {
	if !l.opts.Level.Enabled(level) {
		return
	}
	l.RLock()
	fields := copyFields(l.opts.Fields)
	l.RUnlock()
	fields["level"] = level.String()

	if _, file, line, ok := runtime.Caller(l.opts.CallerSkipCount); ok {
		fields["file"] = fmt.Sprintf("%s:%d", logCallerfilePath(file), line)
	}

	keys := make([]string, 0, len(fields))
	for k, _ := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	metadata := ""
	for _, k := range keys {
		metadata += fmt.Sprintf(" %s=%v", k, fields[k])
	}

	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s %s %v\n", t, metadata, fmt.Sprintf(format, v...))
}

func (l *defaultLogger) String() string {
	return "default"
}

func (l *defaultLogger) Init(opts ...Option) error {
	for _, o := range opts {
		o(&l.opts)
	}
	return nil
}

func copyFields(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}

	return dst
}

// logCallerfilePath returns a package/file:line description of the caller,
// preserving only the leaf directory name and file name.
func logCallerfilePath(loggingFilePath string) string {

	// 为了确保我们在Windows上正确修剪路径，出乎意料地需要使用'/'而不是os.PathSeparator，因为给定的路径源自Go标准库，具体来说是runtime.Caller()，截至2021年3月17日，在Windows上即使返回反斜杠也会返回正斜杠。
	//
	// See https://github.com/golang/go/issues/3335
	// and https://github.com/golang/go/issues/18151
	//
	// for discussion on the issue on Go side.
	idx := strings.LastIndexByte(loggingFilePath, '/')
	if idx == -1 {
		return loggingFilePath
	}

	idx = strings.LastIndexByte(loggingFilePath[:idx], '/')

	if idx == -1 {
		return loggingFilePath
	}

	return loggingFilePath[idx+1:]
}
