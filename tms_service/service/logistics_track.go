package service

import (
	"azoya/nova"
	"fmt"
	"tms_service/models"
	"tms_service/validate"
)

//LogisticsTrackService 模型
type LogisticsTrackService struct {
	BaseService
}

//NewLogisticsTrackService 模型
func NewLogisticsTrackService(c *nova.Context) *LogisticsTrackService {
	return &LogisticsTrackService{BaseService{C: c}}
}

//获取物流轨迹
func (t *LogisticsTrackService) GetTrackItemList(query validate.TrackItemListQuery) (list []models.TrackItem) {
	subscribeRecordModels := models.NewSubscribeRecordModel(t.C)
	trackItemModels := models.NewTrackItemModel(t.C)
	subscribeRecords, _ := subscribeRecordModels.GetList(map[string][]string{"channel_id": {fmt.Sprintf("%v", query.WebsiteID)}, "increment_id": {query.IncrementID}})
	if len(subscribeRecords) > 0 {
		for _, record := range subscribeRecords {
			items, _ := trackItemModels.GetItems(record.ChannelID, map[string]interface{}{"order_id": record.OrderID, "track_number": record.TrackNumber})
			if len(items) > 0 {
				for _, item := range items {
					list = append(list, item)
				}
			}
		}
	}
	return
}
