package dto

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserRegisterResponse struct {
	ID uint64 `json:"id" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	ID       uint64 `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type ResetPasswordRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
