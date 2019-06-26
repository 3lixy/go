package validate

type AddShipmentTrackQuery struct {
	WebsiteID   uint64  `json:"website_id" form:"website_id" binding:"required"`
	IncrementID string  `json:"increment_id" form:"increment_id" binding:"required"`
	CarrierCode string  `json:"carrier_code" form:"carrier_code" binding:"required"`
	TrackNumber string  `json:"track_number" form:"track_number" binding:"required"`
	Length      float64 `json:"length" form:"length"`
	Width       float64 `json:"width" form:"width"`
	Height      float64 `json:"height" form:"height"`
	Weight      float64 `json:"weight" form:"weight"`
}

type DeleteShipmentTrackQuery struct {
	EntityID string `json:"entity_id" form:"entity_id" binding:"required"`
}
