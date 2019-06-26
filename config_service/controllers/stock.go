package controllers

import (
	"azoya/nova"
	"config_service/common"
	"config_service/service"
	"config_service/validate"
	"net/http"
)

type StockControllers struct {
	BaseControllers
}

//NewStockControllers 模型
func NewStockControllers() *StockControllers {
	return &StockControllers{}
}

//Add 仓库新增接口
func (s *StockControllers) Add(c *nova.Context) {

	var query validate.AddQuery

	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		c.Abort()
		return
	}

	stockService := service.NewStockService(c)
	err = stockService.Add(query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, common.Message(http.StatusOK, "success", err))
}

//Update 仓库更新接口
func (s *StockControllers) Update(c *nova.Context) {

	var query validate.UpdateQuery

	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		c.Abort()
		return
	}

	stockService := service.NewStockService(c)
	err = stockService.Update(query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, common.Message(http.StatusOK, "success", err))
}

//List 仓库列表
func (s *StockControllers) List(c *nova.Context) {
	stockService := service.NewStockService(c)
	result := stockService.GetStockList(c.Request.URL.Query())
	c.JSON(http.StatusOK, common.Message(http.StatusOK, common.MsgSuccess, result))
}

//GetDetail 获取仓库详情接口
func (s *StockControllers) Detail(c *nova.Context) {
	var query validate.DetailQuery

	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		c.Abort()
		return
	}

	stockService := service.NewStockService(c)
	stock, err := stockService.GetDetail(query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, common.Message(http.StatusOK, "success", stock))
}
