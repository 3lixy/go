package service

import (
	"azoya/nova"
	"config_service/models"
	"config_service/validate"
)

//StockService 模型
type SystemService struct {
	BaseService
}

//NewSystemService 模型
func NewSystemService(c *nova.Context) *SystemService {
	return &SystemService{BaseService{C: c}}
}

//更新站点系统初始化信息
func (s *SystemService) Update(systemField validate.UpdateSystemQuery) error {
	systemModels := models.NewSystemModels(s.C)
	system, _ := systemModels.GetWebSiteSystemStatus(map[string]interface{}{"website_id": systemField.WebsiteID, "system_id": systemField.SystemID})

	var data models.WebSysStatus
	data.WebsiteID = systemField.WebsiteID
	data.SystemID = systemField.SystemID
	data.Status = systemField.Status
	data.Message = systemField.Message

	var err error
	if system.EntityID > 0 { //存在 则更新
		err = systemModels.Update(map[string]interface{}{"website_id": systemField.WebsiteID, "system_id": systemField.SystemID}, data)
	} else { //不存在 则添加
		err = systemModels.Add(data)
	}
	return err
}

//获取站点系统初始化列表
func (s *SystemService) GetSystemList(params map[string][]string) interface{} {
	var list models.SystemList
	systemModels := models.NewSystemModels(s.C)
	err := systemModels.CheckWebsiteSystem()
	if err != nil {
		panic(err)
	}
	data, count, err := systemModels.GetSystemList(params)
	if err != nil {
		panic(err)
	}
	list.Total = count
	list.Rows = data
	return list
}

//GetDetail 获取详情
func (s *SystemService) GetDetail(entityID int) models.UnionData {
	systemModels := models.NewSystemModels(s.C)
	detail, _ := systemModels.GetSystemDetail(entityID)
	return detail
}
