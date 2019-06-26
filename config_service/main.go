package main

import (
	"azoya/nova"
	"config_service/common"        //公共模块
	m "config_service/middlewares" //中间件
	r "config_service/routers"     //路由配置
	//"config_service/service"       //业务层
	"net/http"
	"time"
)

func main() {
	//框架初始化
	router := nova.Default()
	//捕获异常
	router.Use(m.RecoveryMiddleware())
	//将配置文件载入内存
	common.Init(router.Configer)
	//运行模式
	runmode := router.Configer.String("service::runmode")
	if runmode != "" {
		nova.SetMode(runmode)
	}

	//初始化db
	common.InitDb()

	//注册路由
	r.RegisterRoutes(router)
	//运行端口
	port := router.Configer.String("listen::port")
	//启动
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
