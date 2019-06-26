package controllers

import (
	"azoya/nova"
	"errors"
	"strconv"
	"tms_service/common"
	"tms_service/models"
)

//LineController 物流商controller
type LineController struct {
}

//NewLineController 初始化
func NewLineController() *LineController {
	return &LineController{}
}

//List 返回物流线路
func (line *LineController) List(c *nova.Context) {
	lineModel := models.NewLineModel(c)

	var params models.LineListQueryParams
	var page common.Page
	err := c.Bind(&params)
	pageErr := c.Bind(&page)
	if err != nil && pageErr != nil {
		panic(err)
	}

	lineList, count, err := lineModel.GetList(params, page)
	if err != nil {
		panic(err)
	}

	result := map[string]interface{}{"rows": lineList, "total": count}
	common.ResponseResult(c, result, err)
}

//All 返回所有物流线路
func (line *LineController) All(c *nova.Context) {
	lineModel := models.NewLineModel(c)

	lineList, count, err := lineModel.GetAll()
	if err != nil {
		panic(err)
	}

	result := map[string]interface{}{"rows": lineList, "total": count}
	common.ResponseResult(c, result, err)
}

//Detail 返回物流线路详情
func (line *LineController) Detail(c *nova.Context) {
	lineID, err := strconv.ParseUint(c.Query("entity_id"), 0, 64)
	if err != nil {
		panic(err)
	}

	lineModel := models.NewLineModel(c)

	detail, err := lineModel.GetDetail(lineID)

	common.ResponseResult(c, detail, err)
}

//Delete 删除物流线路
func (line *LineController) Delete(c *nova.Context) {
	lineID, err := strconv.ParseUint(c.PostForm("entity_id"), 0, 64)
	if err != nil {
		panic(err)
	}

	shipmentModel := models.NewShipmentsModel(c)
	var shipmentQueryParams models.ShipmentsListQueryParams
	// shipmentQueryParams := new models.ShipmentsListQueryParams{TransportLineID:lineID}
	shipmentQueryParams.TransportLineID = lineID
	count, _ := shipmentModel.GetTotalRowsCount(shipmentQueryParams)
	if count > 0 {
		err = errors.New("Can't delete the transport line, because already assign shipments")
		common.ResponseResult(c, "Delete Failed", err)
		return
	}

	lineModel := models.NewLineModel(c)

	deleteErr := lineModel.Delete(lineID)

	common.ResponseResult(c, "Delete Success", deleteErr)
}

//Update 更新物流线路信息
func (line *LineController) Update(c *nova.Context) {
	lineID, err := strconv.ParseUint(c.PostForm("entity_id"), 0, 64)

	if err != nil {
		panic(err)
	}

	lineModel := models.NewLineModel(c)
	detail, err := lineModel.GetDetail(lineID)

	if err != nil {
		panic(err)
	}

	bindErr := c.Bind(&detail)
	if bindErr != nil {
		panic(bindErr)
	}

	detail, updateErr := lineModel.Update(detail)

	common.ResponseResult(c, detail, updateErr)
}

//Create 添加物流线路
func (line *LineController) Create(c *nova.Context) {
	var detail models.Line
	bindErr := c.Bind(&detail)

	if bindErr != nil {
		panic(bindErr)
	}

	lineModel := models.NewLineModel(c)
	newLine, createErr := lineModel.Create(detail)

	common.ResponseResult(c, newLine, createErr)
}
