package models

import (
	"azoya/nova"
	"fmt"
	"github.com/jinzhu/gorm"
	"tms_service/common"
	"github.com/caibirdme/yql"
	"math/rand"
)

//LineModel 用于读取line相关数据
type LineModel struct {
	Context *nova.Context
}

//NewLineModel 初始化LineModel
func NewLineModel(c *nova.Context) *LineModel {
	return &LineModel{Context: c}
}

//GetLineWithLogisticsID 根据物流商id来获取物流线路
func (l *LineModel) GetLineWithLogisticsID(logisticsID uint64) ([]Line, int, error) {
	var list []Line
	err := common.GetDb().TmsReadDb.Table(l.TableName()).
		Where("main_logistics_id = ? OR oversea_logistics_id = ? OR inland_logistics_id = ?", logisticsID, logisticsID, logisticsID).
		Find(&list).
		Error

	return list, len(list), err
}

//GetList 根据传入的参数，获取line列表
func (l *LineModel) GetList(params LineListQueryParams, page common.Page) ([]Line, uint64, error) {
	count, err := l.GetTotalRowsCount(params)
	if err != nil {
		panic(err)
	}
	page.TotalCount = count
	page = common.GetValidPage(page)

	if page.Sort == "" {
		page.Sort = l.PrimaryKey()
	}

	if page.Order == "" {
		page.Order = l.defaultOrder()
	}
	var list []Line
	err = common.GetDb().TmsReadDb.Table(l.TableName()).
		Where(params).
		Offset(page.Page).
		Limit(page.Limit).
		Order(fmt.Sprintf("%s  %s", page.Sort, page.Order)).
		Find(&list).
		Error

	return list, page.TotalCount, err
}

//GetAll 获取所有line列表
func (l *LineModel) GetAll() ([]Line, uint64, error) {
	var params LineListQueryParams
	count, err := l.GetTotalRowsCount(params)
	if err != nil {
		panic(err)
	}

	var list []Line
	err = common.GetDb().TmsReadDb.Table(l.TableName()).
		Find(&list).
		Error

	return list, count, err
}

//GetLineListWithoutPage 不分页直接获取line列表
func (l *LineModel) GetLineListWithoutPage(params LineListQueryParams) ([]Line, uint64, error) {
	count, err := l.GetTotalRowsCount(params)
	if err != nil {
		panic(err)
	}

	var list []Line
	err = common.GetDb().TmsReadDb.Table(l.TableName()).
		Where(params).
		Find(&list).
		Error

	return list, count, err
}

//GetTotalRowsCount 获取line的总数
func (l *LineModel) GetTotalRowsCount(params LineListQueryParams) (uint64, error) {
	var count uint64
	err := common.GetDb().TmsReadDb.Table(l.TableName()).
		Where(params).
		Count(&count).
		Error

	return count, err
}

//GetDetail 获取line详情
func (l *LineModel) GetDetail(lineID uint64) (Line, error) {
	var Line Line
	err := common.GetDb().TmsReadDb.Table(l.TableName()).
		Where(l.getWhereString(), lineID).
		Find(&Line).
		Error

	if err != nil && err == gorm.ErrRecordNotFound {
		l.Context.Logger().Error(err.Error())
	}

	return Line, err
}

//Delete 删除line
func (l *LineModel) Delete(lineID uint64) error {
	var Line Line

	Line, err := l.GetDetail(lineID)
	if err != nil {
		return err
	}

	deleteErr := common.GetDb().TmsWriteDb.Table(l.TableName()).
		Where(l.getWhereString(), lineID).
		Delete(&Line).
		Error

	if deleteErr != nil {
		l.Context.Logger().Error(err.Error())
	}

	return deleteErr
}

//Update update line数据
func (l *LineModel) Update(Line Line) (Line, error) {
	Line.UpdatedAt = common.DefaultTime()
	err := common.GetDb().TmsWriteDb.Table(l.TableName()).
		Model(&Line).
		Where(l.getWhereString(), Line.EntityID).
		Update(map[string]interface{}{
		"website_id":Line.WebsiteID,
		"store_id":Line.StoreID,
		"title":Line.Title,
		"main_logistics_id":Line.MainLogisticsID,
		"oversea_logistics_id":Line.OverseaLogisticsID,
		"inland_logistics_id":Line.InlandLogisticsID,
		"customs_declaration_type":Line.CustomsDeclarationType,
		"routes_type":Line.RoutesType,
		"start_country":Line.StartCountry,
		"target_country":Line.TargetCountry,
		"status":Line.Status,
		"delivery_type":Line.DeliveryType,
		"is_push_track_number":Line.IsPushTrackNumber,
		"label_push_type":Line.LabelPushType,
		"label_template":Line.LabelTemplate,
		"is_sync_logistics":Line.IsSyncLogistics,
		"rule_detail":Line.RuleDetail,
		"rule_desc":Line.RuleDesc,
		"probability":Line.Probability,
		"updated_at":Line.UpdatedAt,
	}).Error

	if err != nil {
		l.Context.Logger().Error(err.Error())
	}

	return Line, err
}

//Create create line数据
func (l *LineModel) Create(Line Line) (Line, error) {
	time := common.DefaultTime()
	Line.CreatedAt = time
	Line.UpdatedAt = time
	err := common.GetDb().TmsWriteDb.Table(l.TableName()).Create(&Line).Error

	if err != nil {
		l.Context.Logger().Error(err.Error())
	}

	return Line, err
}

//TableName 数据库表名称
func (l *LineModel) TableName() string {
	return "transport_line"
}

//PrimaryKey 返回Line的主健
func (l *LineModel) PrimaryKey() string {
	return "entity_id"
}

//defaultOrder 默认排序规则
func (l *LineModel) defaultOrder() string {
	return "desc"
}

func (l *LineModel) getWhereString() string {
	return l.PrimaryKey() + "=?"
}

func (l *LineModel) IsMatchLine(line Line, ruleData RuleData) (bool, error) {
	rawYQL := line.RuleDesc
	result, err := yql.Match(rawYQL, map[string]interface{}{
		"sku":         ruleData.Sku,
		"weight":      ruleData.Weight,
		"grand_total": ruleData.GrandTotal,
		"country_id":  ruleData.CountryID,
		"region":      ruleData.Region,
	})
	return result, err
}

func (l *LineModel) GetAwardLine(list []Line) (line Line) {
	type awardLine struct {
		line   Line
		offset int64
		count  int64
	}
	lineSli := make([]*awardLine, 0, len(list))
	var sumCount int64 = 0
	for _, line := range list {
		a := awardLine{
			line:   line,
			offset: sumCount,
			count:  line.Probability,
		}
		//整理所有用户的count数据为数轴
		lineSli = append(lineSli, &a)
		sumCount += line.Probability
	}

	awardIndex := rand.Int63n(sumCount)
	for _, l := range lineSli {
		//判断命中index落在那个路线区间内
		if l.offset+l.count > awardIndex {
			line = l.line
			return
		}
	}
	return
}

//Line 物流商
type Line struct {
	EntityID               uint64 `gorm:"primary_key" json:"entity_id" form:"entity_id"`
	WebsiteID              int    `gorm:"website_id" json:"website_id" form:"website_id" binding:"required"`
	StoreID                int    `gorm:"store_id" json:"store_id" form:"store_id" binding:"required"`
	Title                  string `gorm:"title" json:"title" form:"title" binding:"required"`
	MainLogisticsID        uint64 `gorm:"main_logistics_id" json:"main_logistics_id" form:"main_logistics_id" binding:"required"`
	OverseaLogisticsID     uint64 `gorm:"oversea_logistics_id" json:"oversea_logistics_id" form:"oversea_logistics_id"`
	InlandLogisticsID      uint64 `gorm:"inland_logistics_id" json:"inland_logistics_id" form:"inland_logistics_id"`
	CustomsDeclarationType uint64 `gorm:"customs_declaration_type" json:"customs_declaration_type" form:"customs_declaration_type"  binding:"required"`
	RoutesType             string `gorm:"routes_type" json:"routes_type" form:"routes_type"  binding:"required"`
	StartCountry           string `gorm:"start_country" json:"start_country" form:"start_country" binding:"required"`
	TargetCountry          string `gorm:"target_country" json:"target_country" form:"target_country"`
	Status                 int    `gorm:"status" json:"status" form:"status" binding:"required"`
	DeliveryType           int    `gorm:"delivery_type" json:"delivery_type" form:"delivery_type" binding:"required"`
	IsPushTrackNumber      int    `gorm:"is_push_track_number" json:"is_push_track_number" form:"is_push_track_number" binding:"required"`
	LabelPushType          int    `gorm:"label_push_type" json:"label_push_type" form:"label_push_type" binding:"required"`
	LabelTemplate          int    `gorm:"label_template" json:"label_template" form:"label_template"`
	IsSyncLogistics        int    `gorm:"is_sync_logistics" json:"is_sync_logistics" form:"is_sync_logistics" binding:"required"`
	RuleDetail             string `gorm:"rule_detail" json:"rule_detail" form:"rule_detail"`
	RuleDesc               string `gorm:"rule_desc" json:"rule_desc" form:"rule_desc"`
	Probability            int64  `gorm:"probability" json:"probability" form:"probability"`
	CreatedAt              string `gorm:"created_at" json:"created_at" form:"created_at"`
	UpdatedAt              string `gorm:"updated_at" json:"updated_at" form:"updated_at"`
}

//LineListQueryParams 用于查询list
type LineListQueryParams struct {
	WebsiteID              string `gorm:"website_id" json:"website_id" form:"website_id" binding:"required"`
	StoreID                string `gorm:"store_id" json:"store_id" form:"store_id" `
	Status                 string `gorm:"status" json:"status" form:"status"`
	MainLogisticsID        string `gorm:"main_logistics_id" json:"main_logistics_id" form:"main_logistics_id"`
	InlandLogisticsID      uint64 `gorm:"inland_logistics_id" json:"inland_logistics_id" form:"inland_logistics_id"`
	OverseaLogisticsID     string `gorm:"oversea_logistics_id" json:"oversea_logistics_id" form:"oversea_logistics_id"`
	RoutesType             string `gorm:"routes_type" json:"routes_type" form:"routes_type"`
	CustomsDeclarationType string `gorm:"customs_declaration_type" json:"customs_declaration_type" form:"customs_declaration_type"`
	DeliveryType           string `gorm:"delivery_type" json:"delivery_type" form:"delivery_type"`
}

type RuleData struct {
	Sku        []string `json:"sku"`
	Weight     float64  `json:"weight"`
	GrandTotal float64  `json:"grand_total"`
	CountryID  string   `json:"country_id"`
	Region     string   `json:"region"`
}
