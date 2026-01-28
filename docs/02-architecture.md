# Go 專案架構說明（Laravel 風格）

## 專案結構總覽

```
my-api/
│
├── app/                          # 應用層（核心業務邏輯）
│   ├── app.go                   # 應用容器（Service Container）
│   ├── controllers/             # 控制層
│   │   └── user_controller.go
│   ├── services/                # 業務邏輯層
│   │   └── user_service.go
│   ├── repositories/            # 資料存取層
│   │   └── user_repository.go
│   ├── models/                  # 資料模型
│   │   └── user.go
│   ├── requests/                # 請求驗證（FormRequest）
│   │   └── user_request.go
│   ├── responses/               # 回應 DTO（Resource）
│   │   └── user_response.go
│   ├── middleware/              # 中間件
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   └── traits/                  # 共用功能（Trait）
│       ├── pagination.go
│       └── response_helper.go
│
├── bootstrap/                    # 啟動/連線初始化
│   ├── database.go              # GORM 資料庫連接
│   └── redis.go                 # Redis 連接
│
├── config/                       # 配置模組
│   └── config.go                # 環境變數載入
│
├── database/                     # 資料庫相關
│   ├── migrator.go              # Migration 執行引擎
│   └── migrations/              # Migration 檔案
│       ├── migration.go         # Migration 介面
│       ├── registry.go          # Migration 註冊器
│       └── 000001_create_users_table.go
│
├── routes/                       # 路由定義
│   └── api.go                   # API 路由（RESTful）
│
├── cmd/                          # 命令行工具
│   └── migrate/
│       └── main.go              # Migration CLI
│
├── docs/                         # 專案文件
│
├── main.go                       # 應用程式入口
├── .env                          # 環境變數
├── go.mod                        # Go 模組定義
├── Dockerfile                    # Docker 配置
└── docker-compose.yml
```

---

## 資料夾職責說明

| 資料夾 | 職責 | Laravel 對應 |
|--------|------|--------------|
| `app/controllers/` | 接收 HTTP 請求，呼叫 Service | `app/Http/Controllers` |
| `app/services/` | 業務邏輯處理 | `app/Services` |
| `app/repositories/` | 資料庫操作封裝 | `app/Repositories` |
| `app/models/` | 資料模型定義 | `app/Models` |
| `app/requests/` | 請求驗證規則 | `app/Http/Requests` |
| `app/responses/` | 回應格式化（DTO） | `app/Http/Resources` |
| `app/middleware/` | 中間件（認證、日誌等） | `app/Http/Middleware` |
| `app/traits/` | 共用功能 | `app/Traits` |
| `bootstrap/` | 啟動初始化（DB、Redis） | `bootstrap/` |
| `config/` | 配置載入 | `config/` |
| `database/` | Migration 相關 | `database/` |
| `routes/` | 路由定義 | `routes/` |

---

## 架構分層說明

### 1. Controller（控制器）

**位置：** `app/controllers/`

**職責：**
- 接收 HTTP 請求
- 驗證請求資料
- 呼叫 Service 層
- 回傳響應

```go
func (ctrl *UserController) Store(c *gin.Context) {
    var req requests.CreateUserRequest
    if err := req.Validate(c); err != nil {
        traits.RespondValidationError(c, err)
        return
    }
    user, err := ctrl.app.UserService.CreateUser(&req)
    traits.RespondCreated(c, user, "使用者建立成功")
}
```

### 2. Request（請求驗證）

**位置：** `app/requests/`

**職責：** 類似 Laravel 的 `FormRequest`
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

### 3. Service（業務邏輯層）

**位置：** `app/services/`

**職責：**
- 處理業務邏輯
- 呼叫 Repository 層
- 資料轉換（Model → DTO）
- 業務規則驗證

```go
func (s *userService) CreateUser(req *requests.CreateUserRequest) (*responses.UserResponse, error) {
    // 業務邏輯：檢查 Email 是否已存在
    existingUser, _ := s.userRepo.FindByEmail(req.Email)
    if existingUser != nil {
        return nil, errors.New("電子郵件已被使用")
    }

    user := &models.User{Name: req.Name, Email: req.Email, Age: req.Age}
    s.userRepo.Create(user)

    return responses.NewUserResponse(user), nil
}
```

### 4. Repository（資料存取層）

**位置：** `app/repositories/`

**職責：**
- 資料庫操作（CRUD）
- 查詢邏輯封裝
- 與 ORM 互動

```go
type UserRepository interface {
    Create(user *models.User) error
    FindAll() ([]models.User, error)
    FindByID(id uint) (*models.User, error)
    FindByEmail(email string) (*models.User, error)
    Update(user *models.User) error
    Delete(id uint) error
}
```

### 5. Response（回應 DTO）

**位置：** `app/responses/`

**職責：** 類似 Laravel 的 `Resource`
- 隱藏敏感資訊（如密碼）
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

### 6. Bootstrap（啟動初始化）

**位置：** `bootstrap/`

**職責：** 類似 Laravel 的 `bootstrap/` 和 `ServiceProvider`
- 資料庫連線初始化
- Redis 連線初始化
- 其他服務啟動

```go
// bootstrap/database.go
func InitDB() {
    // 根據 .env 配置連接 MySQL 或 PostgreSQL
    DB, err = gorm.Open(dialector, &gorm.Config{})
}

// bootstrap/redis.go
func InitRedis() {
    RedisClient = redis.NewClient(&redis.Options{...})
}
```

### 7. App Container（應用程式容器）

**位置：** `app/app.go`

**職責：** 類似 Laravel 的 Service Container
- 依賴注入
- 統一管理所有 Service 和 Repository
- 初始化應用程式

```go
type App struct {
    DB             *gorm.DB
    UserRepository repositories.UserRepository
    UserService    services.UserService
}

func NewApp(db *gorm.DB) *App {
    userRepo := repositories.NewUserRepository(db)
    userService := services.NewUserService(userRepo)

    return &App{
        DB:             db,
        UserRepository: userRepo,
        UserService:    userService,
    }
}
```

---

## 資料流向

```
HTTP Request
     │
     ▼
┌─────────────────┐
│     Routes      │  ← 路由匹配
│   (routes/)     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Middleware    │  ← 認證、CORS、日誌
│  (middleware/)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Controller    │  ← 接收請求
│ (controllers/)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Request      │  ← 驗證資料
│  (requests/)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Service      │  ← 業務邏輯
│  (services/)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Repository    │  ← 資料存取
│ (repositories/) │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│     Model       │  ← 資料結構
│   (models/)     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Database     │  ← MySQL/PostgreSQL
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Response     │  ← DTO 轉換
│  (responses/)   │
└────────┬────────┘
         │
         ▼
HTTP Response (JSON)
```

---

## 啟動流程

```go
// main.go
func main() {
    // 1. 載入配置（.env）
    config.LoadConfig()

    // 2. 初始化資料庫連接
    bootstrap.InitDB()

    // 3. 自動執行 migrations
    database.RunMigrations()

    // 4. 初始化 Redis（可選）
    // bootstrap.InitRedis()

    // 5. 建立應用程式容器
    application := app.NewApp(bootstrap.DB)

    // 6. 設置路由
    r := gin.Default()
    routes.SetupRoutes(r, application)

    // 7. 啟動服務
    r.Run(":8080")
}
```

---

## 與 Laravel 的完整對照

| Laravel | Go (本專案) | 說明 |
|---------|-------------|------|
| `app/Http/Controllers` | `app/controllers/` | 控制器 |
| `app/Http/Requests` | `app/requests/` | 請求驗證 |
| `app/Services` | `app/services/` | 業務邏輯 |
| `app/Repositories` | `app/repositories/` | 資料存取 |
| `app/Http/Resources` | `app/responses/` | 回應 DTO |
| `app/Models` | `app/models/` | 資料模型 |
| `app/Http/Middleware` | `app/middleware/` | 中間件 |
| `app/Traits` | `app/traits/` | 共用功能 |
| `bootstrap/app.php` | `bootstrap/` | 啟動初始化 |
| `app/Providers` | `app/app.go` | 依賴注入容器 |
| `config/*.php` | `config/config.go` | 配置 |
| `routes/api.php` | `routes/api.go` | 路由 |
| `database/migrations` | `database/migrations/` | Migration |
| `php artisan` | `cmd/migrate/` | CLI 工具 |

---

## 新增模組範例

### 新增 Product（商品）模組

1. **建立 Model：** `app/models/product.go`
2. **建立 Request：** `app/requests/product_request.go`
3. **建立 Repository：** `app/repositories/product_repository.go`
4. **建立 Service：** `app/services/product_service.go`
5. **建立 Response：** `app/responses/product_response.go`
6. **建立 Controller：** `app/controllers/product_controller.go`
7. **在 `app/app.go` 註冊依賴**
8. **在 `routes/api.go` 新增路由**
9. **建立 Migration：** `database/migrations/000002_create_products_table.go`

---

## 架構優勢

| 優勢 | 說明 |
|------|------|
| **關注點分離** | 每一層職責明確，易於理解 |
| **可測試性** | 每一層都可獨立測試（Mock Repository） |
| **可維護性** | 修改某層不影響其他層 |
| **可擴展性** | 新功能只需新增對應層的檔案 |
| **依賴注入** | 解耦合，易於替換實作 |
| **熟悉度** | Laravel 開發者可以快速上手 |

---

## 相關文件

- [專案設置](01-setup.md)
- [資料庫連接](03-database.md)
- [Migration 系統](04-migration.md)
- [Controller 結構](05-controller-structure.md)

---

**最後更新：** 2026-01-28
