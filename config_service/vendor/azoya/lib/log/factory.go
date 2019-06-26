// 重写""github.com/uber/jaeger/examples/hotrod/pkg/log/factory.go"文件
// 主要是把x/net/context对象更改为gin/context对象
package log

import (
	// "gopkg.in/gin-gonic/gin.v1"
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/opentracing/opentracing-go"
)
// Factory is the default logging wrapper that can create
// logger instances either for a given Context or context-less.
type Factory struct {
	logger *zap.Logger
}

// NewFactory creates a new Factory.
func NewFactory(logger *zap.Logger) Factory {
	return Factory{logger: logger}
}

// Bg creates a context-unaware logger.
func (b Factory) Bg() Logger {
	return logger{
		logger: b.logger,
	}
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func (b Factory) For(ctx context.Context) Logger {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		// TODO for Jaeger span extract trace/span IDs as fields
		return spanLogger{span: span, logger: b.logger}
	}
	return b.Bg()
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (b Factory) With(fields ...zapcore.Field) Factory {
	return Factory{logger: b.logger.With(fields...)}
}