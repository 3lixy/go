package validate

type StoreDetailQuery struct {
	WebsiteID uint64 `form:"website_id"  binding:"required"  validate:"required,gt>0"`
	StoreID   uint64 `form:"store_id"  binding:"required"  validate:"required,gt>0"`
}
