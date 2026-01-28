package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CreateUserRequest - 建立使用者的請求驗證
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"min=0,max=150"`
}

// UpdateUserRequest - 更新使用者的請求驗證
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=2,max=100"`
	Email string `json:"email" binding:"omitempty,email"`
	Age   int    `json:"age" binding:"omitempty,min=0,max=150"`
}

// Validate - 驗證請求資料
func (r *CreateUserRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return err
	}
	return nil
}

// Validate - 驗證更新請求資料
func (r *UpdateUserRequest) Validate(c *gin.Context) error {
	if err := c.ShouldBindJSON(r); err != nil {
		return err
	}
	return nil
}

// FormatValidationError - 格式化驗證錯誤訊息（類似 Laravel 的錯誤格式）
func FormatValidationError(err error) map[string]interface{} {
	errors := make(map[string]interface{})
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			switch e.Tag() {
			case "required":
				errors[field] = field + " 欄位為必填"
			case "email":
				errors[field] = field + " 必須是有效的電子郵件"
			case "min":
				errors[field] = field + " 長度不得小於 " + e.Param()
			case "max":
				errors[field] = field + " 長度不得大於 " + e.Param()
			default:
				errors[field] = field + " 驗證失敗"
			}
		}
	}
	
	return errors
}
