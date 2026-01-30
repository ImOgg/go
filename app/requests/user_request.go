package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CreateUserRequest - 建立使用者的請求驗證
// 類似 Laravel 的 FormRequest，binding tag 等於 rules()
//
// binding tag 驗證規則說明：
// - required：必填欄位
// - min=N：最小長度/數值
// - max=N：最大長度/數值
// - email：必須是有效的 email 格式
// - omitempty：欄位可選，有值才驗證（用於 Update）
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"` // 必填，長度 2~100
	Email string `json:"email" binding:"required,email"`        // 必填，必須是 email 格式
	Age   int    `json:"age" binding:"min=0,max=150"`           // 選填，範圍 0~150
}

// UpdateUserRequest - 更新使用者的請求驗證
// omitempty = 欄位可選，只有在有值時才進行驗證（適合 PATCH/PUT 更新）
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=2,max=100"` // 選填，有值時長度需 2~100
	Email string `json:"email" binding:"omitempty,email"`        // 選填，有值時需是 email 格式
	Age   int    `json:"age" binding:"omitempty,min=0,max=150"`  // 選填，有值時範圍 0~150
}

// Validate - 驗證請求資料
// ShouldBindJSON 會自動：
// 1. 解析 JSON 請求 body
// 2. 根據 binding tag 驗證資料
// 類似 Laravel 的 $request->validate()
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

// FormatValidationError - 格式化驗證錯誤訊息
// 類似 Laravel 的 messages() 方法，將錯誤碼轉成友善的中文訊息
//
// 型別斷言說明：err.(validator.ValidationErrors)
// - 檢查 err 是否為 ValidationErrors 類型
// - ok = true 表示轉型成功，可以遍歷每個欄位的錯誤
//
// e.Tag() 回傳驗證規則名稱，如 "required"、"email"、"min"
// e.Param() 回傳規則參數，如 min=2 中的 "2"
func FormatValidationError(err error) map[string]interface{} {
	errors := make(map[string]interface{})

	// 型別斷言：將 error 轉換為 validator.ValidationErrors
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
