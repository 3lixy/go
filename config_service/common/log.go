package common

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"time"
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

//Log 自定义文件写入
//context 写入内容
//name 生成文件名称
//logType 日志类型：info,error,warning,debug
func Log(context string, name string, logType string) {
	dir := time.Now().Format("2006-01-02")
	path := "logs/" + dir
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterFile, fmt.Sprintf(`{"filename":"logs/%s/%s","level":7,"maxlines":0,"maxsize":0,"daily":false,"maxdays":10}`, dir, name+logFileExtended))
	switch logType {
	case LogTypeError:
		log.Error(context)
	case LogTypeWarning:
		log.Warning(context)
	case LogTypeInfo:
		log.Info(context)
	case LogTypeDebug:
		log.Debug(context)
	default:
		log.Info(context)
	}
}
