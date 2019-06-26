package controllers

import (
	"azoya/nova"
	"config_service/common"
	"config_service/service"
	"config_service/validate"
	"net/http"
	"strconv"
)

type SystemControllers struct {
	BaseControllers
}

//NewSystemControllers 模型
func NewSystemControllers() *SystemControllers {
	return &SystemControllers{}
}

// 更新系统初始化状态
func (s *SystemControllers) Update(c *nova.Context) {
	var query validate.UpdateSystemQuery

	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		c.Abort()
		return
	}

	systemService := service.NewSystemService(c)
	err = systemService.Update(query)
	if err != nil {
		c.JSON(http.StatusOK, common.Message(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, common.Message(http.StatusOK, "success", err))
}

// 获取体统初始化状态列表
func (s *SystemControllers) List(c *nova.Context) {
	systemService := service.NewSystemService(c)
	result := systemService.GetSystemList(c.Request.URL.Query())
	c.JSON(http.StatusOK, common.Message(http.StatusOK, common.MsgSuccess, result))
}

//Detail 获取system init的详情
func (s *SystemControllers) Detail(c *nova.Context) {
	entityID, err := strconv.Atoi(c.Query("entity_id"))
	if err != nil {
		panic(err)
	}
	systemService := service.NewSystemService(c)
	result := systemService.GetDetail(entityID)
	c.JSON(http.StatusOK, common.Message(http.StatusOK, common.MsgSuccess, result))
}
