package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

const (
	//DBMASTER 主库
	DBMASTER = "master"
	//DBSLAVE 从库
	DBSLAVE = "slave"
	//DBALIASNAME db别名
	DBALIASNAME = "db"
)

//dbConnector 数据配置
type dbConnector struct {
	DriverName   string //驱动名称
	UserName     string //用户名
	PassWord     string //密码
	Host         string //主机ip
	Port         int    //端口
	Dbname       string //数据库名称
	MaxIdleConns int    //最大空闲连接数
	MaxOpenConns int    //最大连接数
	DbMode       bool   //日志记录器
	MaxLifetime  int    //最大生存时间
}

//AllDBInstance WriteDb主库 ReadDb从库
type AllDBInstance struct {
	WriteDb *gorm.DB
	ReadDb  *gorm.DB
}

//全局对象
var (
	conns *AllDBInstance
)

// InitDb 初始化db连接
func InitDb() {
	conns = &AllDBInstance{
		WriteDb: connDb(DBMASTER),
		ReadDb:  connDb(DBSLAVE),
	}
}

//GetDb 获取db
func GetDb() *AllDBInstance {
	return conns
}

//ConnDb 操作实例
func connDb(dbType string) *gorm.DB {
	c := GetConfig()
	//配置站点 db+master/slave
	db := DBALIASNAME + "_" + dbType
	//驱动名称
	driverName := c.String(db + "::drivername")
	//端口号
	port, _ := c.Int(db + "::port")
	//最大空闲连接数
	maxIdleConns, _ := c.Int("service::maxidleconns")
	//最大连接
	maxOpenConns, _ := c.Int("service::maxopenconns")
	//日志记录器
	dbMode, _ := c.Bool("service::dbmode")
	//最大生存时间
	maxLifetime, _ := c.Int("service::maxlifetime")

	dbconnector := &dbConnector{
		DriverName:   driverName,
		UserName:     c.String(db + "::username"),
		PassWord:     c.String(db + "::password"),
		Host:         c.String(db + "::host"),
		Port:         port,
		Dbname:       c.String(db + "::dbname"),
		MaxIdleConns: maxIdleConns,
		MaxOpenConns: maxOpenConns,
		DbMode:       dbMode,
		MaxLifetime:  maxLifetime,
	}
	return conn(dbconnector)
}

//conn 连接db实例
func conn(dbconn *dbConnector) *gorm.DB {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbconn.UserName, dbconn.PassWord, dbconn.Host, dbconn.Port, dbconn.Dbname)
	db, err := gorm.Open(dbconn.DriverName, dns)
	if err != nil {
		log.Println("conndb err:", err)
		panic(err)
	}
	//If n <= 0, no idle connections are retained
	if dbconn.MaxIdleConns >= 0 {
		db.DB().SetMaxIdleConns(dbconn.MaxIdleConns)
	}
	// If n <= 0, then there is no limit on the number of open connections.
	if dbconn.MaxOpenConns > 0 {
		db.DB().SetMaxOpenConns(dbconn.MaxOpenConns)
	}
	if dbconn.DbMode {
		db.LogMode(dbconn.DbMode)
	}
	// If d <= 0, connections are reused forever.
	if dbconn.MaxLifetime >= 0 {
		d := time.Duration(dbconn.MaxLifetime) * time.Second
		db.DB().SetConnMaxLifetime(d)
	}
	return db
}
