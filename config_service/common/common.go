package common

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	//LIMIT 分页数量
	LIMIT uint64 = 10

	//PAGE 默认页数
	PAGE uint64 = 1

	//MsgSuccess 接口返回成功
	MsgSuccess = "success"

	//MsgFailed 接口返回失败
	MsgFailed = "failed"

	//timeFormat默认日期格式
	timeFormat = "2006-01-02 15:04:05"
)

//Page 分页
type Page struct {
	Page  uint64
	Limit uint64
}

//Msg 提示
type Msg struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

//DefaultTime 获取默认时间
func DefaultTime(zoneType string) string {
	switch zoneType {
	case "UTC":
		return time.Now().UTC().Format(timeFormat)
	case "PRC":
		return time.Now().Format(timeFormat)
	default:
		return time.Now().UTC().Format(timeFormat)
	}
}

//Message 返回接口响应数组
func Message(status int, msg interface{}, data interface{}) Msg {
	return Msg{
		Status:  status,
		Message: msg,
		Data:    data,
	}
}

//GetParamsPage 获取页码
func GetParamsPage(params map[string][]string) Page {
	var page Page
	page.Page = PAGE
	page.Limit = LIMIT
	if _, ok := params["limit"]; ok {
		paramsLimit, _ := strconv.ParseUint(params["limit"][0], 0, 0)
		if paramsLimit < 5000 {
			page.Limit = paramsLimit
		}
	}

	if _, ok := params["page"]; ok {
		paramsPage, _ := strconv.ParseUint(params["page"][0], 0, 0)
		page.Page = paramsPage
	}

	page.Page = (page.Page - 1) * page.Limit
	return page
}

//GetPageLimit 获取页码
func GetPageLimit(pageNumber uint64, limit uint64) Page {
	var page Page
	page.Page = PAGE
	page.Limit = LIMIT
	if limit > 0 && limit < 5000 {
		page.Limit = limit
	}
	if pageNumber > 0 {
		page.Page = pageNumber
	}

	page.Page = (page.Page - 1) * page.Limit
	return page
}

//GetValidPage 根据查询总条数，获取页码
func GetValidPage(params map[string]string) Page {
	var page Page
	page.Page = PAGE
	page.Limit = LIMIT
	if _, ok := params["limit"]; ok {
		paramsLimit, _ := strconv.ParseUint(params["limit"], 0, 0)
		if paramsLimit < 10001 {
			page.Limit = paramsLimit
		}
	}
	if _, ok := params["page"]; ok {
		paramsPage, _ := strconv.ParseUint(params["page"], 0, 0)
		page.Page = paramsPage
	}
	page.Page = (page.Page - 1) * page.Limit
	return page
}

//GetSameDayTime 获取当天日期
func GetSameDayTime() string {
	return time.Now().Format("20060102")
}

//GetAppointMonthTime 获取指定的时间
func GetAppointMonthTime(timeString string) (start string, end string, err error) {

	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation("2006-01", timeString, loc)

	if err != nil {
		return
	}
	y, m, _ := theTime.Date()
	thisMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
	start = thisMonth.AddDate(0, 0, 0).UTC().Format("2006-01-02 15:04:05")
	end = thisMonth.AddDate(0, +1, 0).UTC().Format("2006-01-02 15:04:05")
	return
}

//GetMonthTime 获取上个月 年月份时间
func GetMonthTime() string {

	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	return thisMonth.AddDate(0, -3, 0).Format("2006-01")
}

//GetSpecificTime() 将GMT时间转换成北京时间
func GetSpecificTime(timeAt string) (string, error) {
	if timeAt == "" || timeAt == "0000-00-00 00:00:00" {
		return "", nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", timeAt)
	if err != nil {
		return "", err
	}
	h, _ := time.ParseDuration("8h")
	return t.Add(h).Format("2006-01-02 15:04:05"), nil
}

//JsonTomap json转换map
func JsonTomap(jsonString string) (map[string]string, error) {
	var result map[string]string
	if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
		return result, err
	}
	return result, nil

}

//StructTomap json转换map
//Struct interface
func StructTomap(Struct interface{}) (map[string]string, error) {
	var result map[string]string
	jsonVal, err := json.Marshal(Struct)
	if err != nil {
		return result, err
	}
	result, err = JsonTomap(fmt.Sprintf("%s", jsonVal))
	return result, err
}

//获取指定月份账期的开始结束时间（精确到时分秒）
func GetSpecificMonthTime(timeString string) (start string, end string, err error) {

	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation("2006-01", timeString, loc)

	if err != nil {
		return
	}
	y, m, _ := theTime.Date()
	thisMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
	start = thisMonth.UTC().Format(timeFormat)
	end = thisMonth.AddDate(0, +1, 0).Add(time.Second * -1).UTC().Format("2006-01-02 15:04:05")
	return
}

//GetLastMonth 获取上个月份
func GetLastMonth(timeString string) string {
	var lastMonth string
	loc, _ := time.LoadLocation("Local")
	timeFormatYearMonth := "2006-01"
	theTime, _ := time.ParseInLocation(timeFormatYearMonth, timeString, loc)
	y, m, _ := theTime.Date()
	thisMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
	lastMonth = thisMonth.UTC().Format(timeFormatYearMonth)
	return lastMonth
}

//GetDateFormat //获取时间格式
func GetDateFormat(timeAt string) (string, error) {

	if regexp.MustCompile(`^\d{4}-\d{1}-\d{1}\s*\d{1,2}:\d{1,2}`).MatchString(timeAt) {
		return "2006-1-2 15:04:05", nil
	}

	if regexp.MustCompile(`^\d{4}-\d{2}-\d{1}\s*\d{1,2}:\d{1,2}`).MatchString(timeAt) {
		return "2006-01-2 15:04:05", nil
	}

	if regexp.MustCompile(`^\d{4}-\d{1}-\d{2}\s*\d{1,2}:\d{1,2}`).MatchString(timeAt) {
		return "2006-1-02 15:04:05", nil
	}

	if regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\s*\d{1,2}:\d{1,2}`).MatchString(timeAt) {
		return "2006-01-02 15:04:05", nil
	}

	if regexp.MustCompile(`^\d{4}[/]\d{1}[/]\d{1}\s*\d{1,2}:\d{1,2}`).MatchString(timeAt) {
		return "2006/1/2 15:04:05", nil
	}

	if regexp.MustCompile(`^\d{4}[/]\d{2}[/]\d{1}\s*\d{1,2}:\d{1,2}`).MatchString(timeAt) {
		return "2006/01/2 15:04:05", nil
	}

	if regexp.MustCompile(`^\d{4}[/]\d{1}[/]\d{2}\s*\d{1,2}:\d{1,2}`).MatchString(timeAt) {
		return "2006/1/02 15:04:05", nil
	}

	if regexp.MustCompile(`^\d{4}[/]\d{2}[/]\d{2}\s*\d{1,2}:\d{1,2}`).MatchString(timeAt) {
		return "2006/01/02 15:04:05", nil
	}
	return "", ErrDateFormat
}

//ConversionTime 北京时间转化UTC时间
func ConversionTime(conversionTime string, format string) string {
	conversionTimeArr := strings.Split(conversionTime, ":")
	if len(conversionTimeArr) == 2 {
		conversionTime = conversionTime + ":00"
	}
	loc, _ := time.LoadLocation("Local")
	dt, _ := time.ParseInLocation(format, conversionTime, loc)
	return dt.Add(time.Hour * -8).Format(format)
}

//CheckCurrencyCode 校验币种
func CheckCurrencyCode(code string) bool {
	currencyCode := []string{
		"AED", "AFN", "ALL", "AMD", "ANG", "AOA", "ARS", "AUD", "AWG", "AZN",
		"BAM", "BBD", "BDT", "BGN", "BHD", "BIF", "BMD", "BND", "BOB", "BRL",
		"BSD", "BTC", "BTN", "BWP", "BYR", "BZD", "CAD", "CDF", "CHF", "CLF",
		"CLP", "CNY", "COP", "CRC", "CUP", "CVE", "CZK", "DJF", "DKK", "DOP",
		"DZD", "EEK", "EGP", "ERN", "ETB", "EUR", "FJD", "FKP", "GBP", "GEL",
		"GGP", "GHS", "GIP", "GMD", "GNF", "GTQ", "GYD", "HKD", "HNL", "HRK",
		"HTG", "HUF", "IDR", "ILS", "IMP", "INR", "IQD", "IRR", "ISK", "JEP",
		"JMD", "JOD", "JPY", "KES", "KGS", "KHR", "KMF", "KPW", "KRW", "KWD",
		"KYD", "KZT", "LAK", "LBP", "LKR", "LRD", "LSL", "LTL", "LVL", "LYD",
		"MAD", "MDL", "MGA", "MKD", "MMK", "MNT", "MOP", "MRO", "MTL", "MUR",
		"MVR", "MWK", "MXN", "MYR", "MZN", "NAD", "NGN", "NIO", "NOK", "NPR",
		"NZD", "OMR", "PAB", "PEN", "PGK", "PHP", "PKR", "PLN", "PYG", "QAR",
		"RON", "RSD", "RUB", "RWF", "SAR", "SBD", "SCR", "SDG", "SEK", "SGD",
		"SHP", "SLL", "SOS", "SRD", "STD", "SVC", "SYP", "SZL", "THB", "TJS",
		"TMT", "TND", "TOP", "TRY", "TTD", "TWD", "TZS", "UAH", "UGX", "USD",
		"UYU", "UZS", "VEF", "VND", "VUV", "WST", "XAF", "XAG", "XAU", "XCD",
		"XDR", "XOF", "XPF", "YER", "ZAR", "ZMK", "ZMW", "ZWL"}
	for _, v := range currencyCode {
		if v == code {
			return true
		}
	}
	return false
}

//Round float64 四舍五入
func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}
