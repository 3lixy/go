package models

import (
	"azoya/nova"
	"github.com/jinzhu/gorm"
	"tms_service/common"
)

//CustomsDeclarationTypeModel 用于读取CustomsDeclarationType相关数据
type CustomsDeclarationTypeModel struct {
	Context *nova.Context
}

//NewCustomsDeclarationTypeModel 初始化CustomsDeclarationTypeModel
func NewCustomsDeclarationTypeModel(c *nova.Context) *CustomsDeclarationTypeModel {
	return &CustomsDeclarationTypeModel{Context: c}
}

//GetList 获取CustomsDeclarationType列表
func (l *CustomsDeclarationTypeModel) GetList() ([]CustomsDeclarationType, error) {
	var list []CustomsDeclarationType
	err := common.GetDb().TmsReadDb.Table(l.TableName()).
		Find(&list).
		Error

	return list, err
}

//GetDetail 获取CustomsDeclarationType详情
func (l *CustomsDeclarationTypeModel) GetDetail(customsDeclarationTypeID uint64) (CustomsDeclarationType, error) {
	var customsDeclarationType CustomsDeclarationType
	err := common.GetDb().TmsReadDb.Table(l.TableName()).
		Where(l.getWhereString(), customsDeclarationTypeID).
		Find(&customsDeclarationType).
		Error

	if err != nil && err == gorm.ErrRecordNotFound {
		l.Context.Logger().Error(err.Error())
	}

	return customsDeclarationType, err
}

//Delete 删除CustomsDeclarationType
func (l *CustomsDeclarationTypeModel) Delete(customsDeclarationTypeID uint64) error {
	var customsDeclarationType CustomsDeclarationType

	customsDeclarationType, err := l.GetDetail(customsDeclarationTypeID)
	if err != nil {
		return err
	}

	deleteErr := common.GetDb().TmsWriteDb.Table(l.TableName()).
		Where(l.getWhereString(), customsDeclarationTypeID).
		Delete(&customsDeclarationType).
		Error

	if deleteErr != nil {
		l.Context.Logger().Error(err.Error())
	}

	return deleteErr
}

//Create create CustomsDeclarationType数据
func (l *CustomsDeclarationTypeModel) Create(customsDeclarationType CustomsDeclarationType) (CustomsDeclarationType, error) {
	err := common.GetDb().TmsWriteDb.Table(l.TableName()).Create(&customsDeclarationType).Error

	if err != nil {
		l.Context.Logger().Error(err.Error())
	}

	return customsDeclarationType, err
}

//TableName 数据库表名称
func (l *CustomsDeclarationTypeModel) TableName() string {
	return "customs_declaration_type"
}

//PrimaryKey 返回CustomsDeclarationType的主健
func (l *CustomsDeclarationTypeModel) PrimaryKey() string {
	return "type_id"
}

func (l *CustomsDeclarationTypeModel) getWhereString() string {
	return l.PrimaryKey() + "=?"
}

//CustomsDeclarationType 清关方式
type CustomsDeclarationType struct {
	TypeID uint64 `gorm:"primary_key" json:"type_id" form:"type_id"`
	Value  string `gorm:"value" json:"value" form:"value" binding:"required"`
}
