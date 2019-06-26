package models

import (
	"azoya/nova"
	"config_service/common"
	"fmt"
	"strconv"
)

//StockModels 模型
type StockModels struct {
	BaseModels
}

//NewStockModels 模型
func NewStockModels(c *nova.Context) *StockModels {
	return &StockModels{BaseModels{C: c}}
}

//Stock 仓库模型
type Stock struct {
	EntityID    uint64 `gorm:"primary_key" json:"entity_id"`
	StockType   uint64 `gorm:"stock_type" json:"stock_type"`
	StockName   string `gorm:"stock_name" json:"stock_name"`
	CompanyName string `gorm:"company_name" json:"company_name"`
	Status      uint64 `gorm:"status" json:"status"`
	Country     string `gorm:"country" json:"country"`
	Province    string `gorm:"province" json:"province"`
	City        string `gorm:"city" json:"city"`
	County      string `gorm:"county" json:"county"`
	AddressOne  string `gorm:"address_one" json:"address_one"`
	AddressTwo  string `gorm:"address_two" json:"address_two"`
	Postcode    string `gorm:"postcode" json:"postcode"`
	LastName    string `gorm:"last_name" json:"last_name"`
	FirstName   string `gorm:"first_name" json:"first_name"`
	Position    string `gorm:"position" json:"position"`
	Telephone   string `gorm:"telephone" json:"telephone"`
	Email       string `gorm:"email" json:"email"`
	Wechat      string `gorm:"wechat" json:"wechat"`
	CreatedAt   string `gorm:"default:'0000-00-00 00:00:00'" json:"created_at"`
	UpdatedAt   string `gorm:"default:'0000-00-00 00:00:00'" json:"updated_at"`
}

//仓库列表
type StockList struct {
	Rows  []Stock `json:"rows"`
	Total int64   `json:"total"`
}

//TableName 仓库表
func (st *Stock) TableName() string {
	return "stock"
}

//GetStock 获取仓库详情
func (s *StockModels) GetStock(params map[string]interface{}) (Stock, error) {
	var result Stock
	err := common.GetDb().ReadDb.Table(result.TableName()).
		Where(params).
		First(&result).
		Error
	if err != nil {
		s.C.Logger().Error(err.Error())
	}
	return result, err
}

//GetList 获取合同列表
func (s *StockModels) GetList(params map[string][]string) (StockList []Stock, count int64, err error) {
	sort := "entity_id"
	orderBy := "asc"

	if _, ok := params["sort"]; ok {
		sort = params["sort"][0]
	}
	if _, ok := params["order"]; ok {
		orderBy = params["order"][0]
	}

	var status uint64

	if _, ok := params["status"]; ok {
		status, _ = strconv.ParseUint(params["status"][0], 10, 64)
	}

	var where string

	if status > 0 {
		where += fmt.Sprintf(" and status = %v", status)
	}

	wherestr := ""
	if where != "" {
		wherestr = "1 = 1" + where
	}
	//page := common.GetParamsPage(params)
	db := common.GetDb().ReadDb

	var stock Stock
	err = db.Model(stock).
		Where(wherestr).
		Order(fmt.Sprintf("%s  %s", sort, orderBy)).
		Find(&StockList).Error

	db.Model(stock).
		Where(wherestr).
		Count(&count)
	return
}

//Add 添加仓库
func (s *StockModels) Add(stock Stock) (err error) {
	db := common.GetDb().WriteDb
	stock.CreatedAt = common.DefaultTime("UTC")
	stock.UpdatedAt = common.DefaultTime("UTC")
	err = db.Create(&stock).Error
	if err != nil {
		return err
	}
	return nil
}

//Update 仓库修改
func (s *StockModels) Update(whereParams map[string]interface{}, stock Stock) error {
	var stockTable Stock
	stock.UpdatedAt = common.DefaultTime("UTC")
	err := common.GetDb().WriteDb.Model(&stockTable).Where(whereParams).Updates(stock).Error
	return err
}

//CheckStock 获取仓库详情
func (s *StockModels) CheckStock(name string, stockID uint64) (Stock, error) {
	var result Stock
	err := common.GetDb().ReadDb.Table(result.TableName()).
		Where("stock_name = ? and entity_id != ?", name, stockID).
		First(&result).
		Error
	if err != nil {
		s.C.Logger().Error(err.Error())
	}
	return result, err
}
