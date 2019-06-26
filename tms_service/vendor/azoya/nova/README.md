#nova
# 简介

nova是我们内部使用的一个golang api框架，现在版本基于gin进行二次开发，在gin的基础上进行改进，精简并且加入一些我们需要的功能，最终形成一个满足我们需求的高性能、健壮的框架。PS:gin使用的是MIT的开源协议，允许我们进行修改并且闭源。

查看gin文档：https://github.com/gin-gonic/gin

# Change Log
## v0.1.0(2017/10/20)
1. 添加入jaeger sdk 
2. 添加promethues监控 
3. 替换json反序列化库 
4. 内置配置文件对象 

# Start to use
与gin使用方式保持一致，但引入的包更换为azoya/nova，下面的代码可以在nova/example/api下找到
```
package main

import (
	"fmt"
	"azoya/nova"
	"azoya/nova/example/api/v1/controllers"
)

func main() {
	router := nova.Default()
	nova.SetMode("debug")

	sample := controllers.NewSampleController()

	router.GET("/sample/loggerinfo", sample.LoggerInfo)
	router.GET("/sample/sqlquery", sample.SQLQuery)
	router.GET("/sample/redis", sample.Redis)
	router.GET("/sample/request", sample.Request)
	router.GET("/sample/startspan", sample.StartSpanFromContext)

	addr := "0.0.0.0"

	port := router.Configer.String("listen::port")

	if port != "" {
		addr = fmt.Sprintf("%s:%s", addr, port)
	}

	router.Run(addr)
}
```

## 注意事项
在引入项目后，会有一部分包在本地没有，这个时候就需要到项目下执行
```
go get
```
然而有一部分库在国内无法访问链接，比如golang.org的库，这个时候可以登陆国内的一些代理来下载
```
https://www.golangtc.com/download/package
```
最后，在编译的时候会报thrift库出错
```
github.com/uber/tchannel-go/thrift/gen-go/meta/meta.go:289: cannot use metaProcessorHealth literal (type *metaProcessorHealth) as type thrift.TProcessorFunction in assignment: 
*metaProcessorHealth does not implement thrift.TProcessorFunction (wrong type for Process method)
```

这个时候需要执行以下操作，把 GOPATH/src/github.com/apache/thrift/ 的分支从master切换到 remotes/origin/0.10.0，问题解决。
```
cd GOPATH/src/github.com/apache/thrift/
git checkout -b remotes/origin/0.10.0 remotes/origin/0.10.0
```

## 配置文件说明
默认配置文件为项目主文件同级的config.conf，暂还不可以指定文件名，使用toml格式。下面的为默认配置项，其余配置根据自己实际需要可以进行增加。
```
[service]
 # service name是指当前运行的服务名称，这个名称应该是全局服务唯一，不然会在jaeger或者其他用到的地方造成混淆,默认项是service_example
name = "current_example"  

[monitor]
# 监控模块是否开启,默认为开启(enable),禁用为(disable)，如果禁用jaeger和prometheus不开启
# 这里还可以拆为更小粒度的控制，就是分别控制jaeger和prometheus是否启用
status = "enable"

[metrics]
# status = enable 为需要验证，默认为需要enable。disable为不需要验证
auth_status = "enable"
# 默认的验证key为 auth,需要auth和token搭配才能请求成功
auth_token = "12121212"

[listen]
# service启动的ip
host = "0.0.0.0" 
# service启动的端口号
port = "9091"  
```

### 配置文件调用
从nova engine中的Configer对象中直接读取就可以。传入的key值为配置文件中的listen::port。
```
router := nova.Default()
port := router.Configer.String("listen::port")
```
也可以在Context对象中调用
```
ctx.Configer().String("listen::port")
```

## Promethues集成
框架对外暴露了一个metrics的地址来给给Promethues来抓取，提供了包括了请求次数、响应速度、CPU、GC回收等指标。
```
http://host:port/metrics
```

## Jaeger使用
jaeger是一个uber开源的一个分布式追踪系统，主要用于解决分布式系统之间调用链的追踪问题，框架中也已经支持，使用框架时可以结合jaeger来解决定位和调试问题。
在使用jaeger时需要在服务器上安装一个agent，其他的收集与显示等不需要独立部署，相关部署见[wiki](http://wiki.i.azoyagroup.com/pages/viewpage.action?pageId=5410114)

### 显性加入span追踪
框架已经有自带整个请求到结束的追踪，如果还想追踪自己某个部分的代码响应时间和参数，可以像下方代码一样显性添加，或者自己封装进类里。
使用步骤：
1. 从当前上下文(Context)中创建出一个span，传入c.Context()和当前的span名称，比如可以传入一个函数名称，用于在jaeger后台进行查看
2. 写业务代码
3. span.Finish()结束这个span，这个span就记录步骤2的执行时间
4. 把这个ctx放入到context中，让整个追踪链串通起来。

```
//StartSpanFromContext 演示 StartSpanFromContext的使用方法
func (controller *SampleController) StartSpanFromContext(c *nova.Context) {
	span, ctx := opentracing.StartSpanFromContext(c.Context(), "test start span from context func")
	//do something
	defer span.Finish()
	c.WithContext(ctx)

	c.JSON(http.StatusOK, nova.H{"code": "success"})
}
```
### 两个服务间如何追踪
需要使用封装好的http库来发送网络请求，在构造网络请求时，会在header中加入trace和span等信息，然后在接收服务时对这些请求数据做处理，则上下两个服务的调用就串联起来了。使用方法如下：
```
//Result 用于接收http请求的返回值
type Result struct {
	Code string `json:"code"`
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

```