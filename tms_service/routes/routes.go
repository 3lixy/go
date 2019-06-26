package routes

import (
	"azoya/nova"
	"tms_service/controllers"
)

//Register 注册路由
func Register(router *nova.Engine) {
	logistics := controllers.NewLogisticsController()
	line := controllers.NewLineController()
	subscribe := controllers.NewSubscribeTypeController()
	customs := controllers.NewCustomsDeclarationTypeController()
	shipments := controllers.NewShipmentsController()
	shipmentTrack := controllers.NewShipmentTrackController()
	logisticsTrack := controllers.NewLogisticsTrackController()

	//物流商
	r := router.Group("/logistics")
	{
		r.GET("/list", logistics.List)
		r.GET("/detail", logistics.Detail)
		r.POST("/delete", logistics.Delete)
		r.POST("/update", logistics.Update)
		r.POST("/create", logistics.Create)
	}

	//物流路线
	r = router.Group("/line")
	{
		r.GET("/list", line.List)
		r.GET("/all", line.All)
		r.GET("/detail", line.Detail)
		r.POST("/delete", line.Delete)
		r.POST("/update", line.Update)
		r.POST("/create", line.Create)
	}

	//api订阅方式
	r = router.Group("/subscribe")
	{
		r.GET("/list", subscribe.List)
	}
	//海关清关方式
	r = router.Group("/customs")
	{
		r.GET("/list", customs.List)
	}

	//海关清关方式
	r = router.Group("/shipments")
	{
		r.GET("/list", shipments.List)
		r.POST("/create", shipments.Create)
		r.POST("/assignwarehouse", shipments.AssignWarehouse)
		r.POST("/assigntransport", shipments.AssignTransport)
		r.POST("/removetransportline", shipments.RemoveTransportLine)
		r.POST("/remove_warehouse", shipments.RemoveWarehouse)
		r.GET("/is_set_warehouse", shipments.IsSetWareHouse)
		r.GET("/detail", shipments.Details)
		r.GET("/down_url", shipments.GetLabelDownUrl)
	}

	//运单管理
	r = router.Group("/shipment_track")
	{
		r.POST("/create", shipmentTrack.Add)
		r.POST("/delete", shipmentTrack.Delete)
		r.GET("/list", shipmentTrack.List)
		r.POST("/add", shipmentTrack.Add)
	}

	//物流轨迹
	r = router.Group("/logistics_track")
	{
		r.GET("/track_items", logisticsTrack.TrackItems)
	}
}
