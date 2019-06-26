package controllers

import (
	"azoya/nova"
	"net/http"
	"tms_service/common"
	"tms_service/service"
	"tms_service/validate"
)

type LogisticsTrackController struct {
	BaseControllers
}

func NewLogisticsTrackController() *LogisticsTrackController {
	return &LogisticsTrackController{}
}

func (t *LogisticsTrackController) TrackItems(c *nova.Context) {
	var query validate.TrackItemListQuery
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		c.Abort()
		return
	}

	logisticsTrackService := service.NewLogisticsTrackService(c)
	stock := logisticsTrackService.GetTrackItemList(query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, common.Message(http.StatusOK, "success", stock))
}
