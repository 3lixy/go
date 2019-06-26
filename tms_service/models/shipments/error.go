package shipments

var (
	//OK 表示成功
	OK = 200

	//ErrorNeedToRemoveWarehouse 发货单已经有仓库，需要删除再分配
	ErrorNeedToRemoveWarehouse = 501

	//ErrorNeedToRemoveTransportLine 发货单已经有物流线路，需要删除再分配
	ErrorNeedToRemoveTransportLine = 502

	//ErrorDifferentStore 这批发货单中有不相同的店铺，需要相同才能分配
	ErrorDifferentStore = 503

	//ErrorTransportLineIDEmpty 物流线路id不存在
	ErrorTransportLineIDEmpty = 504

	//ErrorStatus 状态错误
	ErrorStatus = 505
)
