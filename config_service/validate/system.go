package validate

type UpdateSystemQuery struct {
	WebsiteID	uint64	`form:"website_id"  binding:"required"  validate:"required,gt>0"`
	SystemID	uint64	`form:"system_id"  binding:"required"  validate:"required"`
	Status		uint64	`form:"status"  binding:"required"  validate:"required,gt>0"`
	Message		string	`form:"message"`
}