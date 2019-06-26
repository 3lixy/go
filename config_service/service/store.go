package service

import (
	"azoya/nova"
	"config_service/common"
	"config_service/models"
	"config_service/validate"
)

//StoreService 模型
type StoreService struct {
	BaseService
}

//NewStoreService 模型
func NewStoreService(c *nova.Context) *StoreService {
	return &StoreService{BaseService{C: c}}
}

//Detail 获取店铺详情接口
func (s *StoreService) Detail(query validate.StoreDetailQuery) (storeDetail models.StoreDetail, err error) {
	stockModels := models.NewStockModels(s.C)
	storeModels := models.NewStoreModels(s.C)

	store, _ := storeModels.GetStore(map[string]interface{}{"website_id": query.WebsiteID, "store_id": query.StoreID})
	if store.EntityID <= 0 {
		err = common.ErrStoreNotExist
		return
	}
	stock, _ := stockModels.GetStock(map[string]interface{}{"entity_id": store.StockID})
	if stock.EntityID <= 0 {
		storeDetail.Stock = nil
	} else {

		storeDetail.Stock = &stock
	}
	storeDetail.EntityID = store.EntityID
	storeDetail.StoreID = store.StoreID
	storeDetail.MerchantID = store.MerchantID
	storeDetail.Status = store.Status
	storeDetail.WebsiteID = store.WebsiteID
	storeDetail.Name = store.Name
	storeDetail.HomeUrl = store.HomeUrl
	storeDetail.H5Url = store.H5Url
	storeDetail.ShippingOrigin = store.ShippingOrigin
	storeDetail.TopTips = store.TopTips
	storeDetail.BottomTips = store.BottomTips
	storeDetail.TaxType = store.TaxType
	storeDetail.CountryCode = store.CountryCode
	storeDetail.IsShowFreightDesc = store.IsShowFreightDesc
	storeDetail.FreightDesc = store.FreightDesc
	storeDetail.IsShowTaxDesc = store.IsShowTaxDesc
	storeDetail.TaxDesc = store.TaxDesc
	storeDetail.IsShowDeliveryDesc = store.IsShowDeliveryDesc
	storeDetail.DeliveryDesc = store.DeliveryDesc
	storeDetail.AuthenticationType = store.AuthenticationType
	storeDetail.ExchangeRate = store.ExchangeRate
	storeDetail.Operater = store.Operater
	storeDetail.IsBonded = store.IsBonded
	storeDetail.CreatedAt = store.CreatedAt
	storeDetail.UpdatedAt = store.UpdatedAt
	storeDetail.StockID = store.StockID
	return
}

//GetStockList 获取店铺列表
func (s *StoreService) GetStoreList(params map[string][]string) interface{} {
	storeModels := models.NewStoreModels(s.C)
	stockModels := models.NewStockModels(s.C)
	var result []models.StoreDetail
	stores, err := storeModels.GetAllStore(params)
	if err != nil {
		panic(err)
	}
	stockParams := make(map[string][]string)
	stocks, _, err := stockModels.GetList(stockParams)
	if err != nil {
		panic(err)
	}
	stockMap := make(map[uint64]models.Stock)
	for _, stock := range stocks {
		stockMap[stock.EntityID] = stock
	}
	for _, store := range stores {
		var storeDetail models.StoreDetail
		if value, ok := stockMap[store.StockID]; ok {
			storeDetail.Stock = &value
		} else {
			storeDetail.Stock = nil
		}
		storeDetail.EntityID = store.EntityID
		storeDetail.StoreID = store.StoreID
		storeDetail.MerchantID = store.MerchantID
		storeDetail.Status = store.Status
		storeDetail.WebsiteID = store.WebsiteID
		storeDetail.Name = store.Name
		storeDetail.HomeUrl = store.HomeUrl
		storeDetail.H5Url = store.H5Url
		storeDetail.ShippingOrigin = store.ShippingOrigin
		storeDetail.TopTips = store.TopTips
		storeDetail.BottomTips = store.BottomTips
		storeDetail.TaxType = store.TaxType
		storeDetail.CountryCode = store.CountryCode
		storeDetail.IsShowFreightDesc = store.IsShowFreightDesc
		storeDetail.FreightDesc = store.FreightDesc
		storeDetail.IsShowTaxDesc = store.IsShowTaxDesc
		storeDetail.TaxDesc = store.TaxDesc
		storeDetail.IsShowDeliveryDesc = store.IsShowDeliveryDesc
		storeDetail.DeliveryDesc = store.DeliveryDesc
		storeDetail.AuthenticationType = store.AuthenticationType
		storeDetail.ExchangeRate = store.ExchangeRate
		storeDetail.Operater = store.Operater
		storeDetail.IsBonded = store.IsBonded
		storeDetail.CreatedAt = store.CreatedAt
		storeDetail.UpdatedAt = store.UpdatedAt
		storeDetail.StockID = store.StockID
		result = append(result, storeDetail)
	}
	return result
}
