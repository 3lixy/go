package models

import (
	"azoya/nova"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
	// "strconv"
	"crypto/md5"
	"encoding/hex"
	"github.com/opentracing/opentracing-go"
	"tms_service/common"
	"tms_service/models/line"
	"tms_service/models/shipments"
	"strconv"
	"strings"
)

const (
	//ShipmentIDSplit 用来切割shipment ids
	ShipmentIDSplit = ","
)

//ShipmentsModel 用于读取shipments相关数据
type ShipmentsModel struct {
	Context *nova.Context
}

//NewShipmentsModel 初始化LineModel
func NewShipmentsModel(c *nova.Context) *ShipmentsModel {
	return &ShipmentsModel{Context: c}
}

//GetDetail 获取shipment的详细
func (s *ShipmentsModel) GetDetail(shipmentID uint64) (Shipments, error) {
	var shipment Shipments
	err := common.GetDb().TmsReadDb.Table(s.TableName()).
		Where(s.getWhereString(), shipmentID).
		Find(&shipment).
		Error

	if err != nil && err == gorm.ErrRecordNotFound {
		s.Context.Logger().Error(err.Error())
	}

	return shipment, err
}

//GetList 获取Shipments列表
func (s *ShipmentsModel) GetList(params map[string][]string,page common.Page) ([]Shipments, uint64, error) {
	sort := "entity_id"
	orderBy := "DESC"

	where := make(map[string]interface{})
	if _, ok := params["website_id"]; ok {
		websiteID, _ := strconv.ParseUint(params["website_id"][0], 10, 64)
		if websiteID > 0 {
			where["website_id"] = websiteID
		}
	}
	if _, ok := params["store_id"]; ok {
		store_id, _ := strconv.ParseUint(params["store_id"][0], 10, 64)
		where["store_id"] = store_id
	}
	if _, ok := params["warehouse_id"]; ok {
		warehouse_id, _ := strconv.ParseUint(params["warehouse_id"][0], 10, 64)
		where["warehouse_id"] = warehouse_id
	}
	if _, ok := params["transport_line_id"]; ok {
		transport_line_id, _ := strconv.ParseUint(params["transport_line_id"][0], 10, 64)
		where["transport_line_id"] = transport_line_id
	}
	if _, ok := params["status"]; ok {
		status:= params["status"][0]
		where["status"] = status
	}
	if _, ok := params["main_logistics_id"]; ok {
		main_logistics_id, _ := strconv.ParseUint(params["main_logistics_id"][0], 10, 64)
		where["main_logistics_id"] = main_logistics_id
	}
	if _, ok := params["oversea_logistics_id"]; ok {
		oversea_logistics_id, _ := strconv.ParseUint(params["oversea_logistics_id"][0], 10, 64)
		where["oversea_logistics_id"] = oversea_logistics_id
	}
	if _, ok := params["inland_logistics_id"]; ok {
		inland_logistics_id, _ := strconv.ParseUint(params["inland_logistics_id"][0], 10, 64)
		where["inland_logistics_id"] = inland_logistics_id
	}
	if _, ok := params["customs_declaration_type_id"]; ok {
		customs_declaration_type_id := params["customs_declaration_type_id"][0]
		where["customs_declaration_type_id"] = customs_declaration_type_id
	}

	var incrementID string
	var orderIncrementID string

	if _, ok := params["increment_id"]; ok {
		incrementID = params["increment_id"][0]
	}
	if _, ok := params["order_increment_id"]; ok {
		orderIncrementID = params["order_increment_id"][0]
	}

	var whereSecond string
	if incrementID != "" {
		incrementIDS := strings.Split(incrementID, ",")
		var newStr []string
		for _,s := range incrementIDS{
			newStr = append(newStr,strings.ToUpper(s))
		}
		incrementIDStr := "'" + strings.Join(newStr, "','") + "'"
		whereSecond += fmt.Sprintf(" and increment_id in (%s) ", incrementIDStr)
	}
	if orderIncrementID != "" {
		orderIncrementIDs := strings.Split(orderIncrementID, ",")
		var newStr []string
		for _,s := range orderIncrementIDs{
			newStr = append(newStr,strings.ToUpper(s))
		}
		orderIncrementIDStr := "'" + strings.Join(newStr, "','") + "'"
		whereSecond += fmt.Sprintf(" and order_increment_id in (%s) ", orderIncrementIDStr)
	}
	wherestr := ""
	if whereSecond != "" {
		wherestr = "1 = 1" + whereSecond
	}
	ss := common.GetParamsPage(page)
	var result []Shipments
	err := common.GetDb().TmsReadDb.
		Table(s.TableName()).
		Where(where).
		Where(wherestr).
		Offset(ss.Page).
		Limit(ss.Limit).
		Order(fmt.Sprintf("%s  %s", sort, orderBy)).
		Find(&result).
		Error
	var count uint64
	common.GetDb().TmsReadDb.
		Table(s.TableName()).
		Where(where).
		Where(wherestr).
		Count(&count)

	return result,count,err
}

//GetListByID 通过entity_id来获取shipment list
func (s *ShipmentsModel) GetListByID(shipmentIds []string) ([]Shipments, uint64, error) {
	var list []Shipments
	var count uint64
	err := common.GetDb().TmsReadDb.Table(s.TableName()).
		Where("entity_id in (?)", shipmentIds).
		Find(&list).
		Count(&count).
		Error

	return list, count, err
}

//GetTotalRowsCount 获取shipments的总数
func (s *ShipmentsModel) GetTotalRowsCount(params ShipmentsListQueryParams) (uint64, error) {
	var count uint64
	err := common.GetDb().TmsReadDb.Table(s.TableName()).
		Where(params).
		Count(&count).
		Error

	return count, err
}

//Create create shipments数据
func (s *ShipmentsModel) Create(shipment Shipments) (Shipments, error) {
	span, ctx := opentracing.StartSpanFromContext(s.Context.Context(), "Create shipments")
	s.Context.WithContext(ctx)

	var result Shipments
	err := common.GetDb().TmsReadDb.Table(s.TableName()).
		Where(&Shipments{
		WebsiteID: shipment.WebsiteID,
		OrderID:   shipment.OrderID,
	}).
		First(&result).
		Error

	if err != gorm.ErrRecordNotFound {
		repeatErr := errors.New(shipment.IncrementID + " already exist, don't repeat submit")
		s.Context.Logger().Error(repeatErr.Error())
		return shipment, repeatErr
	}

	tm := common.DefaultTime()
	shipment.CreatedAt = tm
	shipment.UpdatedAt = tm
	shipment.Status = shipments.PrepareAssignWarehouse
	shipment.IncrementID = common.GetIncrementID()

	str := fmt.Sprintf("%s", time.Now().Unix())
	h := md5.New()
	h.Write([]byte(str))
	res := h.Sum(nil)
	shipment.LabelIdentification = hex.EncodeToString(res)

	err = common.GetDb().TmsWriteDb.Table(s.TableName()).Create(&shipment).Error

	if err != nil {
		s.Context.Logger().Error(err.Error())
	}

	defer span.Finish()
	return shipment, err
}

//AssignWarehouse 把发货单分配到对应的仓库
func (s *ShipmentsModel) AssignWarehouse(shipmentIds []string, warehouseID string) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(s.Context.Context(), "Assign warehouse")
	s.Context.WithContext(ctx)

	list, _, err := s.GetListByID(shipmentIds)
	if err != nil {
		panic(err)
	}

	status := shipments.OK
	storeID := list[0].StoreID
	for _, shipment := range list {
		//仓库不为空不可分配，有不同店铺也不能分配
		if shipment.WarehouseID != 0 {
			errText := fmt.Sprintf("%s need to remove warehouse_id first", shipment.IncrementID)
			return shipments.ErrorNeedToRemoveWarehouse, errors.New(errText)
		}

		if shipment.StoreID != storeID {
			errText := fmt.Sprintf("need to select same store assign warehouse")
			return shipments.ErrorDifferentStore, errors.New(errText)
		}

		if shipment.Status != shipments.PrepareAssignWarehouse {
			errText := fmt.Sprintf("status need to be prepare assign warehouse")
			return shipments.ErrorStatus, errors.New(errText)
		}
	}

	err = common.GetDb().TmsWriteDb.Table(s.TableName()).
		Where("entity_id in (?)", shipmentIds).
		Updates(map[string]interface{}{"status": shipments.PrepareAssignTransport, "warehouse_id": warehouseID}).
		Error

	defer span.Finish()
	return status, err
}

//AssignTransport 把发货单分配到对应的仓库
func (s *ShipmentsModel) AssignTransport(shipmentIds []string, transportLineID uint64) (Line, int, error) {
	span, ctx := opentracing.StartSpanFromContext(s.Context.Context(), "Assign transport")
	s.Context.WithContext(ctx)

	lineDetail := Line{}
	list, _, err := s.GetListByID(shipmentIds)
	if err != nil {
		panic(err)
	}

	status := shipments.OK
	storeID := list[0].StoreID
	for _, shipment := range list {
		if shipment.Status != shipments.PrepareAssignTransport {
			errText := fmt.Sprintf("status need to be prepare assign transport line")
			return lineDetail, shipments.ErrorStatus, errors.New(errText)
		}
		if shipment.TransportLineID != 0 {
			errText := fmt.Sprintf("%s need to remove transport line first", shipment.IncrementID)
			return lineDetail, shipments.ErrorNeedToRemoveTransportLine, errors.New(errText)
		}

		if shipment.StoreID != storeID {
			errText := fmt.Sprintf("need to select same store assign warehouse")
			return lineDetail, shipments.ErrorDifferentStore, errors.New(errText)
		}
	}

	if err != nil {
		panic(err)
	}

	//获取物流线路信息
	lineModel := NewLineModel(s.Context)
	lineDetail, lineErr := lineModel.GetDetail(transportLineID)
	if lineErr != nil {
		errText := fmt.Sprintf("get transport line error, message: %s", lineErr.Error())
		return lineDetail, 0, errors.New(errText)
	}

	if lineDetail.Status != line.Enable {
		errText := fmt.Sprintf("this transport line status isn't enable")
		return lineDetail, 0, errors.New(errText)
	}

	//获取物流商信息
	logisticsModel := NewLogisticsModel(s.Context)
	mainLogistics, mainErr := logisticsModel.GetDetail(lineDetail.MainLogisticsID)

	overseaLogistics, oversearErr := logisticsModel.GetDetail(lineDetail.OverseaLogisticsID)

	inlandLogistics, inLandErr := logisticsModel.GetDetail(lineDetail.InlandLogisticsID)

	if mainErr != nil || oversearErr != nil || inLandErr != nil {
		panic(errors.New("transport line logistics not found, please check"))
	}

	//获取海关清关方式信息
	customsModel := NewCustomsDeclarationTypeModel(s.Context)
	customs, customsErr := customsModel.GetDetail(lineDetail.CustomsDeclarationType)

	if customsErr != nil {
		errText := fmt.Sprintf("get customs declaration type error, message: %s", customsErr.Error())
		panic(errText)
	}

	err = common.GetDb().TmsWriteDb.Table(s.TableName()).
		Where("entity_id in (?)", shipmentIds).
		Updates(map[string]interface{}{
		"status":                        shipments.AssignComplete,
		"transport_line_id":             transportLineID,
		"main_logistics_id":             mainLogistics.EntityID,
		"main_logistics_name":           mainLogistics.Name,
		"oversea_logistics_id":          overseaLogistics.EntityID,
		"oversea_logistics_name":        overseaLogistics.Name,
		"inland_logistics_id":           inlandLogistics.EntityID,
		"inland_logistics_name":         inlandLogistics.Name,
		"customs_declaration_type_id":   customs.TypeID,
		"customs_declaration_type_name": customs.Value,
	}).
		Error

	defer span.Finish()
	return lineDetail, status, err
}

//RemoveTransportLine 清除物流线路
func (s *ShipmentsModel) RemoveTransportLine(shipmentIds []string) (int, error) {
	list, _, err := s.GetListByID(shipmentIds)
	if err != nil {
		panic(err)
	}

	status := shipments.OK
	storeID := list[0].StoreID
	for _, shipment := range list {
		//有不同店铺也不能分配
		if shipment.Status != shipments.AssignComplete {
			errText := fmt.Sprintf("status need to be complete")
			return shipments.ErrorStatus, errors.New(errText)
		}

		if shipment.TransportLineID == 0 {
			errText := fmt.Sprintf("%s transport line ID do not existing", shipment.IncrementID)
			return shipments.ErrorTransportLineIDEmpty, errors.New(errText)
		}

		if shipment.StoreID != storeID {
			errText := fmt.Sprintf("need to select same store assign warehouse")
			return shipments.ErrorDifferentStore, errors.New(errText)
		}
	}

	err = common.GetDb().TmsWriteDb.Table(s.TableName()).
		Where("entity_id in (?)", shipmentIds).
		Updates(map[string]interface{}{
		"status":                        shipments.PrepareAssignTransport,
		"transport_line_id":             0,
		"main_logistics_id":             0,
		"main_logistics_name":           "",
		"oversea_logistics_id":          0,
		"oversea_logistics_name":        "",
		"inland_logistics_id":           0,
		"inland_logistics_name":         "",
		"customs_declaration_type_id":   0,
		"customs_declaration_type_name": "",
	}).
		Error

	return status, err
}

//TableName 数据库表名称
func (s *ShipmentsModel) TableName() string {
	return "order_shipments"
}

//PrimaryKey 返回Line的主健
func (s *ShipmentsModel) PrimaryKey() string {
	return "entity_id"
}

//defaultOrder 默认排序规则
func (s *ShipmentsModel) defaultOrder() string {
	return "desc"
}

func (s *ShipmentsModel) getWhereString() string {
	return s.PrimaryKey() + "=?"
}

//清除仓库
func (s *ShipmentsModel) RemoveWarehouse(shipmentIds []string) (int, error) {
	list, _, err := s.GetListByID(shipmentIds)
	if err != nil {
		panic(err)
	}

	status := shipments.OK
	for _, shipment := range list {
		//已分配物流不能分配
		if shipment.Status == shipments.AssignComplete {
			errText := fmt.Sprintf("%s 发货单已分配物流，请先清除物流", shipment.IncrementID)
			return shipments.ErrorStatus, errors.New(errText)
		}

		if shipment.TransportLineID > 0 {
			errText := fmt.Sprintf("%s 发货单已分配物流，请先清除物流", shipment.IncrementID)
			return shipments.ErrorNeedToRemoveTransportLine, errors.New(errText)
		}
	}

	err = common.GetDb().TmsWriteDb.Table(s.TableName()).
		Where("entity_id in (?)", shipmentIds).
		Updates(map[string]interface{}{
		"status":       shipments.PrepareAssignWarehouse,
		"warehouse_id": 0,
	}).Error
	return status, err
}

func (s *ShipmentsModel) IsSetWarehouse(shipmentIds []string) (int, error) {
	list, _, err := s.GetListByID(shipmentIds)
	if err != nil {
		panic(err)
	}

	status := shipments.OK
	storeID := list[0].StoreID
	for _, shipment := range list {
		//必须同一个店铺才可以分配仓库
		if shipment.StoreID != storeID {
			errText := fmt.Sprintf("必须同一个店铺才可以分配仓库")
			return shipments.ErrorDifferentStore, errors.New(errText)
		}
		//状态为待分配仓库才可以分配
		if shipment.Status != shipments.PrepareAssignWarehouse {
			errText := fmt.Sprintf("%s 状态为待分配仓库才可以分配", shipment.IncrementID)
			return shipments.ErrorStatus, errors.New(errText)
		}
	}
	return status, err
}

func (s *ShipmentsModel) GetShipments(params map[string][]string) (Shipments, error) {
	var result Shipments
	where := make(map[string]interface{})
	if _, ok := params["website_id"]; ok {
		websiteId, _ := strconv.ParseUint(params["website_id"][0], 10, 64)
		if websiteId > 0 {
			where["website_id"] = websiteId
		}
	}
	if _, ok := params["order_increment_id"]; ok {
		order_increment_id := params["order_increment_id"][0]
		where["order_increment_id"] = order_increment_id
	}
	if _, ok := params["label_identification"]; ok {
		label_identification := params["label_identification"][0]
		where["label_identification"] = label_identification
	}

	if _, ok := params["entity_id"]; ok {
		entity_id, _ := strconv.ParseUint(params["entity_id"][0], 10, 64)
		where["entity_id"] = entity_id
	}
	err := common.GetDb().TmsReadDb.Table(s.TableName()).Where(where).First(&result).Error
	return result, err
}

func (s *ShipmentsModel) GetShipment(params map[string]interface{}) (shipment Shipments, err error) {
	err = common.GetDb().TmsReadDb.Table(s.TableName()).
		Where(params).
		First(&shipment).Error
	return
}

//ShipmentsListQueryParams 用于查询list
type ShipmentsListQueryParams struct {
	WebsiteID uint64 `gorm:"website_id" json:"website_id" form:"website_id" binding:"required"`
	StoreID   uint64 `gorm:"store_id" json:"store_id" form:"store_id"`
	// OrderIncrementID         string `gorm:"order_increment_id" json:"order_increment_id" form:"order_increment_id"`
	// IncrementID              string `gorm:"increment_id" json:"increment_id" form:"increment_id"`
	WarehouseID              uint64 `gorm:"warehouse_id" json:"warehouse_id" form:"warehouse_id"`
	TransportLineID          uint64 `gorm:"transport_line_id" json:"transport_line_id" form:"transport_line_id"`
	Status                   string `gorm:"status" json:"status" form:"status"`
	MainLogisticsID          uint64 `gorm:"main_logistics_id" json:"main_logistics_id" form:"main_logistics_id"`
	OverseaLogisticsID       uint64 `gorm:"oversea_logistics_id" json:"oversea_logistics_id" form:"oversea_logistics_id"`
	InlandLogisticsID        uint64 `gorm:"inland_logistics_id" json:"inland_logistics_id" form:"inland_logistics_id"`
	CustomsDeclarationTypeID string `gorm:"customs_declaration_type_id" json:"customs_declaration_type_id" form:"customs_declaration_type_id"`
}

//Shipments 发货单
type Shipments struct {
	EntityID                   uint64 `gorm:"primary_key" json:"entity_id" form:"entity_id"`
	WebsiteID                  uint64 `gorm:"website_id" json:"website_id" form:"website_id" binding:"required"`
	StoreID                    uint64 `gorm:"store_id" json:"store_id" form:"store_id"`
	OrderID                    uint64 `gorm:"order_id" json:"order_id" form:"order_id" binding:"required"`
	OrderIncrementID           string `gorm:"order_increment_id" json:"order_increment_id" form:"order_increment_id" binding:"required"`
	IncrementID                string `gorm:"increment_id" json:"increment_id" form:"increment_id"`
	WarehouseID                uint64 `gorm:"warehouse_id" json:"warehouse_id" form:"warehouse_id"`
	TransportLineID            uint64 `gorm:"transport_line_id" json:"transport_line_id" form:"transport_line_id"`
	Status                     string `gorm:"status" json:"status" form:"status"`
	MainLogisticsID            uint64 `gorm:"main_logistics_id" json:"main_logistics_id" form:"main_logistics_id"`
	MainLogisticsName          string `gorm:"main_logistics_name" json:"main_logistics_name" form:"main_logistics_name"`
	OverseaLogisticsID         uint64 `gorm:"oversea_logistics_id" json:"oversea_logistics_id" form:"oversea_logistics_id"`
	OverseaLogisticsName       string `gorm:"oversea_logistics_name" json:"oversea_logistics_name" form:"oversea_logistics_name"`
	InlandLogisticsID          uint64 `gorm:"inland_logistics_id" json:"inland_logistics_id" form:"inland_logistics_id"`
	InlandLogisticsName        string `gorm:"inland_logistics_name" json:"inland_logistics_name" form:"inland_logistics_name"`
	CustomsDeclarationTypeName string `gorm:"customs_declaration_type_name" json:"customs_declaration_type_name" form:"customs_declaration_type_name"`
	StartCountry               string `gorm:"start_country" json:"start_country" form:"start_country"`
	TargetCountry              string `gorm:"target_country" json:"target_country" form:"target_country"`
	LabelIdentification        string `gorm:"label_identification" json:"label_identification" form:"label_identification"`
	CreatedAt                  string `gorm:"created_at" json:"created_at" form:"created_at"`
	UpdatedAt                  string `gorm:"updated_at" json:"updated_at" form:"updated_at"`
}
