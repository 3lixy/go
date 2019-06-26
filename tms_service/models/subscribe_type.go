package models

import (
	"azoya/nova"
	"github.com/jinzhu/gorm"
	"tms_service/common"
)

//SubscribeTypeModel 用于读取SubscribeType相关数据
type SubscribeTypeModel struct {
	Context *nova.Context
}

//NewSubscribeTypeModel 初始化SubscribeTypeModel
func NewSubscribeTypeModel(c *nova.Context) *SubscribeTypeModel {
	return &SubscribeTypeModel{Context: c}
}

//GetList 获取SubscribeType列表
func (l *SubscribeTypeModel) GetList() ([]SubscribeType, error) {
	var list []SubscribeType
	err := common.GetDb().TmsReadDb.Table(l.TableName()).
		Find(&list).
		Error

	return list, err
}

//GetDetail 获取SubscribeType详情
func (l *SubscribeTypeModel) GetDetail(subscribeTypeID uint64) (SubscribeType, error) {
	var subscribeType SubscribeType
	err := common.GetDb().TmsReadDb.Table(l.TableName()).
		Where(l.getWhereString(), subscribeTypeID).
		Find(&subscribeType).
		Error

	if err != nil && err == gorm.ErrRecordNotFound {
		l.Context.Logger().Error(err.Error())
	}

	return subscribeType, err
}

//Delete 删除SubscribeType
func (l *SubscribeTypeModel) Delete(subscribeTypeID uint64) error {
	var subscribeType SubscribeType

	subscribeType, err := l.GetDetail(subscribeTypeID)
	if err != nil {
		return err
	}

	deleteErr := common.GetDb().TmsWriteDb.Table(l.TableName()).
		Where(l.getWhereString(), subscribeTypeID).
		Delete(&subscribeType).
		Error

	if deleteErr != nil {
		l.Context.Logger().Error(err.Error())
	}

	return deleteErr
}

//Create create SubscribeType数据
func (l *SubscribeTypeModel) Create(subscribeType SubscribeType) (SubscribeType, error) {
	err := common.GetDb().TmsWriteDb.Table(l.TableName()).Create(&subscribeType).Error

	if err != nil {
		l.Context.Logger().Error(err.Error())
	}

	return subscribeType, err
}

//TableName 数据库表名称
func (l *SubscribeTypeModel) TableName() string {
	return "subscribe_type"
}

//PrimaryKey 返回SubscribeType的主健
func (l *SubscribeTypeModel) PrimaryKey() string {
	return "type_id"
}

func (l *SubscribeTypeModel) getWhereString() string {
	return l.PrimaryKey() + "=?"
}

//SubscribeType api订阅方式
type SubscribeType struct {
	TypeID uint64 `gorm:"primary_key" json:"type_id" form:"type_id"`
	Value  string `gorm:"value" json:"value" form:"value" binding:"required"`
}
