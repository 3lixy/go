package nova

import (
	//open tracing
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	//azoya
	"azoya/lib/log"
	"azoya/lib/tracing"
	"github.com/opentracing/opentracing-go/ext"
)

//Monitor 初始化一个监控类，里面包含了gin，tracer，等对象，作为主对象往后传递
type Monitor struct {
	//open trace对象，用于创建span等
	tracer opentracing.Tracer

	//日志对象，用于做日志输出
	logger log.Factory

	//prometheus 用于监控
	prometheus *Prometheus
}

//NewDevelopment 初始化Monitor对象，传入参数为当前服务名称,但是logger日志是开发模式
//打印的日志比较友好，便于查看
func NewDevelopment(serviceName string) Monitor {
	logger, _ := zap.NewDevelopment()
	return initMonitor(serviceName, logger)
}

//NewProduction 初始化Monitor对象，传入参数为当前服务名称,但是logger日志是生产模式
//打印的日志是json格式
func NewProduction(serviceName string) Monitor {
	logger, _ := zap.NewProduction()
	return initMonitor(serviceName, logger)
}

//initMonitor 初始化Monitor对象
func initMonitor(serviceName string, logger *zap.Logger) Monitor {
	l := log.NewFactory(
		logger.With(
			zap.String("service", serviceName),
		),
	)

	tracer := tracing.Init(serviceName, l)
	prometheus := NewPrometheus()

	return Monitor{
		tracer:     tracer,
		logger:     l,
		prometheus: prometheus,
	}
}

//Tracer tracer对象可以用来创建span，做Monitor相关的操作
func (s *Monitor) Tracer() opentracing.Tracer {
	return s.tracer
}

//Logger logger使用来做日志操作，如果有open tracing则会记录info到Monitor
func (s *Monitor) Logger() log.Factory {
	return s.logger
}

//Prometheus prometheus对象用来实现监控
func (s *Monitor) Prometheus() *Prometheus {
	return s.prometheus
}

//TraceHandleFunc 用于Monitor中的span的创建，在prometheus后，运行在正常业务逻辑
//开始之前,span name为访问的url
func (s *Monitor) TraceHandleFunc() HandlerFunc {
	return func(c *Context) {
		if c.Request.URL.Path == s.prometheus.MetricsPath {
			c.Next()
			return
		}
		ctx, span := tracing.InitTracing(c.Context(),c.Tracer(),c.Request)
		c = c.WithContext(ctx)

		c.Next()

		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		span.Finish()
	}
}