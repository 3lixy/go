package common

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	KEY_NUM                 int    = 5
	LIMIT                   uint64 = 10
	PAGE                    uint64 = 1
	COUNT                   uint64 = 1
	HTTP_STATUS_OK                 = 200 //成功
	HTTP_STATUS_BAD_REQUEST        = 400 //异常
)

//Page 分页
type Page struct {
	Page       uint64 `json:"page" form:"page"`
	Limit      uint64 `json:"limit" form:"limit"`
	Sort       string `gorm:"sort" json:"sort" form:"sort"`
	Order      string `gorm:"order" json:"order" form:"order"`
	TotalCount uint64 `json:"total_count"`
}

//Msg 提示
type Msg struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

//DefaultTime 获取默认时间
func DefaultTime() string {
	var createdAt string = time.Now().UTC().Format("2006-01-02 15:04:05")
	return createdAt
}

//ErrorMessage 返回接口错误响应数组
func ErrorMessage(status int, msg interface{}) Msg {
	return Msg{
		Status:  status,
		Message: msg,
	}
}

//Message 返回接口响应数组
func Message(status int, msg string, data interface{}) Msg {
	return Msg{
		Status:  status,
		Message: msg,
		Data:    data,
	}
}

//ReduceTime 获取UTC 减去N天的时间
func ReduceTime(dayNum int) string {
	var createdAt string = time.Now().UTC().AddDate(0, 0, -dayNum).Format("2006-01-02 15:04:05")
	return createdAt
}

//GetRandomString 生成随机字符串
func GetRandomString(lens int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lens; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//Substr 截取
//start：正数 - 在字符串的指定位置开始,超出字符串长度强制把start变为字符串长度
//       负数 - 在从字符串结尾的指定位置开始
//       0 - 在字符串中的第一个字符处开始
//length:正数 - 从 start 参数所在的位置返回
//       负数 - 从字符串末端返回

func Substr(str string, start, length int) string {
	if length == 0 {
		return ""
	}
	runeStr := []rune(str)
	lenStr := len(runeStr)

	if start < 0 {
		start = lenStr + start
	}
	if start > lenStr {
		start = lenStr
	}
	end := start + length
	if end > lenStr {
		end = lenStr
	}
	if length < 0 {
		end = lenStr + length
	}
	if start > end {
		start, end = end, start
	}
	return string(runeStr[start:end])
}

//GetValidPage 根据查询总条数，获取页码
func GetValidPage(params Page) Page {

	page := GetParamsPage(params)
	page.TotalCount = params.TotalCount

	if uint64(params.TotalCount) <= page.Page {
		page.Page = 0
	}
	return page
}

//GetParamsPage 获取页码
func GetParamsPage(params Page) Page {

	var page Page

	page.Page = PAGE
	page.Limit = LIMIT

	if params.Limit > 0 {
		if params.Limit < 5000 {
			page.Limit = params.Limit
		}
	}

	if params.Page > 0 {
		page.Page = params.Page
	}

	page.Page = (page.Page - 1) * page.Limit
	return page
}

//DeepsCopy 拷贝MAP数组
func DeepsCopy(value map[string][]string) map[string][]string {
	newMap := make(map[string][]string)
	for k, v := range value {
		newMap[k] = v
	}
	return newMap
}

//ChangeMap 将字符串转换成MAP
func ChangeMap(str string) map[string]string {

	itemSlice := strings.Split(str, ",")
	itemMap := make(map[string]string)
	for _, v := range itemSlice {
		itemMap[v] = v
	}
	return itemMap
}

//StructToMap 结构体转化map
func StructToMap(obj interface{}) []string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var params []string
	for i := 0; i < t.NumField(); i++ {
		params = append(params, InterfaceToString(v.Field(i).Interface()))
	}

	return params
}

//InterfaceToString 类型转换string类型
func InterfaceToString(inter interface{}) string {

	tempStr := ""
	switch inter.(type) {
	case string:
		tempStr = inter.(string)
		break
	case float64:
		tempStr = strconv.FormatFloat(inter.(float64), 'f', -1, 64)
		break
	case int64:
		tempStr = strconv.FormatInt(inter.(int64), 30)
		break
	case int:
		tempStr = strconv.Itoa(inter.(int))
		break
	case uint64:
		tempStr = fmt.Sprintf("%d", inter)
		break
	}
	return tempStr
}

//正则匹配校验手机号
func CheckPhone(phone string) bool {
	rgx := regexp.MustCompile(`^1[0-9]{10}$`)
	return rgx.MatchString(phone)

}

//正则匹配校验邮箱
func CheckEmail(email string) bool {
	rgx := regexp.MustCompile(`^[_a-z0-9-]+(\.[_a-z0-9-]+)*@[a-z0-9-]+(\.[a-z0-9-]+)*(\.[a-z]{2,})$`)
	return rgx.MatchString(email)

}

//正则匹配身份证
func CheckIdentity(identity string) bool {
	if identity == "" {
		return false
	}
	rgx := regexp.MustCompile(`^[1-9]\d{5}[1-9]\d{3}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{3}([0-9]|X)$`)
	return rgx.MatchString(identity)

}

//RoundTo 四舍五入
func RoundTo(f float64, n int) float64 {

	if math.IsNaN(f) == true {
		return 0
	}

	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}

//GetIncrementID 生成发货单的incrementID
func GetIncrementID() string {
	dateTime := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%d", dateTime, 10000+time.Now().Nanosecond()/100000)
}
