package main

import (
	"azoya/nova"
	"fmt"
	"tms_service/common"
	"tms_service/routes"
	"tms_service/middleware"
)

func main() {
	router := nova.Default()
	nova.SetMode(router.Configer.String("service::runmode"))

	common.Init(router.Configer)
	//捕获异常
	router.Use(middleware.RecoveryMiddleware())
	//log需要依赖配置文件
	common.InitLogger()
	defer common.SyncLogger()

	//初始化db
	common.InitDb()

	//注册路由
	routes.Register(router)

	addr := "0.0.0.0"

	port := router.Configer.String("listen::port")

	if port != "" {
		addr = fmt.Sprintf("%s:%s", addr, port)
	}

	router.Run(addr)
}
