package log

import (
	"context"
	"testing"

	klog "github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap/zapcore"
)

// spyLogger 记录所有 Log 调用，用于断言。
type spyLogger struct {
	entries [][]interface{}
}

func (s *spyLogger) Log(_ klog.Level, keyvals ...interface{}) error {
	s.entries = append(s.entries, keyvals)
	return nil
}

func (s *spyLogger) logged() bool { return len(s.entries) > 0 }

func newSpyHelper() (Helper, *spyLogger) {
	spy := &spyLogger{}
	return NewHelper(spy), spy
}

// --- NewHelper ---

func TestNewHelper_NotNil(t *testing.T) {
	h, _ := newSpyHelper()
	if h == nil {
		t.Fatal("NewHelper() returned nil")
	}
}

// --- Helper 方法覆盖 ---

func TestHelper_InfoMethods(t *testing.T) {
	tests := []struct {
		name string
		call func(Helper)
	}{
		{"Info", func(h Helper) { h.Info("msg") }},
		{"Infow", func(h Helper) { h.Infow("key", "val") }},
		{"Infof", func(h Helper) { h.Infof("fmt %s", "arg") }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, spy := newSpyHelper()
			tc.call(h)
			if !spy.logged() {
				t.Errorf("%s: expected log entry, got none", tc.name)
			}
		})
	}
}

func TestHelper_WarnMethods(t *testing.T) {
	tests := []struct {
		name string
		call func(Helper)
	}{
		{"Warn", func(h Helper) { h.Warn("msg") }},
		{"Warnw", func(h Helper) { h.Warnw("key", "val") }},
		{"Warnf", func(h Helper) { h.Warnf("fmt %s", "arg") }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, spy := newSpyHelper()
			tc.call(h)
			if !spy.logged() {
				t.Errorf("%s: expected log entry, got none", tc.name)
			}
		})
	}
}

func TestHelper_ErrorMethods(t *testing.T) {
	tests := []struct {
		name string
		call func(Helper)
	}{
		{"Error", func(h Helper) { h.Error("msg") }},
		{"Errorw", func(h Helper) { h.Errorw("key", "val") }},
		{"Errorf", func(h Helper) { h.Errorf("fmt %s", "arg") }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, spy := newSpyHelper()
			tc.call(h)
			if !spy.logged() {
				t.Errorf("%s: expected log entry, got none", tc.name)
			}
		})
	}
}

func TestHelper_DebugMethods(t *testing.T) {
	tests := []struct {
		name string
		call func(Helper)
	}{
		{"Debug", func(h Helper) { h.Debug("msg") }},
		{"Debugw", func(h Helper) { h.Debugw("key", "val") }},
		{"Debugf", func(h Helper) { h.Debugf("fmt %s", "arg") }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, spy := newSpyHelper()
			tc.call(h)
			if !spy.logged() {
				t.Errorf("%s: expected log entry, got none", tc.name)
			}
		})
	}
}

// --- WithContext ---

func TestHelper_WithContext_ReturnsNewHelper(t *testing.T) {
	h, _ := newSpyHelper()
	h2 := h.WithContext(context.Background())
	if h2 == nil {
		t.Fatal("WithContext() returned nil")
	}
	if h2 == h {
		t.Fatal("WithContext() should return a new Helper instance")
	}
}

func TestHelper_WithContext_StillLogs(t *testing.T) {
	h, spy := newSpyHelper()
	h2 := h.WithContext(context.Background())
	h2.Info("after context")
	if !spy.logged() {
		t.Fatal("Helper from WithContext should still route to original logger")
	}
}

// --- BuildLogger ---

func TestBuildLogger_NotNil(t *testing.T) {
	logger := BuildLogger("svc", "v1.0.0", LevelInfo)
	if logger == nil {
		t.Fatal("BuildLogger() returned nil")
	}
}

func TestBuildLogger_LogsWithoutPanic(t *testing.T) {
	logger := BuildLogger("svc", "v1.0.0", LevelInfo)
	if err := logger.Log(klog.LevelInfo, "msg", "hello"); err != nil {
		t.Errorf("BuildLogger().Log() returned error: %v", err)
	}
}

// --- toZapLevel ---

func TestToZapLevel(t *testing.T) {
	cases := []struct {
		in   Level
		want zapcore.Level
	}{
		{LevelDebug, zapcore.DebugLevel},
		{LevelInfo, zapcore.InfoLevel},
		{LevelWarn, zapcore.WarnLevel},
		{LevelError, zapcore.ErrorLevel},
		{LevelFatal, zapcore.FatalLevel},
		{Level(99), zapcore.InfoLevel}, // default
	}
	for _, tc := range cases {
		if got := toZapLevel(tc.in); got != tc.want {
			t.Errorf("toZapLevel(%d) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
