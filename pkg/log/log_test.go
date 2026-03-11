package log_test

import (
	"context"
	"testing"

	"github.com/ray-dota/backend-mono/pkg/log"

	klog "github.com/go-kratos/kratos/v2/log"
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

func newSpyHelper() (log.Helper, *spyLogger) {
	spy := &spyLogger{}
	return log.NewHelper(spy), spy
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
		call func(log.Helper)
	}{
		{"Info", func(h log.Helper) { h.Info("msg") }},
		{"Infow", func(h log.Helper) { h.Infow("key", "val") }},
		{"Infof", func(h log.Helper) { h.Infof("fmt %s", "arg") }},
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
		call func(log.Helper)
	}{
		{"Warn", func(h log.Helper) { h.Warn("msg") }},
		{"Warnw", func(h log.Helper) { h.Warnw("key", "val") }},
		{"Warnf", func(h log.Helper) { h.Warnf("fmt %s", "arg") }},
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
		call func(log.Helper)
	}{
		{"Error", func(h log.Helper) { h.Error("msg") }},
		{"Errorw", func(h log.Helper) { h.Errorw("key", "val") }},
		{"Errorf", func(h log.Helper) { h.Errorf("fmt %s", "arg") }},
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
		call func(log.Helper)
	}{
		{"Debug", func(h log.Helper) { h.Debug("msg") }},
		{"Debugw", func(h log.Helper) { h.Debugw("key", "val") }},
		{"Debugf", func(h log.Helper) { h.Debugf("fmt %s", "arg") }},
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
	logger := log.BuildLogger("svc", "v1.0.0", "info")
	if logger == nil {
		t.Fatal("BuildLogger() returned nil")
	}
}

func TestBuildLogger_LogsWithoutPanic(t *testing.T) {
	logger := log.BuildLogger("svc", "v1.0.0", "info")
	if err := logger.Log(klog.LevelInfo, "msg", "hello"); err != nil {
		t.Errorf("BuildLogger().Log() returned error: %v", err)
	}
}

func TestBuildLogger_LevelVariants(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "DEBUG", "INFO", "WARN", "ERROR", "unknown"}
	for _, lv := range levels {
		t.Run(lv, func(t *testing.T) {
			logger := log.BuildLogger("svc", "v1.0.0", lv)
			if logger == nil {
				t.Fatalf("BuildLogger(%q) returned nil", lv)
			}
		})
	}
}
