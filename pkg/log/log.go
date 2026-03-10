package log

import (
	"context"

	"github.com/ray-dota/backend-mono/pkg/log/zap"

	klog "github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap/zapcore"
)

// Logger 是底层日志驱动接口，与 github.com/go-kratos/kratos/v2/log.Logger 完全兼容。
// 基础设施层（server/data/biz）通过 wire 注入此类型；业务层使用 Helper。
type Logger = klog.Logger

// Level 日志级别，封装在 pkg/log，外部无需 import zap
type Level int8

const (
	LevelDebug Level = iota - 1
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

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

func NewHelper(logger klog.Logger) Helper {
	return &helper{h: klog.NewHelper(logger)}
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

// Build 组装生产级 Logger，注入服务元信息
// 返回 Logger，框架层和 NewHelper 都可以使用
// log level 从环境变量 LOG_LEVEL 读取，默认 info
func BuildLogger(service, version string, lv Level) Logger {
	return klog.With(
		zap.New(zap.WithLevel(toZapLevel(lv))),
		"service", service,
		"version", version,
	)
}

// toZapLevel 将 pkg/log.Level 转换为 zapcore.Level，仅内部使用
func toZapLevel(l Level) zapcore.Level {
	switch l {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
