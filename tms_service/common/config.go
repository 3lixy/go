package common

import (
	"azoya/nova/config"
)

var _config config.Configer

//初始化配置文件。。
func Init(c config.Configer) bool {

	if c == nil {
		return false
	}
	_config = c
	return true
}

// GetConfig 获取配置文件信息
func GetConfig() config.Configer {
	return _config
}
