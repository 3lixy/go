package controllers

import (
	"azoya/nova"
	"config_service/common"
	"config_service/service"
	"config_service/validate"
	"net/http"
)

type StoreControllers struct {
	BaseControllers
}

//NewStoreControllers 模型
func NewStoreControllers() *StoreControllers {
	return &StoreControllers{}
}

//Detail 获取店铺详情接口
func (s *StoreControllers) Detail(c *nova.Context) {
	var query validate.StoreDetailQuery

	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		c.Abort()
		return
	}

	storeService := service.NewStoreService(c)
	store, err := storeService.Detail(query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, common.Message(http.StatusOK, "success", store))
}

func (s *StoreControllers) List(c *nova.Context) {
	storeService := service.NewStoreService(c)
	result := storeService.GetStoreList(c.Request.URL.Query())
	c.JSON(http.StatusOK, common.Message(http.StatusOK, common.MsgSuccess, result))
}
