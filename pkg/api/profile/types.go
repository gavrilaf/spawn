package profile

type ConfirmDeviceRequest struct {
	Code string `json:"code" form:"code" binding:"required"`
}
