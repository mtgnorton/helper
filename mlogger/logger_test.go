package mlogger

import "testing"

func TestLogger(t *testing.T) {

	l := NewLogger(WithLevel(TraceLevel), WithCallerSkipCount(2), WithFields(map[string]interface{}{"traceID": "trace-1"})).
		Fields(map[string]interface{}{"requestID": "req-1"})

	l.Log(TraceLevel, "test")

	l = NewLogger(WithLevel(InfoLevel), WithCallerSkipCount(2), WithFields(map[string]interface{}{"traceID": "trace-2"})).
		Fields(map[string]interface{}{"requestID": "req-2"})

	l.Log(TraceLevel, "test")

}
