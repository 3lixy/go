package models

import (
	"azoya/nova"
	"config_service/common"
	"fmt"
	"strconv"
)

//StockModels 模型
type SystemModels struct {
	BaseModels
}

//NewSystemModels 模型
func NewSystemModels(c *nova.Context) *SystemModels {
	return &SystemModels{BaseModels{C: c}}
}

type UnionData struct {
	EntityID	uint64	`gorm:"entity_id" json:"entity_id"`
	WebsiteID	uint64	`gorm:"website_id" json:"website_id"`
	SystemID	uint64	`gorm:"system_id" json:"system_id"`
	System		string	`gorm:"system" json:"system_name"`
	Status		uint64	`gorm:"status" json:"status"`
	Message		string	`gorm:"message" json:"message"`
	CreatedAt   string	`gorm:"default:'0000-00-00 00:00:00'" json:"created_at"`
	UpdatedAt   string	`gorm:"default:'0000-00-00 00:00:00'" json:"updated_at"`
}

//站点系统初始化表
func (st *UnionData) TableName() string {
	return "website_init"
}

type Website struct {
	WebsiteID    uint64 `gorm:"website_id"`
	Domain       string `gorm:"domain"`
	Name         string `gorm:"name"`
	Abbreviation string `gorm:"abbreviation"`
	CurrencyCode string `gorm:"currency_code"`
	CurrencyName string `gorm:"currency_name"`
	Status       uint64 `gorm:"status"`
	CreatedAt    string `gorm:"create_at"`
	UpdatedAt    string `gorm:"update_at"`
}

//站点表
func (st *Website) TableName() string {
	return "website"
}

type WebSysStatus struct {
	EntityID  uint64 `gorm:"entity_id" json:"entity_id"`
	WebsiteID uint64 `gorm:"website_id" json:"website_id"`
	SystemID  uint64 `gorm:"system_id" json:"system_id"`
	Status    uint64 `gorm:"status" json:"status"`
	Message   string `gorm:"message" json:"message"`
	CreatedAt string `gorm:"default:'0000-00-00 00:00:00'" json:"created_at"`
	UpdatedAt string `gorm:"default:'0000-00-00 00:00:00'" json:"updated_at"`
}

//站点系统初始化表
func (st *WebSysStatus) TableName() string {
	return "website_init"
}

type System struct {
	EntityID   uint64 `gorm:"entity_id"`
	SystemName string `gorm:"system"`
}

//系统表
func (st *System) TableName() string {
	return "system"
}

type UnionKey struct {
	WebsiteID uint64 `gorm:"website_id"`
	SystemID  uint64 `gorm:"system_id"`
}

//系统列表
type SystemList struct {
	Rows  []UnionData `json:"rows"`
	Total int64       `json:"total"`
}

func (p *SystemModels) GetWebSiteSystemStatus(params map[string]interface{}) (WebSysStatus, error) {
	var result WebSysStatus
	err := common.GetDb().ReadDb.
		Table(result.TableName()).
		Where(params).
		First(&result).
		Error
	if err != nil {
		p.C.Logger().Error(err.Error())
	}
	return result, err
}

func (p *SystemModels) Add(data WebSysStatus) (err error) {
	db := common.GetDb().WriteDb.Table(data.TableName())
	data.CreatedAt = common.DefaultTime("UTC")
	data.UpdatedAt = common.DefaultTime("UTC")
	err = db.Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *SystemModels) Update(whereParams map[string]interface{}, data WebSysStatus) error {
	data.UpdatedAt = common.DefaultTime("UTC")
	err := common.GetDb().WriteDb.Table(data.TableName()).Where(whereParams).Updates(data).Error
	return err
}

func (p *SystemModels) CheckWebsiteSystem() (err error) {
	var websites []Website
	var systems []System

	var system System
	err = common.GetDb().ReadDb.Table(system.TableName()).Find(&systems).Error
	if err != nil {
		return
	}

	var website Website
	err = common.GetDb().ReadDb.Table(website.TableName()).Find(&websites).Error
	if err != nil {
		return
	}

	//获取所有的比较
	var all []UnionKey
	for _, website := range websites {
		for _, system := range systems {
			//common.Log(fmt.Sprintf("web:%v,sys:%v", website.WebsiteID, system.EntityID), "debug", "info")
			all = append(all, UnionKey{WebsiteID: website.WebsiteID, SystemID: system.EntityID})
		}
	}

	var webSystems []WebSysStatus
	var webSysStus WebSysStatus
	err = common.GetDb().ReadDb.Table(webSysStus.TableName()).Find(&webSystems).Error
	if err != nil {
		return
	}

	datas := make(map[UnionKey]uint64)
	for _, one := range webSystems {
		datas[UnionKey{WebsiteID: one.WebsiteID, SystemID: one.SystemID}] = 0
	}

	if len(all) > len(webSystems) {
		for _, one := range all {
			if _, ok := datas[one]; !ok {
				var data WebSysStatus
				data.WebsiteID = one.WebsiteID
				data.SystemID = one.SystemID
				data.Status = 0
				data.CreatedAt = common.DefaultTime("UTC")
				data.UpdatedAt = common.DefaultTime("UTC")
				err = common.GetDb().WriteDb.Table(data.TableName()).Create(&data).Error
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func (p *SystemModels) GetSystemList(params map[string][]string) (dataList []UnionData, count int64, err error) {
	sort := "entity_id"
	orderBy := "desc"

	db := common.GetDb().ReadDb

	status := uint64(999)
	if _, ok := params["status"]; ok {
		status, _ = strconv.ParseUint(params["status"][0], 10, 64)
	}

	var websiteID uint64
	if _, ok := params["website_id"]; ok {
		websiteID, _ = strconv.ParseUint(params["website_id"][0], 10, 64)
	}

	var where string
	if status != 999 {
		where += fmt.Sprintf(" and status = %v", status)
	}

	if websiteID > 0 {
		where += fmt.Sprintf(" and website_id = %v", websiteID)
	}

	whereStr := ""
	if where != "" {
		whereStr = "1 = 1" + where
	}

	page := common.GetParamsPage(params)

	var data UnionData
	err = db.Table(fmt.Sprintf("%s AS w", data.TableName())).
		Select("w.entity_id, w.website_id, w.system_id, s.system, w.status, w.message, w.created_at, w.updated_at").
		Joins("LEFT JOIN system AS s ON s.entity_id = w.system_id ").
		Where(whereStr).
		Offset(page.Page).
		Limit(page.Limit).
		Order(fmt.Sprintf("%s  %s", sort, orderBy)).
		Find(&dataList).Error

	db.Table(data.TableName()).
		Where(whereStr).
		Count(&count)
	return
}

//GetSystemDetail 获取system init的详情
func (p *SystemModels) GetSystemDetail(entityID int) (UnionData, error) {
	var detail UnionData

	db := common.GetDb().ReadDb
	err := db.Table(fmt.Sprintf("%s AS w", detail.TableName())).
		Select("w.entity_id, w.website_id, w.system_id, s.system, w.status, w.message, w.created_at, w.updated_at").
		Joins("LEFT JOIN system AS s ON s.entity_id = w.system_id ").
		Where("w.entity_id = ?", entityID).
		Find(&detail).Error

	return detail, err
}
