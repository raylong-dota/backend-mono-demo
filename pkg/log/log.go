package log

import (
	"context"
	"os"
	"strings"

	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	uzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 是底层日志驱动接口，与 github.com/go-kratos/kratos/v2/log.Logger 完全兼容。
// 基础设施层（server/data/biz）通过 wire 注入此类型；业务层使用 Helper。
type Logger = klog.Logger

// 如果 logger 还没初始化，可以使用默认的全局 logger
// stdout 输出，简单文本格式
func Info(a ...any)                  { klog.Info(a...) }
func Warn(a ...any)                  { klog.Warn(a...) }
func Error(a ...any)                 { klog.Error(a...) }
func Fatal(a ...any)                 { klog.Fatal(a...) }
func Infof(format string, a ...any)  { klog.Infof(format, a...) }
func Warnf(format string, a ...any)  { klog.Warnf(format, a...) }
func Errorf(format string, a ...any) { klog.Errorf(format, a...) }
func Fatalf(format string, a ...any) { klog.Fatalf(format, a...) }

// maxLogSize 单条日志字段最大长度，超出则截断
// 16KB 足以覆盖完整堆栈信息，同时防止异常大消息撑爆日志采集器
const maxLogSize = 16 * 1024 // 16 kb

// truncateCore 包装底层 Core，对超过 maxLogSize 的字段进行截断
type truncateCore struct {
	zapcore.Core
}

func (c *truncateCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if len(entry.Message) > maxLogSize {
		entry.Message = entry.Message[:maxLogSize] + "...[truncated]"
	}
	return c.Core.Check(entry, ce)
}

func (c *truncateCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	if len(entry.Message) > maxLogSize {
		entry.Message = entry.Message[:maxLogSize] + "...[truncated]"
	}
	for i, f := range fields {
		if f.Type == zapcore.StringType && len(f.String) > maxLogSize {
			fields[i].String = f.String[:maxLogSize] + "...[truncated]"
		}
	}
	return c.Core.Write(entry, fields)
}

type option func(*options)
type options struct {
	level zapcore.Level
}

func defaultOptions() *options {
	return &options{
		level: zapcore.InfoLevel,
	}
}

// ParseLevel 将字符串解析为 zapcore.Level
// 支持 debug / info / warn / error / fatal，未识别时返回 InfoLevel。
func parseLevel(lv string) zapcore.Level {
	switch strings.ToUpper(lv) {
	case "DEBUG":
		return zapcore.DebugLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func withLevel(level zapcore.Level) option {
	return func(o *options) {
		o.level = level
	}
}

// newZap 返回 klog.Logger，pkg/log 内部使用
// 固定输出字段：ts / level / msg / caller
func newZap(opts ...option) klog.Logger {
	cfg := defaultOptions()
	for _, o := range opts {
		o(cfg)
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:      "ts",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.RFC3339NanoTimeEncoder,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder, // "biz/user.go:42"
	}

	core := &truncateCore{
		Core: zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.AddSync(os.Stdout), // K8s 体系只输出 stdout
			cfg.level,
		),
	}

	zapLogger := uzap.New(core, uzap.AddCaller(), uzap.AddCallerSkip(3))

	return kzap.NewLogger(zapLogger)
}

// Build 组装生产级 Logger，注入服务元信息
// 返回 Logger，框架层和 NewHelper 都可以使用
func BuildLogger(hostname, service, version string, lv string) Logger {
	return klog.With(
		newZap(withLevel(parseLevel(lv))),
		// 字段名来自约定 https://opentelemetry.io/docs/specs/semconv/resource/service/
		"service.instance.id", hostname,
		"service.name", service,
		"service.version", version,
	)
}

// Helper 业务层唯一依赖的接口
type Helper interface {
	WithContext(ctx context.Context) Helper

	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)

	Debugw(keyvals ...any)
	Infow(keyvals ...any)
	Warnw(keyvals ...any)
	Errorw(keyvals ...any)

	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

// helper 是 Helper 接口的内部实现，对外不可见
type helper struct {
	h *klog.Helper
}

func (h *helper) WithContext(ctx context.Context) Helper {
	return &helper{h: h.h.WithContext(ctx)}
}

func (h *helper) Debug(args ...any) { h.h.Debug(args...) }
func (h *helper) Info(args ...any)  { h.h.Info(args...) }
func (h *helper) Warn(args ...any)  { h.h.Warn(args...) }
func (h *helper) Error(args ...any) { h.h.Error(args...) }

func (h *helper) Debugw(keyvals ...any) { h.h.Debugw(keyvals...) }
func (h *helper) Infow(keyvals ...any)  { h.h.Infow(keyvals...) }
func (h *helper) Warnw(keyvals ...any)  { h.h.Warnw(keyvals...) }
func (h *helper) Errorw(keyvals ...any) { h.h.Errorw(keyvals...) }

func (h *helper) Debugf(format string, args ...any) { h.h.Debugf(format, args...) }
func (h *helper) Infof(format string, args ...any)  { h.h.Infof(format, args...) }
func (h *helper) Warnf(format string, args ...any)  { h.h.Warnf(format, args...) }
func (h *helper) Errorf(format string, args ...any) { h.h.Errorf(format, args...) }

// maxMsgBytes 是单条消息字段允许的最大字节数（16 KiB）。
// 超出部分将被截断并附加截断标记，防止大消息写满日志文件。
const maxMsgBytes = 16 * 1024

func truncateString(s string) string {
	if len(s) <= maxMsgBytes {
		return s
	}
	return s[:maxMsgBytes] + " ...[truncated]"
}

// truncatingLogger 包装底层 Logger，对所有 string 类型的 keyval 值执行截断。
type truncatingLogger struct {
	Logger
}

func (t *truncatingLogger) Log(level klog.Level, keyvals ...any) error {
	for i, v := range keyvals {
		if s, ok := v.(string); ok {
			keyvals[i] = truncateString(s)
		}
	}
	return t.Logger.Log(level, keyvals...)
}

func NewHelper(logger klog.Logger) Helper {
	return &helper{h: klog.NewHelper(&truncatingLogger{logger})}
}
