package config

import (
    "github.com/astaxie/beego/config"
    "path/filepath"
    "strings"
    "os"
)

var configFileName = "config.conf"

// 配置文件的默认值或配置key值
const (
	Enable = "enable"
	Disable = "disable"

	DefaultServiceName = "service_example"
	DefaultMonitorStatus = "enable"
	DefaultMetricsAuthStatus = "disable"

	ServiceNameKey = "service::name"
	MonitorStatusKey = "monitor::status"
    MetricsAuthStatusKey = "metrics::auth_status"
    MetricsAuthKey = "auth"
	MetricsAuthTokenKey = "metrics::auth_token"
)

// Configer 日志对象，用于解析配置文件，输出数据
type Configer interface {
    config.Configer
}

// NewConfig 初始化一个配置文件对象，在应用启动时就加载到内存中
func NewConfig(adapterName string) (Configer){
    iniconf, err := config.NewConfig(adapterName, currentPath())
    if err != nil {
        panic(err)
    }
    return iniconf
}

func currentPath() string {
    appPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        panic(err)
    }
    configPath := filepath.Join(appPath,configFileName)
    return strings.Replace(configPath,"\\","/",-1)
}

// MonitorStatus 获取所有监控启用状态
func MonitorStatus(config Configer) bool {
    if config.DefaultString(MonitorStatusKey,DefaultMonitorStatus) == Enable {
		return true
    } else {
        return false
    }
}

// MetricsAuthStatus 获取metrics监控启用状态
func MetricsAuthStatus(config Configer) bool {
    if config.DefaultString(MetricsAuthStatusKey,DefaultMetricsAuthStatus) == Enable {
		return true
    } else {
        return false
    }
}

// MetricsToken 获取metrics监控接口的验证token
func MetricsToken(config Configer) string {
    return config.String(MetricsAuthTokenKey)
}
