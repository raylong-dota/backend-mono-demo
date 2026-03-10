package zap

import (
	"os"

	kratoszap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	uzap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(*options)

type options struct {
	level zapcore.Level
}

func defaultOptions() *options {
	return &options{
		level: zapcore.InfoLevel,
	}
}

func WithLevel(level zapcore.Level) Option {
	return func(o *options) {
		o.level = level
	}
}

// New 返回 klog.Logger，pkg/log 内部使用
// 固定输出字段：ts / level / caller / msg
func New(opts ...Option) klog.Logger {
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

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout), // K8s 体系只输出 stdout
		cfg.level,
	)

	zapLogger := uzap.New(core, uzap.AddCaller(), uzap.AddCallerSkip(3))

	return kratoszap.NewLogger(zapLogger)
}
