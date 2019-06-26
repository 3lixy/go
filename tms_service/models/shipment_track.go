package models

import (
	"azoya/nova"
	"fmt"
	"strconv"
	"strings"
	"tms_service/common"
)

type ShipmentTrackModel struct {
	Context *nova.Context
}

func NewShipmentTrackModel(c *nova.Context) *ShipmentTrackModel {
	return &ShipmentTrackModel{c}
}

func (s *ShipmentTrackModel) TableName() string {
	return "shipment_track"
}

type ShipmentTrack struct {
	EntityID           uint64  `gorm:"primary_key" json:"entity_id" form:"entity_id"`
	WebsiteID          uint64  `gorm:"website_id" json:"website_id" form:"website_id" binding:"required"`
	OrderID            uint64  `gorm:"order_id" json:"order_id" form:"order_id" binding:"required"`
	IncrementID        string  `gorm:"increment_id" json:"increment_id" form:"increment_id" binding:"required"`
	PartnerIncrementId string  `gorm:"partner_increment_id" json:"partner_increment_id" form:"partner_increment_id"`
	OrderShipmentNo    string  `gorm:"order_shipment_no" json:"order_shipment_no" form:"order_shipment_no"`
	SupplierOrderNo    string  `gorm:"supplier_order_no" json:"supplier_order_no" form:"supplier_order_no"`
	CarrierCode        string  `gorm:"carrier_code" json:"carrier_code" form:"carrier_code" binding:"required"`
	TrackNumber        string  `gorm:"track_number" json:"track_number" form:"track_number" binding:"required"`
	VendorID           uint64  `gorm:"vendor_id" json:"vendor_id" form:"vendor_id"`
	StockID            uint64  `gorm:"stock_id" json:"stock_id" form:"stock_id"`
	Length             float64 `gorm:"length" json:"length" form:"length"`
	Width              float64 `gorm:"width" json:"width" form:"width"`
	Height             float64 `gorm:"height" json:"height" form:"height"`
	Weight             float64 `gorm:"weight" json:"weight" form:"weight"`
	CreatedAt          string  `gorm:"created_at" json:"created_at" form:"created_at"`
	UpdatedAt          string  `gorm:"updated_at" json:"updated_at" form:"updated_at"`
}

//Create
func (s *ShipmentTrackModel) Create(st ShipmentTrack) (error) {
	time := common.DefaultTime()
	st.CreatedAt = time
	st.UpdatedAt = time

	err := common.GetDb().TmsWriteDb.Table(s.TableName()).Create(&st).Error

	if err != nil {
		s.Context.Logger().Error(err.Error())
	}
	return err
}

//GetList
func (s *ShipmentTrackModel) GetList(params map[string][]string, page common.Page) ([]ShipmentTrack, uint64, error) {
	sort := "entity_id"
	orderBy := "DESC"

	where := make(map[string]interface{})
	if _, ok := params["website_id"]; ok {
		websiteID, _ := strconv.ParseUint(params["website_id"][0], 10, 64)
		if websiteID > 0 {
			where["website_id"] = websiteID
		}
	}
	if _, ok := params["vendor_id"]; ok {
		vendorID, _ := strconv.ParseUint(params["vendor_id"][0], 10, 64)
		where["vendor_id"] = vendorID
	}
	if _, ok := params["stock_id"]; ok {
		stockID, _ := strconv.ParseUint(params["stock_id"][0], 10, 64)
		where["stock_id"] = stockID
	}
	var trackNumber string
	var partnerIncrementID string
	var orderShipmentNo string
	var incrementID string
	if _, ok := params["carrier_code"]; ok {
		carrierCode := params["carrier_code"][0]
		where["carrier_code"] = carrierCode
	}
	if _, ok := params["track_number"]; ok {
		trackNumber = params["track_number"][0]
	}
	if _, ok := params["partner_increment_id"]; ok {
		partnerIncrementID = params["partner_increment_id"][0]
	}
	if _, ok := params["order_shipment_no"]; ok {
		orderShipmentNo = params["order_shipment_no"][0]
	}
	if _, ok := params["increment_id"]; ok {
		incrementID = params["increment_id"][0]
	}
	var createdAtStart string
	var createdAtEnd string
	if _, ok := params["created_at_start"]; ok {
		createdAtStart = params["created_at_start"][0]
	}
	if _, ok := params["created_at_end"]; ok {
		createdAtEnd = params["created_at_end"][0]
	}

	var whereSecond string
	if trackNumber != "" {
		trackNumberS := strings.Split(trackNumber, ",")
		trackNumberStr := "'" + strings.Join(trackNumberS, "','") + "'"
		whereSecond += fmt.Sprintf(" and track_number in(%s)", trackNumberStr)
	}
	if partnerIncrementID != "" {
		partnerIncrementIDS := strings.Split(partnerIncrementID, ",")
		partnerIncrementIDStr := "'" + strings.Join(partnerIncrementIDS, "','") + "'"
		whereSecond += fmt.Sprintf(" and partner_increment_id in(%s)", partnerIncrementIDStr)
	}
	if orderShipmentNo != "" {
		orderShipmentNoS := strings.Split(orderShipmentNo, ",")
		orderShipmentNoStr := "'" + strings.Join(orderShipmentNoS, "','") + "'"
		whereSecond += fmt.Sprintf(" and order_shipment_no  in(%s)", orderShipmentNoStr)
	}
	if incrementID != "" {
		incrementIDS := strings.Split(incrementID, ",")
		incrementIDStr := "'" + strings.Join(incrementIDS, "','") + "'"
		whereSecond += fmt.Sprintf(" and increment_id in (%s) ", incrementIDStr)
	}
	if createdAtStart != "" {
		whereSecond += fmt.Sprintf(" and created_at >= %q ", createdAtStart)
	}
	if createdAtEnd != "" {
		whereSecond += fmt.Sprintf(" and created_at <= %q ", createdAtEnd)
	}

	wherestr := ""
	if whereSecond != "" {
		wherestr = "1 = 1" + whereSecond
	}
	ss := common.GetParamsPage(page)
	var result []ShipmentTrack
	err := common.GetDb().TmsReadDb.
		Table(s.TableName()).
		Where(where).
		Where(wherestr).
		Offset(ss.Page).
		Limit(ss.Limit).
		Order(fmt.Sprintf("%s  %s", sort, orderBy)).
		Find(&result).Error

	var count uint64
	common.GetDb().TmsReadDb.
		Table(s.TableName()).
		Where(where).
		Where(wherestr).
		Count(&count)

	return result, count, err
}

func (s *ShipmentTrackModel) GetListByID(shipmentTrackIds []string) ([]ShipmentTrack, error) {
	var list []ShipmentTrack
	err := common.GetDb().TmsReadDb.Table(s.TableName()).
		Where("entity_id in (?)", shipmentTrackIds).
		Find(&list).
		Error

	return list, err
}

func (s *ShipmentTrackModel) Delete(shipmentTrackIds []string) (error) {
	list, err := s.GetListByID(shipmentTrackIds)
	if err != nil {
		return err
	}
	subscribeRecordModel := NewSubscribeRecordModel(s.Context)
	trackItemModel := NewTrackItemModel(s.Context)

	for _, track := range list {
		err = subscribeRecordModel.Delete(map[string]interface{}{"channel_id": track.WebsiteID, "order_id": track.OrderID, "track_number": track.TrackNumber})
		if err != nil {
			return err
		}
		err = trackItemModel.Delete(track.WebsiteID, map[string]interface{}{"order_id": track.OrderID, "track_number": track.TrackNumber})
		if err != nil {
			return err
		}
	}

	var whereStr string
	shipmentTrackIdStr := strings.Join(shipmentTrackIds, ",")
	whereStr = fmt.Sprintf("entity_id in (%s) ", shipmentTrackIdStr)
	err = common.GetDb().TmsWriteDb.
		Table(s.TableName()).
		Where(whereStr).
		Delete(ShipmentTrack{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *ShipmentTrackModel) GetShipmentTrack(params map[string]interface{}) (shipmentTrack ShipmentTrack, err error) {
	err = common.GetDb().TmsReadDb.Table(s.TableName()).
		Where(params).
		First(&shipmentTrack).Error
	return
}
