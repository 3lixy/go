package models

import (
	"azoya/nova"
	"config_service/common"
	"fmt"
	"strconv"
)

//StoreModels 模型
type StoreModels struct {
	BaseModels
}

//NewStockModels 模型
func NewStoreModels(c *nova.Context) *StoreModels {
	return &StoreModels{BaseModels{C: c}}
}

//Stock 仓库模型
type Store struct {
	EntityID           uint64  `gorm:"primary_key" json:"entity_id"`
	StoreID            uint64  `gorm:"store_id" json:"store_id"`
	MerchantID         uint64  `gorm:"merchant_id" json:"merchant_id"`
	WebsiteID          uint64  `gorm:"website_id" json:"website_id"`
	Name               string  `gorm:"name" json:"name"`
	HomeUrl            string  `gorm:"home_url" json:"home_url"`
	H5Url              string  `gorm:"h5_url" json:"h5_url"`
	ShippingOrigin     string  `gorm:"shipping_origin" json:"shipping_origin"`
	TopTips            string  `gorm:"top_tips" json:"top_tips"`
	BottomTips         string  `gorm:"bottom_tips" json:"bottom_tips"`
	TaxType            uint64  `gorm:"tax_type" json:"tax_type"`
	CountryCode        string  `gorm:"country_code" json:"country_code"`
	IsShowFreightDesc  uint64  `gorm:"is_show_freight_desc" json:"is_show_freight_desc"`
	FreightDesc        string  `gorm:"freight_desc" json:"freight_desc"`
	IsShowTaxDesc      uint64  `gorm:"is_show_tax_desc" json:"is_show_tax_desc"`
	TaxDesc            string  `gorm:"tax_desc" json:"tax_desc"`
	IsShowDeliveryDesc uint64  `gorm:"is_show_delivery_desc" json:"is_show_delivery_desc"`
	DeliveryDesc       string  `gorm:"delivery_desc" json:"delivery_desc"`
	AuthenticationType uint64  `gorm:"authentication_type" json:"authentication_type"`
	ExchangeRate       float64 `gorm:"exchange_rate" json:"exchange_rate"`
	Operater           string  `gorm:"operater" json:"operater"`
	Status             uint64  `gorm:"status" json:"status"`
	IsBonded           uint64  `gorm:"is_bonded" json:"is_bonded"`
	StockID            uint64  `gorm:"stock_id" json:"stock_id"`
	CreatedAt          string  `gorm:"default:'0000-00-00 00:00:00'" json:"created_at"`
	UpdatedAt          string  `gorm:"default:'0000-00-00 00:00:00'" json:"updated_at"`
}

type StoreDetail struct {
	EntityID           uint64  `json:"entity_id"`
	StoreID            uint64  `json:"store_id"`
	MerchantID         uint64  `json:"merchant_id"`
	WebsiteID          uint64  `json:"website_id"`
	Name               string  `json:"name"`
	HomeUrl            string  `json:"home_url"`
	H5Url              string  `json:"h5_url"`
	ShippingOrigin     string  `json:"shipping_origin"`
	TopTips            string  `json:"top_tips"`
	BottomTips         string  `json:"bottom_tips"`
	TaxType            uint64  `json:"tax_type"`
	CountryCode        string  `json:"country_code"`
	IsShowFreightDesc  uint64  `json:"is_show_freight_desc"`
	FreightDesc        string  `json:"freight_desc"`
	IsShowTaxDesc      uint64  `json:"is_show_tax_desc"`
	TaxDesc            string  `json:"tax_desc"`
	IsShowDeliveryDesc uint64  `json:"is_show_delivery_desc"`
	DeliveryDesc       string  `json:"delivery_desc"`
	AuthenticationType uint64  `json:"authentication_type"`
	ExchangeRate       float64 `json:"exchange_rate"`
	Operater           string  `json:"operater"`
	Status             uint64  `json:"status"`
	IsBonded           uint64  `json:"is_bonded"`
	StockID            uint64  `json:"stock_id"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
	Stock              *Stock  `json:"stock"`
}

//TableName 店铺表
func (st *Store) TableName() string {
	return "store"
}

//GetStore 获取店铺详情
func (s *StoreModels) GetStore(params map[string]interface{}) (Store, error) {
	var result Store
	err := common.GetDb().ReadDb.Table(result.TableName()).
		Where(params).
		First(&result).
		Error
	if err != nil {
		s.C.Logger().Error(err.Error())
	}
	return result, err
}

func (s *StoreModels) GetAllStore(params map[string][]string) (StoreList []Store, err error) {
	sort := "entity_id"
	orderBy := "asc"

	if _, ok := params["sort"]; ok {
		sort = params["sort"][0]
	}
	if _, ok := params["order"]; ok {
		orderBy = params["order"][0]
	}

	var status uint64
	var websiteID uint64

	if _, ok := params["website_id"]; ok {
		websiteID, _ = strconv.ParseUint(params["website_id"][0], 10, 64)
	}

	if _, ok := params["status"]; ok {
		status, _ = strconv.ParseUint(params["status"][0], 10, 64)
	}

	var where string

	if status > 0 {
		where += fmt.Sprintf(" and status = %v", status)
	}

	if websiteID > 0 {
		where += fmt.Sprintf(" and website_id = %v", websiteID)
	}

	wherestr := ""
	if where != "" {
		wherestr = "1 = 1" + where
	}
	//page := common.GetParamsPage(params)
	db := common.GetDb().ReadDb

	var store Store
	err = db.Model(store).
		Where(wherestr).
		Order(fmt.Sprintf("%s  %s", sort, orderBy)).
		Find(&StoreList).Error
	return
}
