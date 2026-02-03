# Go API 專案文件索引

> 這是一個採用 **Laravel 風格分層架構**的 Go Web API 專案

---

## 快速導覽

| 編號 | 文件 | 說明 | 適合誰看 |
|------|------|------|----------|
| 01 | [專案設置](01-setup.md) | 初始化專案、安裝依賴 | 新手入門 |
| 02 | [架構說明](02-architecture.md) | 專案結構、分層設計、資料流 | 想了解全貌 |
| 03 | [資料庫連接](03-database.md) | GORM、MySQL、環境配置 | 資料庫操作 |
| 04 | [Migration 系統](04-migration.md) | 資料表版本控制、建立/回滾 | 資料庫變更 |
| 05 | [Controller 結構](05-controller-structure.md) | Controller 組織方式、與 Laravel 對比 | 寫 API |
| 06 | [常用命令](06-commands.md) | go build、go run、go mod 等 | 快速參考 |
| 07 | [常見問題](07-troubleshooting.md) | 錯誤排查、解決方案 | 遇到問題 |
| 08 | [JWT 認證](08-jwt-authentication.md) | Token 驗證、登入/註冊 | 身份驗證 |
| 09 | [測試指南](09-testing.md) | 單元測試、Table-Driven 測試 | 品質保證 |
| 10 | [併發處理](10-concurrency.md) | Goroutines、Channels、Worker Pool | 任務隊列 |
| 11 | [日誌系統](11-logging.md) | zerolog 結構化日誌、Request ID | 除錯追蹤 |

---

## 建議閱讀順序

### 第一次接觸這個專案？

```
01-setup.md → 02-architecture.md → 03-database.md
```

### 想新增功能？

```
02-architecture.md → 05-controller-structure.md → 04-migration.md
```

### 快速查詢？

```
06-commands.md 或 07-troubleshooting.md
```

---

## 專案架構速覽

```
my-api/
├── app/                    # 應用層（核心業務邏輯）
│   ├── app.go             # 應用容器（依賴注入）
│   ├── controllers/       # 控制器
│   ├── services/          # 業務邏輯層
│   ├── repositories/      # 資料存取層
│   ├── models/            # 資料模型
│   ├── requests/          # 請求驗證
│   ├── responses/         # 回應 DTO
│   ├── middleware/        # 中間件
│   └── traits/            # 共用功能
│
├── bootstrap/              # 啟動/連線初始化
│   ├── database.go        # GORM 資料庫連接
│   ├── redis.go           # Redis 連接
│   └── logger.go          # Logger 初始化
│
├── config/                 # 配置模組
│   └── config.go          # 環境變數載入
│
├── database/               # 資料庫相關
│   ├── migrator.go        # Migration 執行器
│   └── migrations/        # Migration 檔案
│
├── routes/                 # 路由定義
│   └── api.go             # API 路由
│
├── cmd/                    # 命令行工具
│   └── migrate/           # Migration CLI
│
├── storage/                # 儲存目錄
│   └── logs/              # 日誌檔案
│
├── docs/                   # 文件（你正在看的）
│
├── main.go                 # 入口點
├── .env                    # 環境變數
└── docker-compose.yml      # Docker 配置
```

---

## 技術棧

| 類型 | 技術 |
|------|------|
| Web 框架 | Gin |
| ORM | GORM |
| 資料庫 | MySQL / PostgreSQL |
| 快取 | Redis |
| 日誌 | zerolog + lumberjack |
| 容器化 | Docker |

---

## 與 Laravel 對照表

| Laravel | Go (本專案) |
|---------|-------------|
| `app/Http/Controllers` | `app/controllers/` |
| `app/Services` | `app/services/` |
| `app/Repositories` | `app/repositories/` |
| `app/Http/Requests` | `app/requests/` |
| `app/Http/Resources` | `app/responses/` |
| `app/Providers` | `bootstrap/` |
| `database/migrations` | `database/migrations/` |
| `routes/api.php` | `routes/api.go` |
| `config/*.php` | `config/config.go` + `.env` |

---

## 變更記錄

查看 [CHANGELOG.md](CHANGELOG.md) 了解專案的版本變更歷史。

---

**最後更新：** 2026-02-03
