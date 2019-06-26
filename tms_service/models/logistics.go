package models

import (
	"azoya/nova"
	"fmt"
	"github.com/jinzhu/gorm"
	"tms_service/common"
)

//LogisticsModel 用于读取logistics相关数据
type LogisticsModel struct {
	Context *nova.Context
}

//NewLogisticsModel 初始化LogisticsModel
func NewLogisticsModel(c *nova.Context) *LogisticsModel {
	return &LogisticsModel{Context: c}
}

//GetList 获取logistics列表
func (l *LogisticsModel) GetList(params ListQueryParams) ([]Logistics, uint64, error) {
	count, err := l.GetTotalRowsCount(params)
	if err != nil {
		panic(err)
	}

	pageModel := common.Page{Page: params.Page, Limit: params.Limit, TotalCount: count}
	page := common.GetValidPage(pageModel)

	if params.Sort == "" {
		params.Sort = l.PrimaryKey()
	}

	if params.Order == "" {
		params.Order = l.defaultOrder()
	}
	var list []Logistics
	err = common.GetDb().TmsReadDb.Table(l.TableName()).
		Where("company_name LIKE ?", "%"+params.CompanyName+"%").
		Where("name LIKE ?", "%"+params.Name+"%").
		Where(&ListQueryParams{Status: params.Status}).
		Offset(page.Page).
		Limit(page.Limit).
		Order(fmt.Sprintf("%s  %s", params.Sort, params.Order)).
		Find(&list).
		Error

	return list, count, err
}

//GetTotalRowsCount 获取logistics的总数
func (l *LogisticsModel) GetTotalRowsCount(params ListQueryParams) (uint64, error) {
	var count uint64
	err := common.GetDb().TmsReadDb.Table(l.TableName()).
		Where("company_name LIKE ?", "%"+params.CompanyName+"%").
		Where("name LIKE ?", "%"+params.Name+"%").
		Where(&ListQueryParams{Status: params.Status}).
		Count(&count).
		Error

	return count, err
}

//GetDetail 获取logistics详情
func (l *LogisticsModel) GetDetail(logisticsID uint64) (Logistics, error) {
	var logistics Logistics
	err := common.GetDb().TmsReadDb.Table(l.TableName()).
		Where(l.getWhereString(), logisticsID).
		Find(&logistics).
		Error

	if err != nil && err == gorm.ErrRecordNotFound {
		l.Context.Logger().Error(err.Error())
	}

	return logistics, err
}

//Delete 删除logistics
func (l *LogisticsModel) Delete(logisticsID uint64) error {
	var logistics Logistics

	logistics, err := l.GetDetail(logisticsID)
	if err != nil {
		return err
	}

	deleteErr := common.GetDb().TmsWriteDb.Table(l.TableName()).
		Where(l.getWhereString(), logisticsID).
		Delete(&logistics).
		Error

	if deleteErr != nil {
		l.Context.Logger().Error(err.Error())
	}

	return deleteErr
}

//Update update logistics数据
func (l *LogisticsModel) Update(logistics Logistics) (Logistics, error) {
	logistics.UpdatedAt = common.DefaultTime()
	err := common.GetDb().TmsWriteDb.Table(l.TableName()).
		Model(&logistics).
		Where(l.getWhereString(), logistics.EntityID).
		Update(logistics).
		Error

	if err != nil {
		l.Context.Logger().Error(err.Error())
	}

	return logistics, err
}

//Create create logistics数据
func (l *LogisticsModel) Create(logistics Logistics) (Logistics, error) {
	time := common.DefaultTime()
	logistics.CreatedAt = time
	logistics.UpdatedAt = time
	err := common.GetDb().TmsWriteDb.Table(l.TableName()).Create(&logistics).Error

	if err != nil {
		l.Context.Logger().Error(err.Error())
	}

	return logistics, err
}

//TableName 数据库表名称
func (l *LogisticsModel) TableName() string {
	return "logistics"
}

//PrimaryKey 返回Logistics的主健
func (l *LogisticsModel) PrimaryKey() string {
	return "entity_id"
}

//defaultOrder 默认排序规则
func (l *LogisticsModel) defaultOrder() string {
	return "desc"
}

func (l *LogisticsModel) getWhereString() string {
	return l.PrimaryKey() + "=?"
}

//Logistics 物流商
type Logistics struct {
	EntityID        uint64 `gorm:"primary_key" json:"entity_id" form:"entity_id"`
	Code            string `gorm:"code" json:"code" form:"code" binding:"required"`
	Type            string `gorm:"type" json:"type" form:"type" binding:"required"`
	Name            string `gorm:"name" json:"name" form:"name" binding:"required"`
	CompanyName     string `gorm:"company_name" json:"company_name" form:"company_name" binding:"required"`
	Country         string `gorm:"country" json:"country" form:"country" binding:"required"`
	Province        string `gorm:"province" json:"province" form:"province"`
	City            string `gorm:"city" json:"city" form:"city"`
	County          string `gorm:"county" json:"county" form:"county"`
	Address         string `gorm:"address" json:"address" form:"address" binding:"required"`
	ZipCode         string `gorm:"zip_code" json:"zip_code" form:"zip_code"`
	APIType         string `gorm:"api_type" json:"api_type" form:"api_type" binding:"required"`
	SubscribeTypeID int    `gorm:"subscribe_type_id" json:"subscribe_type_id" form:"subscribe_type_id" binding:"required"`
	Status          int    `gorm:"status" json:"status" form:"status" binding:"required"`
	QueryURL        string `gorm:"query_url" json:"query_url" form:"query_url" binding:"required"`
	Firstname       string `gorm:"firstname" json:"firstname" form:"firstname" binding:"required"`
	Lastname        string `gorm:"lastname" json:"lastname" form:"lastname" binding:"required"`
	Cellphone       string `gorm:"cellphone" json:"cellphone" form:"cellphone"`
	Email           string `gorm:"email" json:"email" form:"email" binding:"required,email"`
	Position        string `gorm:"position" json:"position" form:"position"`
	WechatID        string `gorm:"wechat_id" json:"wechat_id" form:"wechat_id"`
	Note            string `gorm:"note" json:"note" form:"note"`
	CarrierCode     string `gorm:"carrier_code" json:"carrier_code" form:"carrier_code"`
	CreatedAt       string `gorm:"created_at" json:"created_at" form:"created_at"`
	UpdatedAt       string `gorm:"updated_at" json:"updated_at" form:"updated_at"`
}

//ListQueryParams 用于查询list
type ListQueryParams struct {
	Name        string `gorm:"name" json:"name" form:"name"`
	CompanyName string `gorm:"company_name" json:"company_name" form:"company_name" `
	Status      string `gorm:"status" json:"status" form:"status"`
	Page        uint64 `gorm:"page" json:"page" form:"page"`
	Sort        string `gorm:"sort" json:"sort" form:"sort"`
	Order       string `gorm:"order" json:"order" form:"order"`
	Limit       uint64 `gorm:"limit" json:"limit" form:"limit"`
}

func (l *LogisticsModel) GetLogistics(params map[string]interface{}) (logistics Logistics, err error) {
	err = common.GetDb().TmsReadDb.Table(l.TableName()).
		Where(params).
		First(&logistics).Error
	return
}
