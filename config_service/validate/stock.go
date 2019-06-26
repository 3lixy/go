package validate

type AddQuery struct {
	StockType   uint64 `form:"stock_type"  binding:"required"  validate:"required,gt>0"`
	StockName   string `form:"stock_name"  binding:"required"  validate:"required"`
	CompanyName string `form:"company_name"  binding:"required"  validate:"required"`
	Status      uint64 `form:"status"  binding:"required"  validate:"required,gt>0"`
	Country     string `form:"country"  binding:"required"  validate:"required"`
	Province    string `form:"province"`
	City        string `form:"city"`
	County      string `form:"county"`
	AddressOne  string `form:"address_one"  binding:"required"  validate:"required"`
	AddressTwo  string `form:"address_two"`
	Postcode    string `form:"postcode"`
	LastName    string `form:"last_name"  binding:"required"  validate:"required"`
	FirstName   string `form:"first_name"  binding:"required"  validate:"required"`
	Position    string `form:"position"`
	Telephone   string `form:"telephone"  binding:"required"  validate:"required"`
	Email       string `form:"email"  binding:"required"  validate:"required"`
	Wechat      string `form:"wechat"`
}

type DetailQuery struct {
	EntityID uint64 `form:"entity_id"  binding:"required"  validate:"required,gt>0"`
}

type UpdateQuery struct {
	EntityID    uint64 `form:"entity_id"  binding:"required"  validate:"required,gt>0"`
	StockType   uint64 `form:"stock_type"  binding:"required"  validate:"required,gt>0"`
	StockName   string `form:"stock_name"  binding:"required"  validate:"required"`
	CompanyName string `form:"company_name"  binding:"required"  validate:"required"`
	Status      uint64 `form:"status"  binding:"required"  validate:"required,gt>0"`
	Country     string `form:"country"  binding:"required"  validate:"required"`
	Province    string `form:"province"`
	City        string `form:"city"`
	County      string `form:"county"`
	AddressOne  string `form:"address_one"  binding:"required"  validate:"required"`
	AddressTwo  string `form:"address_two"`
	Postcode    string `form:"postcode"`
	LastName    string `form:"last_name"  binding:"required"  validate:"required"`
	FirstName   string `form:"first_name"  binding:"required"  validate:"required"`
	Position    string `form:"position"`
	Telephone   string `form:"telephone"  binding:"required"  validate:"required"`
	Email       string `form:"email"  binding:"required"  validate:"required"`
	Wechat      string `form:"wechat"`
}
