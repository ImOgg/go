# Go 專案架構說明（Laravel 風格）

## 📁 新的專案結構

```
app/
├── app.go              # 應用程式容器（Service Container）
├── requests/           # 請求驗證（FormRequest）
│   └── user_request.go
├── services/           # 業務邏輯層（Service Layer）
│   └── user_service.go
├── repositories/       # 資料存取層（Repository Pattern）
│   └── user_repository.go
├── responses/          # 回應 DTO（Resource）
│   └── user_response.go
└── traits/             # 共用功能（Trait）
    ├── pagination.go
    └── response_helper.go
```

## 🏗️ 架構分層說明

### 1️⃣ **Controller（控制器）**
- 接收 HTTP 請求
- 驗證請求資料
- 呼叫 Service 層
- 回傳響應

```go
func (ctrl *UserController) Store(c *gin.Context) {
    var req requests.CreateUserRequest
    if err := req.Validate(c); err != nil {
        // 處理驗證錯誤
        return
    }
    user, err := ctrl.app.UserService.CreateUser(&req)
    traits.RespondCreated(c, user, "使用者建立成功")
}
```

### 2️⃣ **Request（請求驗證）**
類似 Laravel 的 `FormRequest`，負責：
- 定義驗證規則
- 驗證請求資料
- 格式化錯誤訊息

```go
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required,min=2,max=100"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"min=0,max=150"`
}
```

### 3️⃣ **Service（業務邏輯層）**
類似 Laravel 的 Service，負責：
- 處理業務邏輯
- 呼叫 Repository 層
- 資料轉換（Model → DTO）
- 業務規則驗證

```go
func (s *userService) CreateUser(req *requests.CreateUserRequest) (*responses.UserResponse, error) {
    // 檢查 Email 是否已存在
    existingUser, _ := s.userRepo.FindByEmail(req.Email)
    if existingUser != nil {
        return nil, errors.New("電子郵件已被使用")
    }
    
    // 建立使用者
    user := &models.User{Name: req.Name, Email: req.Email, Age: req.Age}
    s.userRepo.Create(user)
    
    return responses.NewUserResponse(user), nil
}
```

### 4️⃣ **Repository（資料存取層）**
類似 Laravel 的 Repository Pattern，負責：
- 資料庫操作（CRUD）
- 查詢邏輯封裝
- 與 ORM 互動

```go
type UserRepository interface {
    Create(user *models.User) error
    FindAll() ([]models.User, error)
    FindByID(id uint) (*models.User, error)
    Update(user *models.User) error
    Delete(id uint) error
}
```

### 5️⃣ **Response（回應 DTO）**
類似 Laravel 的 `Resource`，負責：
- 隱藏敏感資訊
- 格式化回應資料
- 統一 API 回應格式

```go
type UserResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    // 不包含敏感資料如密碼
}
```

### 6️⃣ **Trait（共用功能）**
類似 Laravel 的 Trait，提供：
- 分頁功能
- 統一的回應輔助函式
- 其他可重用的功能

```go
// 分頁
traits.Paginate(db, page, perPage, &users)

// 統一回應
traits.RespondSuccess(c, data, "成功")
traits.RespondError(c, 400, "失敗", errors)
```

### 7️⃣ **App Container（應用程式容器）**
類似 Laravel 的 Service Container，負責：
- 依賴注入
- 統一管理所有 Service 和 Repository
- 初始化應用程式

```go
app := app.NewApp(database.DB)
// app.UserService 可在任何 Controller 中使用
```

## 🚀 使用方式

### 測試新架構的 API

```bash
# 1. 啟動應用程式
go run main.go

# 2. 測試 API（新版本在 /api/v2）

# 建立使用者
curl -X POST http://localhost:8080/api/v2/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "張三",
    "email": "zhang@example.com",
    "age": 25
  }'

# 取得所有使用者
curl http://localhost:8080/api/v2/users

# 取得單一使用者
curl http://localhost:8080/api/v2/users/1

# 更新使用者
curl -X PUT http://localhost:8080/api/v2/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "張三豐", "age": 30}'

# 刪除使用者
curl -X DELETE http://localhost:8080/api/v2/users/1
```

## 🔄 資料流向

```
Client Request
    ↓
Controller (接收請求)
    ↓
Request (驗證資料)
    ↓
Service (業務邏輯)
    ↓
Repository (資料存取)
    ↓
Model (資料模型)
    ↓
Database
    ↓
Response DTO (格式化回應)
    ↓
Client Response
```

## 🆚 與 Laravel 的對照

| Laravel | Go (這個專案) | 說明 |
|---------|--------------|------|
| `FormRequest` | `app/requests/` | 請求驗證 |
| `$fillable` | Request Struct | 控制可寫入欄位 |
| `Service` | `app/services/` | 業務邏輯 |
| `Repository` | `app/repositories/` | 資料存取 |
| `Resource` | `app/responses/` | 回應 DTO |
| `Trait` | `app/traits/` | 共用功能 |
| `Container` | `app/app.go` | 依賴注入容器 |

## ✨ 優勢

1. **關注點分離**：每一層職責明確
2. **可測試性**：每一層都可獨立測試
3. **可維護性**：修改容易，不影響其他層
4. **可擴展性**：輕鬆新增功能
5. **依賴注入**：解耦合，易於替換實作

## 📝 新增其他功能範例

### 新增商品 (Product) 模組

1. 建立 Model: `models/product.go`
2. 建立 Request: `app/requests/product_request.go`
3. 建立 Repository: `app/repositories/product_repository.go`
4. 建立 Service: `app/services/product_service.go`
5. 建立 Response: `app/responses/product_response.go`
6. 建立 Controller: `controllers/product/product.go`
7. 在 `app/app.go` 註冊依賴
8. 在 `routes/v2_routes.go` 新增路由

完全模組化！
