package common

import (
	"azoya/lib/log"
	"azoya/nova"
)

var (
	//日志类型

	//LogTypeInfo 提示
	LogTypeInfo = "info"

	//LogTypeWarning 警告
	LogTypeWarning = "warning"

	//LogTypeError 错误
	LogTypeError = "error"

	//LogTypeDebug 调试
	LogTypeDebug = "debug"

	//logFileExtended 日志文件后缀
	logFileExtended = ".log"
)
var _logger *log.Logger

func InitLogger() {
	var err error
	var cfg *log.Config
	mode := GetConfig().String("service::runmode")
	if mode == nova.ReleaseMode {
		cfg = log.ProductionConfig()
	} else {
		cfg = log.DevelopmentConfig()
	}

	cfg.Program = GetConfig().String("service::name")

	_logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}

// 刷新已缓冲的日志
// 在应用程序退出之前应注意调用Sync
func SyncLogger() {
	_logger.Sync()
}

func GetLogger() *log.Logger {
	return _logger
}
