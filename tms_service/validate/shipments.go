package validate

type LabelUrlQuery struct {
	WebsiteID   uint64 `form:"website_id"  binding:"required"  validate:"required,gt>0"`
	IncrementID string `form:"increment_id"  binding:"required"  validate:"required"`
}
