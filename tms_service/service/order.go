package service

import (
	"azoya/nova"
	"fmt"
	"tms_service/common"
	"net/http"
	"errors"
	"net/url"
	"strings"
)

//OrderService 模型
type OrderService struct {
	BaseService
}

//NewOrderService 模型
func NewOrderService(c *nova.Context) *OrderService {
	return &OrderService{BaseService{C: c}}
}

const (
	//ORDERCANCELRETURNSTOCK 是否释放库存1是0否
	ORDERCANCELRETURNSTOCK = "1"
	//DOWNLOADEDYES 商户拉取是否拉取订单是yes，否no
	DOWNLOADEDYES = "yes"
	//TLC同步状态 1成功
	ISSYNCYES = 2
	//FTP同步状态 1成功
	FTPSTATUSYES = 1
	NEW          = "new"
	//PENDDINGSTATUS 未付款
	PENDDINGSTATUS = "pending"
	//PROCESSINGSTATUS 已付款
	PROCESSINGSTATUS = "processing"
	//SHIPMENTSSTATUS 已发货
	SHIPMENTSSTATUS = "shipments"
	//CANCELEDSTATUS 已取消
	CANCELEDSTATUS = "canceled"
	//CHANGEORDERSTATUS 改变订单状态
	CHANGEORDERSTATUS = "change order status"
	//OPERATOR 操作者
	OPERATOR = "operator"

	//OrderTypeSame = 0
	OrderTypeSame = 0
	//OrderTypeCut order_type 砍单新建
	OrderTypeCut = 2
	//ORDERTYPEOOS 缺货订单
	OrderTypeOos = 3

	//身份证信息 已验证
	IdentityVerified = "verified"
	//身份证信息 未验证
	IdentityUnverified = "unverified"

	EXPORTLIMIT = 10000
)

type OrderSearch struct {
	EntityID             uint64  `gorm:"entity_id" json:"entity_id"`
	IncrementID          string  `gorm:"increment_id" json:"increment_id"`
	CustomerID           uint64  `gorm:"customer_id" json:"customer_id"`
	StockID              string  `gorm:"stock_id" json:"stock_id"`
	GrandTotal           float64 `gorm:"grand_total" json:"grand_total"`
	Status               string  `gorm:"status" json:"status"`
	IdentityStatus       string  `gorm:"identity_status" json:"identity_status"`
	PromotionSource      string  `gorm:"promotion_source" json:"promotion_source"`
	CreatedAt            string  `gorm:"created_at" json:"created_at"`
	UpdatedAt            string  `gorm:"updated_at" json:"updated_at"`
	PaidAt               string  `gorm:"paid_at" json:"paid_at"`
	Firstname            string  `gorm:"firstname" json:"firstname"`
	ParentID             uint64  `gorm:"parent_id" json:"parent_id"`
	OrderSource          string  `gorm:"order_source" json:"order_source"`
	OriginalSymbol       string  `gorm:"original_symbol" json:"original_symbol"`
	OrderCurrencyCode    string  `gorm:"order_currency_code" json:"order_currency_code"`
	OriginalCurrencyCode string  `gorm:"original_currency_code" json:"original_currency_code"`
	OriginalExchangeRate float64 `gorm:"original_exchange_rate" json:"original_exchange_rate"`
	RelationParentID     string  `gorm:"relation_parent_id" json:"relation_parent_id"`
	TotalFee             float64 `gorm:"total_fee" json:"total_fee"`
	OrderType            int     `gorm:"order_type" json:"order_type"`
	Level                uint64  `gorm:"level" json:"level"`
	ThreePartOrderNumber string  `gorm:"three_part_order_number" json:"three_part_order_number"`
	CompletedAt          string  `gorm:"completed_at" json:"completed_at"`
	CanceledAt           string  `gorm:"canceled_at" json:"canceled_at"`
	SupplierOrderNo      string  `gorm:"supplier_order_no" json:"supplier_order_no"`
	ShippingType         string  `gorm:"shipping_type" json:"shipping_type"`
	VendorID             uint64  `gorm:"vendor_id" json:"vendor_id"`
}

type OrderSearchResult struct {
	Status  uint64          `json:"status"`
	Message string          `json:"message"`
	Data    OrderSearchList `json:"data"`
}

type OrderSearchList struct {
	Rows  []OrderSearch `json:"rows"`
	Total int64         `json:"total"`
}
func (o *OrderService) OrderQuickSearch(websiteID uint64, params map[string]interface{}, page uint64, limit uint64) (orderSearchList []OrderSearch, err error) {
	var result OrderSearchResult
	c := common.GetConfig()
	orderURL := c.String("oms::order_quick_search")
	v := url.Values{}
	v.Set("website_id", fmt.Sprintf("%v", websiteID))
	for key, value := range params {
		v.Set(key, fmt.Sprintf("%v", value))
	}
	v.Set("page", fmt.Sprintf("%v", page))
	v.Set("limit", fmt.Sprintf("%v", limit))
	endpoint := "order/search"
	err = o.C.Client().GetJSON(o.C,endpoint,orderURL+"?"+v.Encode(), &result)
	if err != nil {
		fmt.Println(err)
	}
	if result.Status == http.StatusOK {
		if len(result.Data.Rows) > 0 {
			orderSearchList = result.Data.Rows
		} else {
			err = common.ErrGetOrder
		}
	} else {
		err = errors.New(common.InterfaceToString(result.Message))
	}
	return
}

func (o *OrderService) OrderStatusToShipment(websiteID uint64, orderID uint64) (err error) {
	c := common.GetConfig()
	cancelURL := c.String("oms::status_to_shipment")
	ssoUserName := c.String("oms::sso_user_name")
	params := url.Values{}
	params.Set("website_id", fmt.Sprintf("%v", websiteID))
	params.Set("order_id", fmt.Sprintf("%v", orderID))
	params.Set("sso_user_id", fmt.Sprintf("%v", 1))
	params.Set("sso_user_name", ssoUserName)
	vcode := params.Encode()
	payload := strings.NewReader(vcode)
	endpoint := "order/shipment"
	var result Result
	err = o.C.Client().PostJSON(o.C,endpoint,cancelURL, &result,payload)
	if err != nil {
		fmt.Println(err)
	}
	if result.Status != http.StatusOK {
		err = errors.New(common.InterfaceToString(result.Message))
	}
	return
}
