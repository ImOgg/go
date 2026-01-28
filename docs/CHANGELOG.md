# 變更記錄 (Changelog)

所有重要的專案變更都會記錄在這個檔案中。

格式基於 [Keep a Changelog](https://keepachangelog.com/zh-TW/1.0.0/)。

---

## [Unreleased]

### 待辦事項

#### 核心功能
- [ ] JWT 認證（Auth）- 完整的登入/登出/Token 刷新
- [ ] Queue Job - 背景任務處理（類似 Laravel Queue）
- [ ] WebSocket - 即時通訊支援
- [ ] Email - 郵件發送（SMTP、第三方服務）
- [ ] File Storage - 檔案上傳（本地、S3、雲端）
- [ ] Cache - 快取策略（Redis 快取層封裝）
- [ ] Logging - 結構化日誌系統（類似 Laravel Log）

#### 安全性
- [ ] CORS 設定優化 - 目前是允許全部，需要限制來源
- [ ] XSS 防護 - 輸入過濾、輸出編碼
- [ ] CSRF 防護 - Token 驗證
- [ ] Rate Limiting - API 請求限流
- [ ] SQL Injection 防護 - 參數化查詢檢查
- [ ] Input Validation - 更完整的輸入驗證

#### 測試
- [ ] 單元測試（Unit Test）- Repository、Service 層
- [ ] 整合測試（Integration Test）- API 端對端
- [ ] Mock 測試 - 資料庫 Mock

#### 文件
- [ ] Swagger API 文件
- [ ] Postman Collection
- [ ] 部署文件（Docker、K8s）

#### 其他
- [ ] 多語系（i18n）
- [ ] 排程任務（Scheduler）- 類似 Laravel Task Scheduling
- [ ] Event/Listener - 事件驅動架構
- [ ] Notification - 通知系統（Email、SMS、Push）

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

**最後更新：** 2026-01-28
