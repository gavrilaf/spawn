package profile

type ConfirmDeviceRequest struct {
	Code string `json:"code" form:"code" binding:"required"`
}

type UpdateCountryRequest struct {
	Country string `json:"country" form:"country" binding:"required"`
}

type UpdatePersonalInfoRequest struct {
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`
	BirthDate string `json:"birth_date" form:"birth_date" binding:"required"`
}
