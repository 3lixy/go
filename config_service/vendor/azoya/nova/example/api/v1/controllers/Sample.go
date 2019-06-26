package controllers

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"

	"azoya/nova"
	"azoya/lib/tracing"
)

//Result 用于接收http请求的返回值
type Result struct {
	Code string `json:"code"`
}

//SampleController controller 都继承Monitor
type SampleController struct {
}

//NewSampleController 初始化，把初始化的Monitor对象放入
func NewSampleController() *SampleController {
	return &SampleController{}
}

//LoggerInfo 普通日志记录的示例，在span下面追加log
func (controller *SampleController) LoggerInfo(c *nova.Context) {
	c.Logger().Info("abcde")
	c.JSON(http.StatusOK, nova.H{"code": "aaa"})
}

//SQLQuery 数据库查询的示例
func (controller *SampleController) SQLQuery(c *nova.Context) {
	//新建一个tracer，那在jaeger里面就是一个单独的tracer，可以在ui上进行筛选
	if span := opentracing.SpanFromContext(c.Context()); span != nil {
		tracer := tracing.Init("SQL QUERY", c.LoggerFactory())
		span := tracer.StartSpan("SQL SELECT", opentracing.ChildOf(span.Context()))
		tags.SpanKindRPCClient.Set(span)
		tags.PeerService.Set(span, "mysql")
		sqlQuery := "SELECT * FROM customer WHERE customer_id=1111"
		span.SetTag("sql.query", sqlQuery)
		c.Logger().Info(sqlQuery)
		defer span.Finish()
		c.WithSpan(span)
	}
	c.JSON(http.StatusOK, nova.H{"code": "sql query"})
}

//Redis redis读取的示例
func (controller *SampleController) Redis(c *nova.Context) {
	if span := opentracing.SpanFromContext(c.Context()); span != nil {
		//直接使用旧的tracer，然后进行新建span，这样子不能在ui上进行筛选
		span := c.Tracer().StartSpan("redis select", opentracing.ChildOf(span.Context()))
		span.SetTag("redis.key", "test")
		defer span.Finish()
		c.WithSpan(span)
	}

	c.JSON(http.StatusOK, nova.H{"code": "redis"})
}

//StartSpanFromContext 演示 StartSpanFromContext的使用方法
func (controller *SampleController) StartSpanFromContext(c *nova.Context) {
	span, ctx := opentracing.StartSpanFromContext(c.Context(), "test start span from context func")
	//do something
	defer span.Finish()

	c.WithContext(ctx)

	c.JSON(http.StatusOK, nova.H{"code": "success"})
}

//Request 发送http请求
func (controller *SampleController) Request(c *nova.Context) {
	var result Result
	url := "http://127.0.0.1:9091/sample/redis"
	if err := c.Client().GetJSON(c.Context(), "/getSample/Redis", url, &result); err != nil {
		c.JSON(http.StatusOK, err)
	}

	c.JSON(http.StatusOK, nova.H{"code": result.Code})
}
