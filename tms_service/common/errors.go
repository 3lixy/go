package common

import (
	"errors"
)

var (
	ErrShipment              = errors.New("发货单已分配物流，请先清除物流")
	ErrWarehouse             = errors.New("仓库已分配")
	ErrStore                 = errors.New("勾选的店铺不一致")
	ErrGetOrder              = errors.New("获取订单信息错误")
	ErrOrderShipmentNotExist = errors.New("发货单不存在")
	ErrLogisticsNotExist     = errors.New("物流商不存在")
)
