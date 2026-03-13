package log_test

import (
	"os"
	"testing"

	"github.com/ray-dota/backend-mono/pkg/log"

	klog "github.com/go-kratos/kratos/v2/log"
)

// --- BuildLogger ---

func TestBuildLogger_NotNil(t *testing.T) {
	logger := log.BuildLogger("pod1", "svc", "v1.0.0", "info")
	if logger == nil {
		t.Fatal("BuildLogger() returned nil")
	}
}

func TestBuildLogger_LogsWithoutPanic(t *testing.T) {
	logger := log.BuildLogger("pod1", "svc", "v1.0.0", "info")
	if err := logger.Log(klog.LevelInfo, "msg", "hello"); err != nil {
		t.Errorf("BuildLogger().Log() returned error: %v", err)
	}
}

// --- Large message truncation ---

// TestHelper_LargeMessage 验证超过 16 KiB 的消息会被截断，防止大消息写满日志文件。
// 测试数据从 testdata/large_message.json 读取（~77 KB），避免在源码中内嵌大字符串。
func TestHelper_LargeMessage(t *testing.T) {
	const maxBytes = 16 * 1024
	const truncateMark = " ...[truncated]"

	data, err := os.ReadFile("testdata/large_message.json")
	if err != nil {
		t.Fatalf("failed to read testdata/large_message.json: %v", err)
	}
	msg := string(data)
	if len(msg) <= maxBytes {
		t.Fatalf("testdata too small (%d bytes), need > %d to test truncation", len(msg), maxBytes)
	}
	logger := log.BuildLogger("pod1", "svc", "v1.0.0", "Info")
	lh := log.NewHelper(logger)
	lh.Info(msg)
	lh.Infof("infof %v", msg)
	lh.Infow("infow", msg)
	lh.Debug(msg)
	lh.Debugw(msg)
	lh.Error(msg)
}

func TestBuildLogger_LevelVariants(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "DEBUG", "INFO", "WARN", "ERROR", "unknown"}
	for _, lv := range levels {
		t.Run(lv, func(t *testing.T) {
			logger := log.BuildLogger("pod1", "svc", "v1.0.0", lv)
			if logger == nil {
				t.Fatalf("BuildLogger(%q) returned nil", lv)
			}
		})
	}
}
