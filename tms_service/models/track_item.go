package models

import (
	"azoya/nova"
	"fmt"
	"strconv"
	"tms_service/common"
)

//TrackItemModels 模型
type TrackItemModel struct {
	BaseModel
}

//NewTrackItemModels 模型
func NewTrackItemModel(c *nova.Context) *TrackItemModel {
	return &TrackItemModel{BaseModel{}}
}

type TrackItem struct {
	EntityID      uint64 `gorm:"primary_key" json:"entity_id"`
	OrderID       uint64 `gorm:"column:order_id" json:"order_id"`
	TrackNumber   string `gorm:"column:track_number" json:"track_number"`
	CarriersCode  string `gorm:"column:carriers_code" json:"carriers_code"`
	Country       string `gorm:"column:country" json:"country"`
	State         string `gorm:"column:state" json:"state"`
	StatusCode    string `gorm:"column:status_code" json:"status_code"`
	Address       string `gorm:"column:address" json:"address"`
	CompanyName   string `gorm:"column:company_name" json:"company_name"`
	EventCode     string `gorm:"column:event_code" json:"event_code"`
	EventDateTime string `gorm:"column:event_date_time" json:"event_date_time"`
	Description   string `gorm:"column:description" json:"description"`
	IsNum         uint64 `gorm:"column:is_num" json:"is_num"`
	CreatedAt     string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     string `gorm:"column:updated_at" json:"updated_at"`
}

func (t *TrackItem) TableName(WebsiteID uint64) string {
	tableName := "logistics_track_item_" + strconv.Itoa(int(WebsiteID))
	return tableName
}

func (i *TrackItemModel) GetItems(websiteID uint64, params map[string]interface{}) (items []TrackItem, err error) {
	sort := "event_date_time"
	orderBy := "desc"
	var result TrackItem
	err = common.GetDb().TrackReadDb.Table(result.TableName(websiteID)).
		Where(params).
		Order(fmt.Sprintf("%s  %s", sort, orderBy)).
		Find(&items).Error
	return
}

func (i *TrackItemModel) Delete(websiteID uint64,where map[string]interface{}) (err error) {
	var trackItem TrackItem
	db := common.GetDb().TrackWriteDb
	err = db.Table(trackItem.TableName(websiteID)).Where(where).Delete(TrackItem{}).Error
	return
}
