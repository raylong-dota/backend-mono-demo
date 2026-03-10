package zap

import (
	"testing"

	klog "github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap/zapcore"
)

func TestNew_NotNil(t *testing.T) {
	logger := New()
	if logger == nil {
		t.Fatal("New() returned nil")
	}
}

func TestNew_ImplementsKlogLogger(t *testing.T) {
	var _ klog.Logger = New()
}

func TestNew_LogsWithoutPanic(t *testing.T) {
	logger := New()
	if err := logger.Log(klog.LevelInfo, "msg", "hello"); err != nil {
		t.Errorf("Log() returned error: %v", err)
	}
}

func TestWithLevel(t *testing.T) {
	cases := []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
	}
	for _, lvl := range cases {
		logger := New(WithLevel(lvl))
		if logger == nil {
			t.Errorf("New(WithLevel(%v)) returned nil", lvl)
		}
	}
}

func TestDefaultOptions(t *testing.T) {
	cfg := defaultOptions()
	if cfg.level != zapcore.InfoLevel {
		t.Errorf("default level = %v, want InfoLevel", cfg.level)
	}
}
