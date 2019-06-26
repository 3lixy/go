package controllers

import (
	"azoya/nova"
	"tms_service/common"
	"tms_service/models"
)

//CustomsDeclarationTypeController 清关方式
type CustomsDeclarationTypeController struct {
}

//NewCustomsDeclarationTypeController 初始化
func NewCustomsDeclarationTypeController() *CustomsDeclarationTypeController {
	return &CustomsDeclarationTypeController{}
}

//List 返回物流商列表数据
func (customsDeclarationType *CustomsDeclarationTypeController) List(c *nova.Context) {
	customsDeclarationTypeModel := models.NewCustomsDeclarationTypeModel(c)

	customsDeclarationTypeList, err := customsDeclarationTypeModel.GetList()
	if err != nil {
		panic(err)
	}

	result := map[string]interface{}{"rows": customsDeclarationTypeList}
	common.ResponseResult(c, result, err)
}
