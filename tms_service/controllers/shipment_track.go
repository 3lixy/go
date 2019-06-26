package controllers

import (
	"azoya/nova"
	"tms_service/common"
	"tms_service/models"
	"net/http"
	"tms_service/validate"
	"tms_service/service"
)

//运单管理 STController
type ShipmentTrackController struct {
	BaseControllers
}

func NewShipmentTrackController() *ShipmentTrackController {
	return &ShipmentTrackController{}
}

//Create
func (s *ShipmentTrackController) Add(c *nova.Context) {
	var query validate.AddShipmentTrackQuery
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		return
	}
	shipmentTrackService := service.NewShipmentTrackService(c)
	err = shipmentTrackService.Add(query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, common.Message(http.StatusOK, "success", err))

}

//List
func (s *ShipmentTrackController) List(c *nova.Context) {
	stModel := models.NewShipmentTrackModel(c)
	var page common.Page
	stLine, count, err := stModel.GetList(c.Request.URL.Query(), page)
	if err != nil {
		panic(err)
	}
	result := map[string]interface{}{"rows": stLine, "total": count}
	common.ResponseResult(c, result, err)
}

//Delete
func (s *ShipmentTrackController) Delete(c *nova.Context) {
	var query validate.DeleteShipmentTrackQuery
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		return
	}
	shipmentTrackService := service.NewShipmentTrackService(c)
	err = shipmentTrackService.Delete(query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, common.Message(http.StatusOK, "success", err))
}



