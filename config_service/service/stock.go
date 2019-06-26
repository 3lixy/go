package service

import (
	"azoya/nova"
	"config_service/common"
	"config_service/models"
	"config_service/validate"
)

//StockService 模型
type StockService struct {
	BaseService
}

//NewStockService 模型
func NewStockService(c *nova.Context) *StockService {
	return &StockService{BaseService{C: c}}
}

//Add 仓库信息增加
func (s *StockService) Add(stockField validate.AddQuery) error {

	stockModels := models.NewStockModels(s.C)
	stock, _ := stockModels.GetStock(map[string]interface{}{"stock_name": stockField.StockName})
	var err error
	if stock.EntityID > 0 {
		return common.ErrStockNameIsExist
	}
	var stockData models.Stock
	stockData.StockType = stockField.StockType
	stockData.StockName = stockField.StockName
	stockData.CompanyName = stockField.CompanyName
	stockData.Status = stockField.Status
	stockData.Country = stockField.Country
	stockData.Province = stockField.Province
	stockData.City = stockField.City
	stockData.County = stockField.County
	stockData.AddressOne = stockField.AddressOne
	stockData.AddressTwo = stockField.AddressTwo
	stockData.Postcode = stockField.Postcode
	stockData.FirstName = stockField.FirstName
	stockData.LastName = stockField.LastName
	stockData.Position = stockField.Position
	stockData.Telephone = stockField.Telephone
	stockData.Email = stockField.Email
	stockData.Wechat = stockField.Wechat
	err = stockModels.Add(stockData)
	return err
}

//Update 更新仓库信息
func (s *StockService) Update(stockField validate.UpdateQuery) error {
	stockModels := models.NewStockModels(s.C)
	stock, _ := stockModels.GetStock(map[string]interface{}{"entity_id": stockField.EntityID})
	//仓库不存在 无法更新
	var err error
	if stock.EntityID <= 0 {
		return common.ErrStockNotExist
	}
	check, _ := stockModels.CheckStock(stockField.StockName, stock.EntityID)
	if check.EntityID > 0 {
		return common.ErrStockNameIsExist
	}
	var stockData models.Stock
	stockData.StockType = stockField.StockType
	stockData.StockName = stockField.StockName
	stockData.CompanyName = stockField.CompanyName
	stockData.Status = stockField.Status
	stockData.Country = stockField.Country
	stockData.Province = stockField.Province
	stockData.City = stockField.City
	stockData.County = stockField.County
	stockData.AddressOne = stockField.AddressOne
	stockData.AddressTwo = stockField.AddressTwo
	stockData.Postcode = stockField.Postcode
	stockData.FirstName = stockField.FirstName
	stockData.LastName = stockField.LastName
	stockData.Position = stockField.Position
	stockData.Telephone = stockField.Telephone
	stockData.Email = stockField.Email
	stockData.Wechat = stockField.Wechat

	err = stockModels.Update(map[string]interface{}{"entity_id": stock.EntityID}, stockData)
	return err
}

//GetStockList 获取仓库列表
func (s *StockService) GetStockList(params map[string][]string) interface{} {
	var list models.StockList
	stockModels := models.NewStockModels(s.C)
	data, count, err := stockModels.GetList(params)
	if err != nil {
		panic(err)
	}
	list.Total = count
	list.Rows = data
	return list
}

//GetDetail 获取仓库详情接口
func (s *StockService) GetDetail(query validate.DetailQuery) (stock models.Stock, err error) {
	stockModels := models.NewStockModels(s.C)
	stock, err = stockModels.GetStock(map[string]interface{}{"entity_id": query.EntityID})
	return
}
