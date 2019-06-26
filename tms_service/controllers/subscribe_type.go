package controllers

import (
	"azoya/nova"
	"tms_service/common"
	"tms_service/models"
)

//SubscribeTypeController 物流商controller
type SubscribeTypeController struct {
}

//NewSubscribeTypeController 初始化
func NewSubscribeTypeController() *SubscribeTypeController {
	return &SubscribeTypeController{}
}

//List 返回物流商列表数据
func (SubscribeType *SubscribeTypeController) List(c *nova.Context) {
	SubscribeTypeModel := models.NewSubscribeTypeModel(c)

	SubscribeTypeList, err := SubscribeTypeModel.GetList()
	if err != nil {
		panic(err)
	}

	result := map[string]interface{}{"rows": SubscribeTypeList}
	common.ResponseResult(c, result, err)
}
