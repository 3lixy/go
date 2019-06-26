package controllers

import (
	"azoya/nova"
	"errors"
	"strconv"
	"tms_service/common"
	"tms_service/models"
	// "fmt"
)

//LogisticsController 物流商controller
type LogisticsController struct {
}

//NewLogisticsController 初始化
func NewLogisticsController() *LogisticsController {
	return &LogisticsController{}
}

//List 返回物流商列表数据
func (logistics *LogisticsController) List(c *nova.Context) {
	logisticsModel := models.NewLogisticsModel(c)

	var params models.ListQueryParams
	err := c.Bind(&params)
	if err != nil {
		panic(err)
	}

	logisticsList, count, err := logisticsModel.GetList(params)
	if err != nil {
		panic(err)
	}

	result := map[string]interface{}{"rows": logisticsList, "total": count}

	common.ResponseResult(c, result, err)
}

//Detail 返回物流商详情
func (logistics *LogisticsController) Detail(c *nova.Context) {
	logisticsID, err := strconv.ParseUint(c.Query("entity_id"), 0, 64)
	if err != nil {
		panic(err)
	}

	logisticsModel := models.NewLogisticsModel(c)

	detail, err := logisticsModel.GetDetail(logisticsID)

	common.ResponseResult(c, detail, err)
}

//Delete 删除物流商
func (logistics *LogisticsController) Delete(c *nova.Context) {
	logisticsID, err := strconv.ParseUint(c.PostForm("entity_id"), 0, 64)
	if err != nil {
		panic(err)
	}

	lineModel := models.NewLineModel(c)

	_, count, err := lineModel.GetLineWithLogisticsID(logisticsID)
	if count > 0 {
		getLineErr := errors.New("this logistics already assign transport line")
		common.ResponseResult(c, nil, getLineErr)
		return
	}

	logisticsModel := models.NewLogisticsModel(c)

	deleteErr := logisticsModel.Delete(logisticsID)

	common.ResponseResult(c, nil, deleteErr)
}

//Update 更新物流商数据
func (logistics *LogisticsController) Update(c *nova.Context) {
	logisticsID, err := strconv.ParseUint(c.PostForm("entity_id"), 0, 64)

	if err != nil {
		panic(err)
	}

	logisticsModel := models.NewLogisticsModel(c)
	_, err = logisticsModel.GetDetail(logisticsID)

	if err != nil {
		panic(err)
	}

	var updateLogistics models.Logistics
	bindErr := c.Bind(&updateLogistics)
	if bindErr != nil {
		panic(bindErr)
	}

	updateLogistics, updateErr := logisticsModel.Update(updateLogistics)

	common.ResponseResult(c, updateLogistics, updateErr)
}

//Create 添加物流商数据
func (logistics *LogisticsController) Create(c *nova.Context) {
	var detail models.Logistics
	bindErr := c.Bind(&detail)

	if bindErr != nil {
		panic(bindErr)
	}

	logisticsModel := models.NewLogisticsModel(c)
	newLogistics, createErr := logisticsModel.Create(detail)

	common.ResponseResult(c, newLogistics, createErr)
}
