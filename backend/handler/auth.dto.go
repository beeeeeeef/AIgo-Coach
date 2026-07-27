package handler

type SendRegisterCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginByPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
type LoginByCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type LoginResponse struct{

}

type MeResponse struct {
	
}