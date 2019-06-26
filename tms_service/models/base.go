package models

import (
	"azoya/nova"
)

//BaseModel 模型 Db代表magento数据库db对象(分主库和从库) tmsDb指tms数据库db对象(分主库和从库)
type BaseModel struct {
	C *nova.Context
}
