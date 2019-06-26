package common

import (
	"errors"
)

var (
	ErrStockNameIsExist = errors.New("仓库名称已存在")
	ErrStockNotExist    = errors.New("仓库不存在")
	ErrDateFormat       = errors.New("时间格式错误")
	ErrStoreNotExist    = errors.New("店铺不存在")
)
