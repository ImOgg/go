package requests

import "github.com/gin-gonic/gin"

// RegisterRequest - 使用者註冊請求驗證
type RegisterRequest struct {
	Name            string `json:"name" binding:"required,min=2,max=100"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8,max=72"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
	Age             int    `json:"age" binding:"omitempty,min=0,max=150"`
}

// Validate - 驗證註冊請求
func (r *RegisterRequest) Validate(c *gin.Context) error {
	return c.ShouldBindJSON(r)
}

// LoginRequest - 使用者登入請求驗證
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Validate - 驗證登入請求
func (r *LoginRequest) Validate(c *gin.Context) error {
	return c.ShouldBindJSON(r)
}
