package nova

import (
	//open tracing
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	//azoya
	"azoya/lib/log"
	"azoya/lib/tracing"
)

//Monitor 初始化一个监控类，里面包含了gin，tracer，等对象，作为主对象往后传递
type Monitor struct {
	//open trace对象，用于创建span等
	tracer opentracing.Tracer

	//日志对象，用于做日志输出
	logger *log.Logger

	//prometheus 用于监控
	prometheus *Prometheus
}

//NewDevelopment 初始化Monitor对象，传入参数为当前服务名称,但是logger日志是开发模式
//打印的日志比较友好，便于查看
func NewDevelopment(serviceName string) Monitor {
	cfg := log.DevelopmentConfig()
	cfg.Program = serviceName
	cfg.FileEnabled = false
	cfg.StdoutEnabled = true
	logger, _ := cfg.Build()
	return initMonitor(serviceName, logger)
}

//NewProduction 初始化Monitor对象，传入参数为当前服务名称,但是logger日志是生产模式
//打印的日志是json格式
func NewProduction(serviceName string) Monitor {
	cfg := log.ProductionConfig()
	cfg.Program = serviceName
	cfg.FileEnabled = false
	cfg.StdoutEnabled = true
	logger, _ := cfg.Build()
	return initMonitor(serviceName, logger)
}

//initMonitor 初始化Monitor对象
func initMonitor(serviceName string, logger *log.Logger) Monitor {

	tracer := tracing.Init(serviceName, logger)
	prometheus := NewPrometheus()

	return Monitor{
		tracer:     tracer,
		logger:     logger,
		prometheus: prometheus,
	}

}

//Tracer tracer对象可以用来创建span，做Monitor相关的操作
func (s *Monitor) Tracer() opentracing.Tracer {
	return s.tracer
}

//Logger logger使用来做日志操作，如果有open tracing则会记录info到Monitor
func (s *Monitor) Logger() *log.Logger {
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
			return
		}
		ctx, span := tracing.InitTracing(c.Context(), c.Tracer(), c.Request)
		c = c.WithContext(ctx)

		c.Next()

		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		span.Finish()
	}
}
