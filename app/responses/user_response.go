package responses

import (
	"time"
	"my-api/app/models"
)

// UserResponse - 使用者回應 DTO（隱藏敏感資訊）
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUserResponse - 將 Model 轉換為 Response DTO
func NewUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ApiResponse - 統一的 API 回應格式
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// NewSuccessResponse - 成功回應
func NewSuccessResponse(data interface{}, message string) ApiResponse {
	return ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse - 錯誤回應
func NewErrorResponse(errors interface{}, message string) ApiResponse {
	return ApiResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}
