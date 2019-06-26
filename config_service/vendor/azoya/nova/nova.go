// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package nova

import (
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"

	//azoya
	"azoya/lib/log"
	"azoya/lib/tracing"
	"azoya/nova/config"
)

// Version is Framework's version
const Version = "v1.0.0"

var default404Body = []byte("404 page not found")
var default405Body = []byte("405 method not allowed")

type HandlerFunc func(*Context)
type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. ie. the last handler is the main own.
func (c HandlersChain) Last() HandlerFunc {
	length := len(c)
	if length > 0 {
		return c[length-1]
	}
	return nil
}

type (
	RoutesInfo []RouteInfo
	RouteInfo  struct {
		Method  string
		Path    string
		Handler string
	}

	// Engine is the framework's instance, it contains the muxer, middleware and configuration settings.
	// Create an instance of Engine, by using New() or Default()
	Engine struct {
		RouterGroup
		secureJsonPrefix string
		allNoRoute  HandlersChain
		allNoMethod HandlersChain
		noRoute     HandlersChain
		noMethod    HandlersChain
		pool        sync.Pool
		trees       methodTrees

		// Enables automatic redirection if the current route can't be matched but a
		// handler for the path with (without) the trailing slash exists.
		// For example if /foo/ is requested but a route only exists for /foo, the
		// client is redirected to /foo with http status code 301 for GET requests
		// and 307 for all other request methods.
		RedirectTrailingSlash bool

		// If enabled, the router tries to fix the current request path, if no
		// handle is registered for it.
		// First superfluous path elements like ../ or // are removed.
		// Afterwards the router does a case-insensitive lookup of the cleaned path.
		// If a handle can be found for this route, the router makes a redirection
		// to the corrected path with status code 301 for GET requests and 307 for
		// all other request methods.
		// For example /FOO and /..//Foo could be redirected to /foo.
		// RedirectTrailingSlash is independent of this option.
		RedirectFixedPath bool

		// If enabled, the router checks if another method is allowed for the
		// current route, if the current request can not be routed.
		// If this is the case, the request is answered with 'Method Not Allowed'
		// and HTTP status code 405.
		// If no other Method is allowed, the request is delegated to the NotFound
		// handler.
		HandleMethodNotAllowed bool
		ForwardedByClientIP    bool

		// If enabled, the url.RawPath will be used to find parameters.
		UseRawPath bool

		// If true, the path value will be unescaped.
		// If UseRawPath is false (by default), the UnescapePathValues effectively is true,
		// as url.Path gonna be used, which is already unescaped.
		UnescapePathValues bool

		// Value of 'maxMemory' param that is given to http.Request's ParseMultipartForm
		// method call.
		MaxMultipartMemory int64

		//加入jaeger追踪，日志，prometheus对象
		//open trace对象，用于创建span等
		Tracer opentracing.Tracer

		//日志对象，用于做日志输出
		Logger log.Factory

		//prometheus 用于监控
		Metrics *Prometheus

		//tracing client，用这个对象发送的http请求，封装了jaeger的追踪记录
		Client *tracing.HTTPClient

		//服务名称，标志当前启动的是什么服务
		serviceName string

		Configer config.Configer
	}
)

var _ IRouter = &Engine{}

// New returns a new blank Engine instance without any middleware attached.
// By default the configuration is:
// - RedirectTrailingSlash:  true
// - RedirectFixedPath:      false
// - HandleMethodNotAllowed: false
// - ForwardedByClientIP:    true
func New(configAdapterName string) *Engine {
	debugPrintWARNINGNew()
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		secureJsonPrefix:       "while(1);",
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      false,
		HandleMethodNotAllowed: false,
		ForwardedByClientIP:    true,
		UseRawPath:             false,
		UnescapePathValues:     true,
		trees:                  make(methodTrees, 0, 9),
		Tracer:                 nil,
		Metrics:                nil,
		// Logger:                 nil,
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}

	//初始化日志
	engine.Configer = config.NewConfig(configAdapterName)
	engine.setServiceName(engine.Configer.String(config.ServiceNameKey))

	return engine
}

// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine {
	//这里的ini是指读取的本地配置文件，这里应该有一个扩展，传入consul后，可以从配置中心读取配置文件
	engine := New("ini")
	engine.Use(Logger(), Recovery())
	engine.allocateMonitor()

	return engine
}

func (engine *Engine) allocateMonitor() *Engine {
	//添加上jaeger追踪和prometheus监控
	monitor := NewProduction(engine.getServiceName())

	engine.Metrics = monitor.Prometheus()

	opentracing.InitGlobalTracer(monitor.Tracer())

	engine.Tracer = monitor.Tracer()
	engine.Logger = monitor.Logger()
	engine.Client = &tracing.HTTPClient{
		Client: &http.Client{Transport: &nethttp.Transport{}},
		Tracer: monitor.Tracer(),
	}

	//这里判断监控是否开启，但是jaeger的对象已经实例化了
	//如果想在程序中使用也是可以的
	if config.MonitorStatus(engine.Configer) {
		//Use 把Monitor和prometheus使用到gin上
		engine.Use(monitor.TraceHandleFunc())
		engine.useMetrics()
	}

	return engine
}

// useMetrics 运用metrics，分为需要验证和不需要两种
func (engine *Engine) useMetrics() {
	if config.MetricsAuthStatus(engine.Configer) {
		engine.Metrics.UseWithAuth(engine,Accounts{
			config.MetricsAuthKey : config.MetricsToken(engine.Configer),
		})
	} else {
		engine.Metrics.Use(engine)
	}
}


func (engine *Engine) allocateContext() *Context {
	return &Context{engine: engine}
}

func (engine *Engine) setServiceName(serviceName string) {
	if serviceName == "" {
		engine.serviceName = config.DefaultServiceName
	} else {
		engine.serviceName = serviceName
	}

	debugPrint("Current service name is %s\n", engine.serviceName)
}

func (engine *Engine) getServiceName() string {
	if len(engine.serviceName) == 0 {
		return config.DefaultServiceName
	}

	return engine.serviceName
}

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.noRoute = handlers
	engine.rebuild404Handlers()
}

// NoMethod sets the handlers called when... TODO.
func (engine *Engine) NoMethod(handlers ...HandlerFunc) {
	engine.noMethod = handlers
	engine.rebuild405Handlers()
}

// Use attachs a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.Use(middleware...)
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
	return engine
}

func (engine *Engine) rebuild404Handlers() {
	engine.allNoRoute = engine.combineHandlers(engine.noRoute)
}

func (engine *Engine) rebuild405Handlers() {
	engine.allNoMethod = engine.combineHandlers(engine.noMethod)
}

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")

	debugPrintRoute(method, path, handlers)
	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)
}

// Routes returns a slice of registered routes, including some useful information, such as:
// the http method, path and the handler name.
func (engine *Engine) Routes() (routes RoutesInfo) {
	for _, tree := range engine.trees {
		routes = iterate("", tree.method, routes, tree.root)
	}
	return routes
}

func iterate(path, method string, routes RoutesInfo, root *node) RoutesInfo {
	path += root.path
	if len(root.handlers) > 0 {
		routes = append(routes, RouteInfo{
			Method:  method,
			Path:    path,
			Handler: nameOfFunction(root.handlers.Last()),
		})
	}
	for _, child := range root.children {
		routes = iterate(path, method, routes, child)
	}
	return routes
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) Run(addr ...string) (err error) {
	defer func() { debugPrintError(err) }()

	address := resolveAddress(addr)
	debugPrint("Listening and serving HTTP on %s\n", address)
	err = http.ListenAndServe(address, engine)
	return
}

// RunTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests.
// It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunTLS(addr, certFile, keyFile string) (err error) {
	debugPrint("Listening and serving HTTPS on %s\n", addr)
	defer func() { debugPrintError(err) }()

	err = http.ListenAndServeTLS(addr, certFile, keyFile, engine)
	return
}

// RunUnix attaches the router to a http.Server and starts listening and serving HTTP requests
// through the specified unix socket (ie. a file).
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunUnix(file string) (err error) {
	debugPrint("Listening and serving HTTP on unix:/%s", file)
	defer func() { debugPrintError(err) }()

	os.Remove(file)
	listener, err := net.Listen("unix", file)
	if err != nil {
		return
	}
	defer listener.Close()
	err = http.Serve(listener, engine)
	return
}

// ServeHTTP conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = req
	c.reset()

	engine.handleHTTPRequest(c)

	engine.pool.Put(c)
}

// HandleContext re-enter a context that has been rewritten.
// This can be done by setting c.Request.Path to your new target.
// Disclaimer: You can loop yourself to death with this, use wisely.
func (engine *Engine) HandleContext(c *Context) {
	c.reset()
	engine.handleHTTPRequest(c)
	engine.pool.Put(c)
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	path := c.Request.URL.Path
	unescape := false
	if engine.UseRawPath && len(c.Request.URL.RawPath) > 0 {
		path = c.Request.URL.RawPath
		unescape = engine.UnescapePathValues
	}

	// Find root of the tree for the given HTTP method
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method == httpMethod {
			root := t[i].root
			// Find route in tree
			handlers, params, tsr := root.getValue(path, c.Params, unescape)
			if handlers != nil {
				c.handlers = handlers
				c.Params = params
				c.Next()
				c.writermem.WriteHeaderNow()
				return
			}
			if httpMethod != "CONNECT" && path != "/" {
				if tsr && engine.RedirectTrailingSlash {
					redirectTrailingSlash(c)
					return
				}
				if engine.RedirectFixedPath && redirectFixedPath(c, root, engine.RedirectFixedPath) {
					return
				}
			}
			break
		}
	}

	if engine.HandleMethodNotAllowed {
		for _, tree := range engine.trees {
			if tree.method != httpMethod {
				if handlers, _, _ := tree.root.getValue(path, nil, unescape); handlers != nil {
					c.handlers = engine.allNoMethod
					serveError(c, 405, default405Body)
					return
				}
			}
		}
	}
	c.handlers = engine.allNoRoute
	serveError(c, 404, default404Body)
}

var mimePlain = []string{MIMEPlain}

func serveError(c *Context, code int, defaultMessage []byte) {
	c.writermem.status = code
	c.Next()
	if !c.writermem.Written() {
		if c.writermem.Status() == code {
			c.writermem.Header()["Content-Type"] = mimePlain
			c.Writer.Write(defaultMessage)
		} else {
			c.writermem.WriteHeaderNow()
		}
	}
}

func redirectTrailingSlash(c *Context) {
	req := c.Request
	path := req.URL.Path
	code := 301 // Permanent redirect, request with GET method
	if req.Method != "GET" {
		code = 307
	}

	if length := len(path); length > 1 && path[length-1] == '/' {
		req.URL.Path = path[:length-1]
	} else {
		req.URL.Path = path + "/"
	}
	debugPrint("redirecting request %d: %s --> %s", code, path, req.URL.String())
	http.Redirect(c.Writer, req, req.URL.String(), code)
	c.writermem.WriteHeaderNow()
}

func redirectFixedPath(c *Context, root *node, trailingSlash bool) bool {
	req := c.Request
	path := req.URL.Path

	fixedPath, found := root.findCaseInsensitivePath(
		cleanPath(path),
		trailingSlash,
	)
	if found {
		code := 301 // Permanent redirect, request with GET method
		if req.Method != "GET" {
			code = 307
		}
		req.URL.Path = string(fixedPath)
		debugPrint("redirecting request %d: %s --> %s", code, path, req.URL.String())
		http.Redirect(c.Writer, req, req.URL.String(), code)
		c.writermem.WriteHeaderNow()
		return true
	}
	return false
}