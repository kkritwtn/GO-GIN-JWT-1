package dto

type LoginAdminRequest struct {
	Username string `json:"username" binding:"required,min=4"`
	Password string `json:"password" binding:"required,min=6"`
}
