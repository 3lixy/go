package service

import (
	"azoya/nova"
	"tms_service/validate"
	"tms_service/models"
	"tms_service/common"
	"strings"
	"strconv"
)

//ShipmentTrackService 模型
type ShipmentTrackService struct {
	BaseService
}

//NewShipmentTrackService 模型
func NewShipmentTrackService(c *nova.Context) *ShipmentTrackService {
	return &ShipmentTrackService{BaseService{C: c}}
}

func (s *ShipmentTrackService) Add(query validate.AddShipmentTrackQuery) (err error) {
	orderService := NewOrderService(s.C)
	shipmentsModel := models.NewShipmentsModel(s.C)
	logisticsModel := models.NewLogisticsModel(s.C)
	shipmentTrackModel := models.NewShipmentTrackModel(s.C)
	subscribeRecordModel := models.NewSubscribeRecordModel(s.C)
	orderList, err := orderService.OrderQuickSearch(query.WebsiteID, map[string]interface{}{"increment_id": query.IncrementID}, 0, 1)
	if len(orderList) > 0 { //判断订单是否存在
		order := orderList[0]
		orderShipment, _ := shipmentsModel.GetShipment(map[string]interface{}{"website_id": query.WebsiteID, "order_id": order.EntityID})
		if orderShipment.EntityID <= 0 { //判断发货单是否存在
			err = common.ErrOrderShipmentNotExist
			return err
		}
		logistics, _ := logisticsModel.GetLogistics(map[string]interface{}{"code": query.CarrierCode})
		if logistics.EntityID <= 0 { //判断物流方式是否存在
			err = common.ErrLogisticsNotExist
			return err
		}
		//写入运单号
		var shipmentTrack models.ShipmentTrack
		shipmentTrack, _ = shipmentTrackModel.GetShipmentTrack(map[string]interface{}{"website_id": query.WebsiteID, "order_id": order.EntityID, "carrier_code": query.CarrierCode, "track_number": query.TrackNumber})
		if shipmentTrack.EntityID <= 0 {
			shipmentTrack.WebsiteID = query.WebsiteID
			shipmentTrack.OrderID = order.EntityID
			shipmentTrack.IncrementID = order.IncrementID
			shipmentTrack.SupplierOrderNo = order.SupplierOrderNo
			shipmentTrack.OrderShipmentNo = orderShipment.IncrementID
			shipmentTrack.PartnerIncrementId = order.ThreePartOrderNumber
			shipmentTrack.VendorID = order.VendorID
			uStockID, _ := strconv.ParseUint(order.StockID, 10, 64)
			shipmentTrack.StockID = uStockID
			shipmentTrack.CarrierCode = query.CarrierCode
			shipmentTrack.TrackNumber = query.TrackNumber
			shipmentTrack.Length = query.Length
			shipmentTrack.Height = query.Height
			shipmentTrack.Weight = query.Weight
			shipmentTrack.Width = query.Width
			err = shipmentTrackModel.Create(shipmentTrack)
			if err != nil {
				return err
			}
		}
		//写入物流轨迹记录
		var subscribeRecord models.SubscribeRecord
		subscribeRecord, _ = subscribeRecordModel.GetSubscribeRecord(map[string]interface{}{"channel_id": query.WebsiteID, "order_id": order.EntityID, "carriers_code": logistics.CarrierCode, "track_number": query.TrackNumber})
		if subscribeRecord.EntityID <= 0 {
			subscribeRecord.OrderID = order.EntityID
			subscribeRecord.IncrementID = query.IncrementID
			subscribeRecord.TrackNumber = query.TrackNumber
			subscribeRecord.ChannelID = query.WebsiteID
			subscribeRecord.CarriersCode = logistics.CarrierCode
			subscribeRecord.CarriersTitle = logistics.Name
			if logistics.SubscribeTypeID > 1 { //判断是否点阅物流轨迹
				subscribeRecord.IsSubscribe = 1
			}
			err = subscribeRecordModel.Add(subscribeRecord)
			if err != nil {
				return err
			}
		}
		//已付款订单发起订单状态变更已发货
		if order.Status == PROCESSINGSTATUS {
			err = orderService.OrderStatusToShipment(query.WebsiteID, order.EntityID)
			if err != nil {
				return err
			}
			//发送邮件短信
		}
	} else {
		err = common.ErrGetOrder
	}
	return
}

func (s *ShipmentTrackService) Delete(query validate.DeleteShipmentTrackQuery) (err error) {
	shipmentTrackModel := models.NewShipmentTrackModel(s.C)
	entityIds := strings.Split(query.EntityID, ",")
	err = shipmentTrackModel.Delete(entityIds)
	return
}
