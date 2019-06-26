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