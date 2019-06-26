package routers

import (
	"azoya/nova"
	"config_service/controllers"
)

//RegisterRoutes 注册路由
func RegisterRoutes(router *nova.Engine) {
	stock := controllers.NewStockControllers()
	store := controllers.NewStoreControllers()
	system := controllers.NewSystemControllers()
	stockRouter := router.Group("/stock")
	{
		stockRouter.POST("/add", stock.Add)
		stockRouter.POST("/update", stock.Update)
		stockRouter.GET("/list", stock.List)
		stockRouter.GET("/detail", stock.Detail)
	}

	storeRouter := router.Group("/store")
	{
		storeRouter.GET("/detail", store.Detail)
		storeRouter.GET("/list", store.List)
	}

	systemRouter := router.Group("/system")
	{
		systemRouter.POST("/update", system.Update)
		systemRouter.GET("/list", system.List)
		systemRouter.GET("/detail", system.Detail)
	}
}
