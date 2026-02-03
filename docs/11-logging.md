# Go 結構化日誌系統

本專案使用 [zerolog](https://github.com/rs/zerolog) 實作類似 Laravel Log 的結構化日誌系統。

## 目錄

- [環境變數設定](#環境變數設定)
- [基本使用](#基本使用)
- [在 Controller 中使用](#在-controller-中使用)
- [在 Service 中使用](#在-service-中使用)
- [日誌等級](#日誌等級)
- [輸出範例](#輸出範例)
- [檔案結構](#檔案結構)

---

## 環境變數設定

在 `.env` 檔案中設定：

```env
# 日誌等級: debug, info, warn, error, fatal
LOG_LEVEL=debug

# 輸出格式: console (彩色), json
LOG_FORMAT=console

# 輸出目標: stdout, file, both
LOG_OUTPUT=both

# 日誌檔案路徑
LOG_FILE_PATH=storage/logs/app.log

# 日誌輪替設定
LOG_MAX_SIZE=100      # 單檔最大 MB
LOG_MAX_BACKUPS=30    # 保留檔案數
LOG_MAX_AGE=30        # 保留天數
LOG_COMPRESS=true     # 壓縮舊檔
```

### 建議配置

| 環境 | LOG_LEVEL | LOG_FORMAT | LOG_OUTPUT |
|------|-----------|------------|------------|
| 開發 | debug | console | stdout 或 both |
| 測試 | info | json | both |
| 生產 | info | json | file 或 both |

---

## 基本使用

### 方式一：從 Gin Context 取得（推薦）

在 Controller 或 Middleware 中，使用這個方式可以自動帶上 `request_id`：

```go
import "my-api/app/pkg/logger"

func (ctrl *UserController) Store(c *gin.Context) {
    // 從 context 取得 Logger（已帶有 request_id）
    log := logger.FromGinContext(c)

    log.Info("使用者建立成功", map[string]interface{}{
        "user_id": user.ID,
        "email":   user.Email,
    })
}
```

### 方式二：使用全域 Logger

在 Service 或其他地方，使用全域 Logger：

```go
import "my-api/bootstrap"

func (s *userService) CreateUser(req *requests.CreateUserRequest) error {
    bootstrap.Log.Info("建立使用者", map[string]interface{}{
        "email": req.Email,
    })
    // ...
}
```

### 方式三：使用 logger.Global()

```go
import "my-api/app/pkg/logger"

logger.Global().Error("發生錯誤", map[string]interface{}{
    "error": err.Error(),
})
```

---

## 在 Controller 中使用

```go
package controllers

import (
    "my-api/app/pkg/logger"
)

func (ctrl *UserController) Store(c *gin.Context) {
    log := logger.FromGinContext(c)

    // 驗證失敗
    if err := req.Validate(c); err != nil {
        log.Warning("驗證失敗", map[string]interface{}{
            "errors": err.Error(),
        })
        return
    }

    // 開始處理（Debug 等級，僅開發時顯示）
    log.Debug("開始建立使用者", map[string]interface{}{
        "email": req.Email,
        "name":  req.Name,
    })

    // 操作失敗
    if err != nil {
        log.Error("使用者建立失敗", map[string]interface{}{
            "email": req.Email,
            "error": err.Error(),
        })
        return
    }

    // 操作成功
    log.Info("使用者建立成功", map[string]interface{}{
        "user_id": user.ID,
        "email":   user.Email,
    })
}
```

---

## 在 Service 中使用

```go
package services

import "my-api/bootstrap"

func (s *authService) Login(req *requests.LoginRequest) (*responses.AuthResponse, error) {
    log := bootstrap.Log

    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        log.Warning("登入失敗: 使用者不存在", map[string]interface{}{
            "email": req.Email,
        })
        return nil, errors.New("帳號或密碼錯誤")
    }

    log.Info("使用者登入成功", map[string]interface{}{
        "user_id": user.ID,
        "email":   user.Email,
    })

    return response, nil
}
```

---

## 日誌等級

| 等級 | 方法 | 用途 |
|------|------|------|
| Debug | `log.Debug()` | 開發除錯資訊，生產環境不顯示 |
| Info | `log.Info()` | 一般資訊，如成功操作 |
| Warning | `log.Warning()` | 警告，如驗證失敗、找不到資源 |
| Error | `log.Error()` | 錯誤，如資料庫操作失敗 |
| Fatal | `log.Fatal()` | 致命錯誤，會終止程式 |

### 使用建議

```go
// Debug - 開發除錯
log.Debug("準備執行查詢", map[string]interface{}{"sql": query})

// Info - 重要操作成功
log.Info("使用者登入成功", map[string]interface{}{"user_id": 1})

// Warning - 可預期的失敗（不是 bug）
log.Warning("登入失敗：密碼錯誤", map[string]interface{}{"email": "test@test.com"})

// Error - 非預期的錯誤（可能是 bug）
log.Error("資料庫連線失敗", map[string]interface{}{"error": err.Error()})

// Fatal - 無法恢復的錯誤（程式會終止）
log.Fatal("無法載入設定檔", map[string]interface{}{"path": configPath})
```

---

## 輸出範例

### Console 格式（開發環境）

```
2026-02-03 11:22:44 INF Logger 初始化成功 format=console level=debug output=both
2026-02-03 11:22:44 INF Server starting env=development port=8080
2026-02-03 11:22:55 INF Request request_id=29b62573-... method=GET path=/health status=200 latency=29.457µs
2026-02-03 11:23:16 WRN Client Error request_id=934e2c4d-... method=POST path=/api/login status=401
```

### JSON 格式（生產環境）

```json
{"level":"info","time":"2026-02-03T11:22:44Z","message":"Logger 初始化成功","format":"json","output":"file"}
{"level":"info","time":"2026-02-03T11:22:55Z","request_id":"29b62573-...","method":"GET","path":"/health","status":200,"message":"Request"}
{"level":"warn","time":"2026-02-03T11:23:16Z","request_id":"934e2c4d-...","method":"POST","path":"/api/login","status":401,"message":"Client Error"}
```

---

## 檔案結構

```
app/
├── pkg/
│   └── logger/
│       ├── logger.go      # Logger 核心封裝
│       └── context.go     # Gin Context 輔助函數
├── middleware/
│   ├── logger.go          # HTTP 請求日誌中間件
│   └── request_id.go      # Request ID 中間件
bootstrap/
└── logger.go              # Logger 初始化
storage/
└── logs/
    └── app.log            # 日誌檔案
```

---

## 進階用法

### 建立子 Logger

```go
// 帶有固定欄位的子 Logger
userLog := log.WithContext(map[string]interface{}{
    "module": "user",
    "action": "create",
})
userLog.Info("開始處理")  // 自動帶上 module 和 action

// 帶有 Request ID 的子 Logger
reqLog := log.WithRequestID("abc-123")

// 帶有錯誤的子 Logger
errLog := log.WithError(err)
errLog.Error("操作失敗")  // 自動帶上 error 欄位
```

### 與 Laravel Log 對照

| Laravel | Go (本專案) |
|---------|-------------|
| `Log::info('message')` | `log.Info("message", nil)` |
| `Log::info('message', ['key' => 'value'])` | `log.Info("message", map[string]interface{}{"key": "value"})` |
| `Log::error('error', ['exception' => $e])` | `log.Error("error", map[string]interface{}{"error": err.Error()})` |
| `Log::channel('daily')->info()` | 透過環境變數設定輸出目標 |

---

## 依賴套件

```go
github.com/rs/zerolog                  // 高性能結構化日誌
gopkg.in/natefinch/lumberjack.v2       // 日誌檔案輪替
github.com/google/uuid                 // UUID 生成（Request ID）
```

---

**最後更新：** 2026-02-03
