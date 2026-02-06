# 參數驗證設定指南

> 本文檔說明如何驗證前端傳來的參數，確保數據安全性和完整性

---

## 📌 概述

當前端發送請求時，後端需要驗證這些參數：
- **資料格式** - Email、電話號碼、URL 等格式是否正確
- **必填檢查** - 必須提供的欄位是否存在
- **長度限制** - 字串長度、數值範圍是否符合
- **業務規則** - Email 是否已被使用、年齡是否合法等

---

## 🛠️ 常用工具

### 1. validator 套件

安裝依賴：
```bash
go get github.com/go-playground/validator/v10
```

### 2. 專案中的驗證方式

查看現有範例：
- `app/requests/user_request.go` - 使用者相關驗證
- `app/requests/auth_request.go` - 認證相關驗證

---

## 📝 設定驗證規則

### 步驟 1: 定義驗證結構體

在 `app/requests/` 下創建驗證文件：

```go
// app/requests/user_request.go

package requests

type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"required,min=18,max=120"`
    Phone    string `json:"phone" validate:"required,e164"`
}

type UpdateUserRequest struct {
    Name  string `json:"name" validate:"omitempty,min=2,max=50"`
    Email string `json:"email" validate:"omitempty,email"`
    Age   int    `json:"age" validate:"omitempty,min=18,max=120"`
}
```

### 步驟 2: 創建驗證幫助函數

在 `app/requests/validator.go` 中：

```go
package requests

import (
    "fmt"
    "github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
}

// ValidateStruct 驗證結構體並返回錯誤訊息
func ValidateStruct(data interface{}) error {
    err := validate.Struct(data)
    if err == nil {
        return nil
    }

    validationErrors := err.(validator.ValidationErrors)
    var messages []string

    for _, fieldError := range validationErrors {
        messages = append(messages, formatValidationError(fieldError))
    }

    return fmt.Errorf("validation failed: %v", messages)
}

// formatValidationError 格式化驗證錯誤訊息
func formatValidationError(fe validator.FieldError) string {
    field := fe.Field()
    tag := fe.Tag()

    switch tag {
    case "required":
        return fmt.Sprintf("%s 是必填欄位", field)
    case "min":
        return fmt.Sprintf("%s 最少需要 %s 個字符", field, fe.Param())
    case "max":
        return fmt.Sprintf("%s 最多只能 %s 個字符", field, fe.Param())
    case "email":
        return fmt.Sprintf("%s 格式不正確", field)
    case "e164":
        return fmt.Sprintf("%s 電話格式不正確", field)
    default:
        return fmt.Sprintf("%s 驗證失敗: %s", field, tag)
    }
}
```

---

## 🔍 常見驗證規則

| 規則 | 說明 | 例子 |
|------|------|------|
| `required` | 必填 | `validate:"required"` |
| `email` | Email 格式 | `validate:"email"` |
| `min=n` | 最小值/長度 | `validate:"min=2"` |
| `max=n` | 最大值/長度 | `validate:"max=50"` |
| `e164` | 國際電話格式 | `validate:"e164"` |
| `url` | URL 格式 | `validate:"url"` |
| `numeric` | 數字 | `validate:"numeric"` |
| `alpha` | 純字母 | `validate:"alpha"` |
| `alphanumeric` | 字母和數字 | `validate:"alphanumeric"` |
| `omitempty` | 可選（非空時驗證） | `validate:"omitempty,email"` |
| `gt=n` | 大於 | `validate:"gt=0"` |
| `gte=n` | 大於等於 | `validate:"gte=0"` |
| `lt=n` | 小於 | `validate:"lt=100"` |
| `lte=n` | 小於等於 | `validate:"lte=100"` |

---

## 💡 在 Controller 中使用驗證

### 基本用法

```go
// app/controllers/user_controller.go

package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "app/requests"
    "app/services"
)

type UserController struct {
    userService *services.UserService
}

// Store 創建使用者
func (uc *UserController) Store(c *gin.Context) {
    var req requests.CreateUserRequest

    // 綁定 JSON 數據
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "無效的請求格式",
            "error": err.Error(),
        })
        return
    }

    // 驗證數據
    if err := requests.ValidateStruct(req); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "message": "驗證失敗",
            "error": err.Error(),
        })
        return
    }

    // 調用 Service 創建使用者
    user, err := uc.userService.CreateUser(req.Name, req.Email, req.Age)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "message": "創建失敗",
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "創建成功",
        "data": user,
    })
}
```

### 帶有中間件的驗證

```go
// 建立驗證中間件
func ValidateRequest(reqType interface{}) gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := c.ShouldBindJSON(&reqType); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "message": "無效的請求格式",
                "error": err.Error(),
            })
            c.Abort()
            return
        }

        if err := requests.ValidateStruct(reqType); err != nil {
            c.JSON(http.StatusUnprocessableEntity, gin.H{
                "message": "驗證失敗",
                "error": err.Error(),
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

---

## ✅ 業務規則驗證

除了格式驗證，還需要驗證業務規則：

```go
// app/requests/validator.go - 自訂驗證函數

package requests

import (
    "github.com/go-playground/validator/v10"
    "app/repositories"
)

func RegisterCustomValidators(validate *validator.Validate, userRepo repositories.UserRepository) {
    // 驗證 Email 是否已存在
    validate.RegisterValidationFunc("email_unique", func(fl validator.FieldLevel) bool {
        email := fl.Field().String()
        exists, _ := userRepo.EmailExists(email)
        return !exists
    })

    // 驗證年齡是否為成人
    validate.RegisterValidationFunc("adult", func(fl validator.FieldLevel) bool {
        age := int(fl.Field().Int())
        return age >= 18
    })
}
```

使用自訂驗證：

```go
type RegisterRequest struct {
    Email string `json:"email" validate:"required,email,email_unique"`
    Age   int    `json:"age" validate:"required,adult"`
}
```

---

## 📋 驗證清單

在實作參數驗證時，檢查以下項目：

### 設定階段
- [ ] 已在 `app/requests/` 下定義驗證結構體
- [ ] 各欄位有適當的驗證標籤
- [ ] 中文欄位名稱對應正確
- [ ] 必填和可選欄位標記清楚

### 實作階段
- [ ] Controller 中有調用 `ValidateStruct()`
- [ ] 驗證失敗時返回 422（UnprocessableEntity）狀態碼
- [ ] 錯誤訊息清晰易懂（中文）
- [ ] 敏感欄位（如密碼）不在錯誤中洩露

### 測試階段
- [ ] 測試必填欄位驗證
- [ ] 測試格式驗證（Email、Phone）
- [ ] 測試長度限制
- [ ] 測試邊界值
- [ ] 測試業務規則驗證

### 文檔階段
- [ ] API 文檔列出所有驗證規則
- [ ] 提供請求範例
- [ ] 說明可能的錯誤回應

---

## 🧪 測試驗證規則

### 單元測試例子

```go
// app/requests/user_request_test.go

package requests

import (
    "testing"
)

func TestCreateUserRequest_Validation(t *testing.T) {
    tests := []struct {
        name    string
        req     CreateUserRequest
        wantErr bool
    }{
        {
            name:    "Valid request",
            req:     CreateUserRequest{Name: "John", Email: "john@example.com", Age: 25},
            wantErr: false,
        },
        {
            name:    "Missing name",
            req:     CreateUserRequest{Email: "john@example.com", Age: 25},
            wantErr: true,
        },
        {
            name:    "Invalid email",
            req:     CreateUserRequest{Name: "John", Email: "invalid", Age: 25},
            wantErr: true,
        },
        {
            name:    "Age too young",
            req:     CreateUserRequest{Name: "John", Email: "john@example.com", Age: 15},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateStruct(tt.req)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateStruct() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## 🚀 最佳實踐

### 1. 分層驗證

```
前端驗證（UX）→ 請求驗證（格式）→ 業務驗證（規則）→ 資料庫唯一性
```

### 2. 統一錯誤回應格式

```go
{
    "message": "驗證失敗",
    "errors": {
        "email": ["Email 格式不正確", "Email 已被使用"],
        "age": ["年齡必須 >= 18"]
    }
}
```

### 3. 記錄驗證失敗

```go
if err := requests.ValidateStruct(req); err != nil {
    logger.Warning("Validation failed",
        "user_ip": c.ClientIP(),
        "endpoint": c.FullPath(),
        "error": err.Error(),
    )
}
```

### 4. 避免資訊洩露

```go
// ❌ 錯誤：暴露內部實現
"error": "User with email john@example.com already exists"

// ✅ 正確：通用訊息
"error": "Email 已被使用，請使用其他 Email"
```

---

## 📚 相關資源

- [validator 文檔](https://github.com/go-playground/validator)
- [現有驗證例子](../app/requests/)
- [API 錯誤處理](05-controller-structure.md)

---

**最後更新**: 2026-02-06
