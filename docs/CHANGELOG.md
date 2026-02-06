# 變更記錄 (Changelog)

所有重要的專案變更都會記錄在這個檔案中。

格式基於 [Keep a Changelog](https://keepachangelog.com/zh-TW/1.0.0/)。

---

## [Unreleased]

### 待辦事項

#### 核心功能
- [x] JWT 認證（Auth）- 完整的登入/登出
- [ ] Refresh Token - 自動續期機制
- [ ] 角色權限（RBAC）- 使用者角色與權限管理
- [ ] Token 黑名單 - 用 Redis 實作登出失效
- [ ] 社交登入 - Google、GitHub OAuth
- [ ] Queue Job - 背景任務處理（類似 Laravel Queue）
- [ ] WebSocket - 即時通訊支援
- [ ] Email - 郵件發送（SMTP、第三方服務）
- [ ] File Storage - 檔案上傳（本地、S3、雲端）
- [ ] Cache - 快取策略（Redis 快取層封裝）
- [x] Logging - 結構化日誌系統（類似 Laravel Log）

#### 安全性
- [ ] CORS 設定優化 - 目前是允許全部，需要限制來源
- [ ] XSS 防護 - 輸入過濾、輸出編碼
- [ ] CSRF 防護 - Token 驗證
- [ ] Rate Limiting - API 請求限流
- [ ] SQL Injection 防護 - 參數化查詢檢查
- [x] Input Validation - 更完整的輸入驗證

#### 測試
- [x] 單元測試範例 - `app/utils/jwt_test.go`
- [x] Mock 測試 - `app/services/user_service_test.go`
- [ ] 單元測試（Unit Test）- Repository 層
- [ ] 整合測試（Integration Test）- API 端對端

#### 文件
- [x] 參數驗證指南 - `docs/12-validation.md`
- [ ] Swagger API 文件
- [ ] Postman Collection
- [ ] 部署文件（Docker、K8s）

#### 其他
- [ ] 多語系（i18n）
- [ ] 排程任務（Scheduler）- 類似 Laravel Task Scheduling
- [ ] Event/Listener - 事件驅動架構
- [ ] Notification - 通知系統（Email、SMS、Push）

---

## [0.6.0] - 2026-02-03

### 新增 - 結構化日誌系統

#### 新增
- `app/pkg/logger/logger.go` - Logger 核心封裝
  - 使用 zerolog 高性能日誌套件
  - 支援多日誌等級：debug, info, warn, error, fatal
  - Laravel 風格的便捷方法：`Debug()`, `Info()`, `Warning()`, `Error()`, `Fatal()`
  - 支援結構化 context 附加資訊
  - `WithRequestID()` - 建立帶有 Request ID 的子 Logger
  - `WithContext()` - 建立帶有自訂欄位的子 Logger
  - `WithError()` - 建立帶有錯誤資訊的子 Logger

- `app/pkg/logger/context.go` - Context 輔助函數
  - `FromGinContext()` - 從 Gin context 取得 Logger
  - `ToGinContext()` - 將 Logger 存入 Gin context
  - `GetRequestID()` - 取得 Request ID

- `bootstrap/logger.go` - Logger 初始化
  - 全域 `bootstrap.Log` 變數

- `app/middleware/request_id.go` - Request ID 中間件
  - 為每個請求生成唯一 UUID
  - 支援從 `X-Request-ID` header 取得（分散式追蹤）
  - 自動設定 Response Header

- `storage/logs/` - 日誌存放目錄

#### 變更
- `config/config.go` - 新增 LogConfig 配置結構
- `app/app.go` - 注入 Logger 到 App 容器
- `app/middleware/logger.go` - 改用 zerolog 結構化日誌
- `routes/api.go` - 新增全域中間件（Recovery、RequestID、Logger、CORS）
- `main.go` - 新增 Logger 初始化，改用 `gin.New()`

#### 環境變數
```env
LOG_LEVEL=debug        # 日誌等級
LOG_FORMAT=console     # console 或 json
LOG_OUTPUT=stdout      # stdout, file, both
LOG_FILE_PATH=storage/logs/app.log
LOG_MAX_SIZE=100       # 單檔最大 MB
LOG_MAX_BACKUPS=30     # 保留檔案數
LOG_MAX_AGE=30         # 保留天數
LOG_COMPRESS=true      # 壓縮舊檔
```

#### 依賴套件
- `github.com/rs/zerolog` - 高性能結構化日誌
- `gopkg.in/natefinch/lumberjack.v2` - 日誌檔案輪替
- `github.com/google/uuid` - UUID 生成

#### 日誌輸出範例

開發環境 (console):
```
2026-02-03 14:30:45 INF Request request_id=abc-123 status=200 method=GET path=/api/users latency=12ms
```

生產環境 (JSON):
```json
{"level":"info","time":"2026-02-03T14:30:50+08:00","request_id":"abc-123","status":200,"method":"GET","path":"/api/users","message":"Request"}
```

---

## [0.5.0] - 2026-02-02

### 新增 - 測試框架與文件

#### 新增
- `app/utils/jwt_test.go` - 單元測試範例
  - `TestHashPassword` - 測試密碼加密
  - `TestCheckPassword` - 測試密碼驗證
  - `TestHashPassword_DifferentHashes` - 測試 bcrypt salt 機制
  - `TestCheckPassword_TableDriven` - Table-Driven 測試範例
  - `TestHashPassword_EmptyPassword` - 邊界條件測試

- `app/services/user_service_test.go` - Mock 測試範例
  - `mockUserRepository` - 模擬 UserRepository interface
  - `TestUserService_CreateUser` - 測試新增使用者（含 Email 重複檢查）
  - `TestUserService_GetUserByID` - 測試取得單一使用者
  - `TestUserService_UpdateUser` - 測試更新使用者
  - `TestUserService_DeleteUser` - 測試刪除使用者
  - `TestUserService_GetAllUsers` - 測試取得所有使用者

- `docs/09-testing.md` - 測試指南文件
  - Go 測試基礎
  - 與 Laravel PHPUnit 對比
  - Table-Driven 測試模式
  - Mock 測試原理與範例
  - 測試環境設定說明
  - 常用測試指令

#### 移除
- `models/` 根目錄資料夾 - 刪除空的舊資料夾（已遷移至 `app/models/`）

#### 文件更新
- `docs/00-index.md` - 新增測試指南索引

#### 測試指令
```bash
# 執行測試
docker exec my-go-app go test -v ./app/utils/...

# 執行整個專案測試
docker exec my-go-app go test -v ./...
```

---

## [0.4.0] - 2026-01-29

### 新增 - JWT 認證系統

#### 新增
- `app/utils/jwt.go` - JWT 工具函式
  - `GenerateToken()` - 產生 JWT Token
  - `ValidateToken()` - 驗證 JWT Token
  - `HashPassword()` - 密碼加密（bcrypt）
  - `CheckPassword()` - 密碼驗證

- `app/controllers/auth_controller.go` - 認證控制器
  - `Register()` - 使用者註冊
  - `Login()` - 使用者登入
  - `Logout()` - 使用者登出
  - `Me()` - 取得當前用戶資訊

- `app/services/auth_service.go` - 認證服務層

- `app/requests/auth_request.go` - 認證請求驗證
  - `RegisterRequest` - 註冊請求驗證
  - `LoginRequest` - 登入請求驗證

- `app/responses/auth_response.go` - 認證回應 DTO

- `database/migrations/000002_add_password_to_users.go` - 新增 password 欄位

#### 變更
- `app/models/user.go` - 新增 Password 欄位（json:"-" 隱藏）
- `app/middleware/auth.go` - 實作真正的 JWT 驗證
- `app/app.go` - 註冊 AuthService
- `config/config.go` - 新增 JWTConfig 設定
- `routes/api.go` - 新增認證路由
- `.env` - 新增 JWT_SECRET、JWT_EXPIRY_HOURS

#### API 端點
| 方法 | 路徑 | 說明 | 需驗證 |
|------|------|------|--------|
| POST | `/api/register` | 使用者註冊 | 否 |
| POST | `/api/login` | 使用者登入 | 否 |
| POST | `/api/logout` | 使用者登出 | 是 |
| GET | `/api/me` | 取得當前用戶 | 是 |

#### 依賴套件
- `github.com/golang-jwt/jwt/v5` - JWT 處理
- `golang.org/x/crypto/bcrypt` - 密碼加密

---

## [0.3.0] - 2026-01-28

### 重構 - 資料夾結構優化

#### 新增
- `bootstrap/` 資料夾 - 統一管理啟動/連線初始化
  - `bootstrap/database.go` - GORM 資料庫連接
  - `bootstrap/redis.go` - Redis 連接

#### 變更
- `database/migrate_simple.go` → `database/migrator.go`（重新命名）
- `main.go` - 更新 import 路徑，使用 `bootstrap.InitDB()` 和 `bootstrap.DB`

#### 移除
- `database/db.go` - 移至 `bootstrap/database.go`
- `database/redis.go` - 移至 `bootstrap/redis.go`
- `database/sql_raw.go` - 刪除（未使用的備用方案）
- `controllers/` 根目錄 - 刪除舊的學習用測試代碼
  - `controllers/hello/hello.go`
  - `controllers/test/test.go`

#### 文件更新
- 重新編號所有 docs 文件（01-setup.md, 02-architecture.md, ...）
- 新增 `docs/00-index.md` 目錄索引
- 更新 `docs/02-architecture.md` 反映新的資料夾結構
- 新增 `docs/CHANGELOG.md` 變更記錄

---

## [0.2.0] - 2026-01-27

### 新增 - Laravel 風格分層架構

#### 新增
- `app/` 應用層完整結構
  - `app/app.go` - 應用容器（依賴注入）
  - `app/controllers/user_controller.go` - 使用者控制器
  - `app/services/user_service.go` - 使用者業務邏輯
  - `app/repositories/user_repository.go` - 使用者資料存取（含 Interface）
  - `app/models/user.go` - 使用者模型
  - `app/requests/user_request.go` - 請求驗證
  - `app/responses/user_response.go` - 回應 DTO
  - `app/middleware/` - 中間件（auth, cors, logger）
  - `app/traits/` - 共用功能（pagination, response_helper）

- `database/migrations/` - 自訂 Migration 系統
  - `migration.go` - Migration 介面定義
  - `registry.go` - Migration 註冊器
  - `000001_create_users_table.go` - 使用者表 Migration

- `cmd/migrate/main.go` - Migration CLI 工具
  - `migrate` - 執行所有待執行的 migrations
  - `rollback` - 回滾最後一個 migration
  - `status` - 查看 migration 狀態
  - `make <name>` - 建立新的 migration 檔案

- `routes/api.go` - RESTful API 路由

#### 變更
- 重構專案結構為 Laravel 風格分層架構
- `main.go` - 整合應用容器和路由設定

---

## [0.1.0] - 2026-01-23

### 初始版本

#### 新增
- 專案初始化
- Gin Web 框架整合
- GORM ORM 整合
- MySQL / PostgreSQL 支援
- Redis 連接支援
- Docker 配置
- 環境變數配置（.env）
- 基本的 Hello World API

#### 文件
- `docs/setup.md` - 專案設置指南
- `docs/database.md` - 資料庫連接說明
- `docs/commands.md` - 常用命令
- `docs/troubleshooting.md` - 常見問題

---

## 版本說明

- **主版本號 (Major)**：不相容的 API 變更
- **次版本號 (Minor)**：新增功能（向下相容）
- **修訂號 (Patch)**：Bug 修復（向下相容）

---

## 貢獻者

- 專案作者

---

**最後更新：** 2026-02-03
