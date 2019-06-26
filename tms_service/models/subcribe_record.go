package models

import (
	"azoya/lib/log"
	"azoya/nova"
	"fmt"
	"strconv"
	"strings"
	"tms_service/common"
)

//SubscribeRecordModels 模型
type SubscribeRecordModel struct {
	BaseModel
}

//NewSubscribeRecordModels 模型
func NewSubscribeRecordModel(c *nova.Context) *SubscribeRecordModel {
	return &SubscribeRecordModel{BaseModel{}}
}

//TableName 表名
func (r *SubscribeRecord) TableName() string {
	return "logistics_subscribe_record"
}

//SubscribeRecord 订阅结果数据
type SubscribeRecord struct {
	EntityID        uint64 `gorm:"primary_key" json:"entity_id"`
	OrderID         uint64 `gorm:"column:order_id" json:"order_id"`
	TrackNumber     string `gorm:"column:track_number" json:"track_number"`
	Status          uint64 `gorm:"column:status" json:"status"`
	ChannelID       uint64 `gorm:"column:channel_id" json:"channel_id"`
	IncrementID     string `gorm:"column:increment_id" json:"increment_id"`
	CarriersCode    string `gorm:"column:carriers_code" json:"carriers_code"`
	CarriersTitle   string `gorm:"column:carriers_title" json:"carriers_title"`
	Address         string `gorm:"column:address" json:"address"`
	Description     string `gorm:"column:description" json:"description"`
	IsSubscribe     uint64 `gorm:"column:is_subscribe" json:"is_subscribe"`
	ReSubscribeTime uint64 `gorm:"column:re_subscribe_time" json:"re_subscribe_time"`
	PushStatus      uint64 `gorm:"column:push_status" json:"push_status"`
	PushTime        uint64 `gorm:"column:push_time" json:"push_time"`
	CreatedAt       string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       string `gorm:"column:updated_at" json:"updated_at"`
}

const (
	LOGISTIC_TRACK_FINISH uint64 = 100 //物流追踪结束
	SUBSCRIBE_SUCCESS     uint64 = 200 //物流信息订阅成功
	CALLBACK_FAILED       uint64 = 300 //回调数据失败
	SUBSCRIBE_OFF         uint64 = 400 //物流订阅未开启
)

//获取订阅记录列表
func (s *SubscribeRecordModel) GetList(params map[string][]string) (records []SubscribeRecord, err error) {
	sort := "created_at"
	orderBy := "desc"
	var channelID uint64
	var status uint64
	var trackNumber string
	var CarriersCode string
	var incrementID string
	if _, ok := params["channel_id"]; ok {
		channelID, _ = strconv.ParseUint(params["channel_id"][0], 10, 64)
	}
	if _, ok := params["status"]; ok {
		status, _ = strconv.ParseUint(params["status"][0], 10, 64)
	}
	if _, ok := params["increment_id"]; ok {
		incrementID = params["increment_id"][0]
	}
	if _, ok := params["track_number"]; ok {
		trackNumber = params["track_number"][0]
	}
	if _, ok := params["carriers_code"]; ok {
		CarriersCode = params["carriers_code"][0]
	}
	var where string
	if channelID > 0 {
		where += fmt.Sprintf(" and channel_id = %v", channelID)
	}
	if incrementID != "" {
		incrementIDS := strings.Split(incrementID, ",")
		incrementIDStr := "'" + strings.Join(incrementIDS, "','") + "'"
		where += fmt.Sprintf(" and increment_id in(%s)", incrementIDStr)
	}
	if trackNumber != "" {
		trackNumbers := strings.Split(trackNumber, ",")
		trackNumberStr := "'" + strings.Join(trackNumbers, "','") + "'"
		where += fmt.Sprintf(" and track_number in(%s)", trackNumberStr)
	}
	if status > 0 {
		where += fmt.Sprintf(" and status = %v", status)
	}
	if CarriersCode != "" {
		where += fmt.Sprintf(" and carriers_code = %q", CarriersCode)
	}
	wherestr := ""
	if where != "" {
		wherestr = "1 = 1" + where
	}
	db := common.GetDb().TrackReadDb
	err = db.Table("logistics_subscribe_record").
		Select("*").
		Where(wherestr).
		Order(fmt.Sprintf("%s  %s", sort, orderBy)).
		Find(&records).Error
	if err != nil {
		common.GetLogger().Error("db error", log.String("error msg", err.Error()))
	}
	return
}
func (s *SubscribeRecordModel) Add(subscribeRecord SubscribeRecord) (err error) {
	db := common.GetDb().TrackWriteDb
	subscribeRecord.Address = ""
	subscribeRecord.Status = 400
	subscribeRecord.ReSubscribeTime = 0
	subscribeRecord.PushStatus = 1
	subscribeRecord.PushTime = 0
	subscribeRecord.Description = ""
	subscribeRecord.UpdatedAt = common.DefaultTime()
	subscribeRecord.CreatedAt = common.DefaultTime()
	err = db.Create(&subscribeRecord).Error
	return
}

func (s *SubscribeRecordModel) Delete(where map[string]interface{}) (err error) {
	var subscribeRecord SubscribeRecord
	db := common.GetDb().TrackWriteDb
	err = db.Table(subscribeRecord.TableName()).Where(where).Delete(SubscribeRecord{}).Error
	return
}

func (s *SubscribeRecordModel) GetSubscribeRecord(params map[string]interface{}) (subscribeRecord SubscribeRecord, err error) {
	var result SubscribeRecord
	err = common.GetDb().TrackReadDb.Table(result.TableName()).
		Where(params).
		First(&subscribeRecord).Error
	return
}
