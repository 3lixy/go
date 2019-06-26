package controllers

import (
	"azoya/nova"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"tms_service/common"
	"tms_service/models"
	"tms_service/models/line"
	// "time"
	"encoding/json"
	"tms_service/validate"
)

//ShipmentsController 物流商controller
type ShipmentsController struct {
}

//NewShipmentsController 初始化
func NewShipmentsController() *ShipmentsController {
	return &ShipmentsController{}
}

//List 返回发货列表数据
func (Shipments *ShipmentsController) List(c *nova.Context) {
	ShipmentsModel := models.NewShipmentsModel(c)
	var page common.Page
	ShipmentsList, count, err := ShipmentsModel.GetList(c.Request.URL.Query(),page)
	if err != nil {
		panic(err)
	}

	result := map[string]interface{}{"rows": ShipmentsList, "total": count}
	common.ResponseResult(c, result, err)
}

//Create 创建发货单
func (Shipments *ShipmentsController) Create(c *nova.Context) {
	var detail models.Shipments
	bindErr := c.Bind(&detail)
	common.GetLogger().Info(fmt.Sprintf("shipment data:%v", detail))

	if bindErr != nil {
		panic(bindErr)
	}
	//清空warehouse_id，创建发货单时不写入仓库id，后续分配仓库再写入
	detail.WarehouseID = 0

	shipmentsModel := models.NewShipmentsModel(c)

	shipment, err := shipmentsModel.Create(detail)
	if err != nil {
		panic(err)
	}

	shipmentIds := []string{fmt.Sprintf("%d", shipment.EntityID)}

	//如果有仓库id，就直接分配仓库
	if warehouseID := c.PostForm("warehouse_id"); warehouseID != "" && warehouseID != "0" {
		_, err = shipmentsModel.AssignWarehouse(shipmentIds, warehouseID)

		//分配仓库成功则分配物流线路
		if err == nil {
			var params models.LineListQueryParams
			params.WebsiteID = fmt.Sprintf("%d", detail.WebsiteID)
			params.StoreID = fmt.Sprintf("%d", detail.StoreID)
			params.Status = fmt.Sprintf("%d", line.Enable)

			//如果能获取到物流线路，也直接分配物流线路
			lineModel := models.NewLineModel(c)
			lineList, _, getLineErr := lineModel.GetLineListWithoutPage(params)

			if getLineErr != nil {
				common.GetLogger().Error(fmt.Sprintf("line err:%v", getLineErr))
			} else if len(lineList) > 0 {
				var finalLine models.Line
				if len(lineList) == 1 {
					finalLine = lineList[0]
				} else {
					ruleDataJson := c.PostForm("rule_data")
					var ruleData models.RuleData
					var matchLineList []models.Line
					err = json.Unmarshal([]byte(ruleDataJson), &ruleData)
					common.GetLogger().Info(fmt.Sprintf("json err:%v", ruleDataJson))
					common.GetLogger().Info(fmt.Sprintf("rule data json:%v", err))
					common.GetLogger().Info(fmt.Sprintf("rule data:%v", ruleData))
					for _, l := range lineList {
						isMatch, _ := lineModel.IsMatchLine(l, ruleData)
						if isMatch {
							matchLineList = append(matchLineList, l)
						}
						common.GetLogger().Info(fmt.Sprintf("line data:%v,is_match:%v", l,isMatch))
					}
					if len(matchLineList) > 1 {
						finalLine = lineModel.GetAwardLine(matchLineList)
					} else if len(matchLineList) == 1 {
						finalLine = matchLineList[0]
					}
				}
				common.GetLogger().Info(fmt.Sprintf("final_line data:%v", finalLine))
				shipmentsModel.AssignTransport(shipmentIds, finalLine.EntityID)
			}
		}
	}

	shipmentDetail, _ := shipmentsModel.GetDetail(shipment.EntityID)
	common.ResponseResult(c, shipmentDetail, err)
}

//AssignWarehouse 给发货单分配仓库
func (Shipments *ShipmentsController) AssignWarehouse(c *nova.Context) {
	shipmentIDs := strings.Split(c.PostForm("shipment_ids"), models.ShipmentIDSplit)
	warehouseID := c.PostForm("warehouse_id")

	validateShipmentID(shipmentIDs)

	if warehouseID == "" {
		errText := "warehouse id is empty, please check"
		panic(errText)
	}

	shipmentsModel := models.NewShipmentsModel(c)

	status, err := shipmentsModel.AssignWarehouse(shipmentIDs, warehouseID)

	if err != nil {
		c.JSON(http.StatusOK, common.ErrorMessage(status, err.Error()))
		return
	}

	common.ResponseResult(c, nil, err)
}

//AssignTransport 给发货单分配物流
func (Shipments *ShipmentsController) AssignTransport(c *nova.Context) {
	shipmentIDs := strings.Split(c.PostForm("shipment_ids"), models.ShipmentIDSplit)
	transportLineID := c.PostForm("transport_line_id")
	validateShipmentID(shipmentIDs)

	if transportLineID == "" {
		errText := "transport line id is empty, please check"
		panic(errText)
	}

	shipmentsModel := models.NewShipmentsModel(c)
	lineID, _ := strconv.ParseUint(transportLineID, 0, 64)
	lineDetail, status, err := shipmentsModel.AssignTransport(shipmentIDs, lineID)

	if err != nil {
		c.JSON(http.StatusOK, common.ErrorMessage(status, err.Error()))
		return
	}

	common.ResponseResult(c, lineDetail, err)
}

//RemoveTransportLine 清除物流线路分配
func (Shipments *ShipmentsController) RemoveTransportLine(c *nova.Context) {
	shipmentIDs := strings.Split(c.PostForm("shipment_ids"), models.ShipmentIDSplit)
	validateShipmentID(shipmentIDs)

	shipmentsModel := models.NewShipmentsModel(c)
	status, err := shipmentsModel.RemoveTransportLine(shipmentIDs)

	if err != nil {
		c.JSON(http.StatusOK, common.ErrorMessage(status, err.Error()))
		return
	}

	common.ResponseResult(c, nil, err)
}

func validateShipmentID(shipments []string) {
	if len(shipments) == 1 && shipments[0] == "" {
		errText := "shipment ids is empty, please check"
		panic(errText)
	}
}

func (Shipments *ShipmentsController) RemoveWarehouse(c *nova.Context) {
	shipmentIDs := strings.Split(c.PostForm("shipment_ids"), models.ShipmentIDSplit)
	validateShipmentID(shipmentIDs)
	shipmentsModel := models.NewShipmentsModel(c)
	status, err := shipmentsModel.RemoveWarehouse(shipmentIDs)
	if err != nil {
		c.JSON(http.StatusOK, common.ErrorMessage(status, err.Error()))
		return
	}
	common.ResponseResult(c, nil, err)
}

func (Shipments *ShipmentsController) IsSetWareHouse(c *nova.Context) {
	shipmentIDs := strings.Split(c.Query("shipment_ids"), models.ShipmentIDSplit)
	validateShipmentID(shipmentIDs)
	shipmentsModel := models.NewShipmentsModel(c)
	status, err := shipmentsModel.IsSetWarehouse(shipmentIDs)
	if err != nil {
		c.JSON(http.StatusOK, common.ErrorMessage(status, err.Error()))
		return
	}
	common.ResponseResult(c, nil, err)
}

func (Shipments *ShipmentsController) Details(c *nova.Context) {
	smodel := models.NewShipmentsModel(c)

	result, err := smodel.GetShipments(c.Request.URL.Query())
	common.ResponseResult(c, result, err)
}

func (Shipments *ShipmentsController) GetLabelDownUrl(c *nova.Context) {
	var query validate.LabelUrlQuery
	err := c.Bind(&query)
	if err != nil {
		common.ResponseResult(c, nil, err)
		return
	}

	websiteId:= c.Request.URL.Query().Get("website_id")
	incrementId:= c.Request.URL.Query().Get("increment_id")
	shipmentsModel := models.NewShipmentsModel(c)
	result, err := shipmentsModel.GetShipments(map[string][]string{"website_id":{websiteId},"order_increment_id":{incrementId}})
	if err != nil {
		common.ResponseResult(c, nil, err)
		return
	}
	config := common.GetConfig()
	url:=config.String("url::label_down_url")
	common.ResponseResult(c, map[string]string{"url":url+result.LabelIdentification}, err)
	return
}