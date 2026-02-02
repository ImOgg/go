# Go 測試指南

> Go 內建測試框架，無需額外安裝套件（不像 PHP 需要 PHPUnit）

---

## 與 Laravel PHPUnit 對比

| 項目 | Laravel PHPUnit | Go |
|------|-----------------|-----|
| 測試目錄 | `tests/Unit/`, `tests/Feature/` | **與原始碼同目錄** |
| 檔案命名 | `*Test.php` | `*_test.go` |
| 測試框架 | 需安裝 PHPUnit | 內建 `testing` 套件 |
| 執行命令 | `php artisan test` | `go test` |
| 測試方法 | `public function test_xxx()` | `func TestXxx(t *testing.T)` |

---

## 測試檔案位置

Go 的測試檔案放在 **與被測試程式碼相同的目錄**：

```
app/
├── utils/
│   ├── jwt.go           # 主程式
│   └── jwt_test.go      # 測試（同一目錄）
├── services/
│   ├── user_service.go
│   └── user_service_test.go
```

---

## 基本測試結構

### 1. 最簡單的測試

```go
// jwt_test.go
package utils

import "testing"

func TestHashPassword(t *testing.T) {
    password := "mySecret123"

    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("HashPassword() 發生錯誤: %v", err)
    }

    if hash == "" {
        t.Error("HashPassword() 回傳空字串")
    }
}
```

### 2. Table-Driven 測試（推薦）

類似 PHPUnit 的 Data Provider：

```go
func TestCheckPassword_TableDriven(t *testing.T) {
    originalPassword := "testPassword"
    hash, _ := HashPassword(originalPassword)

    // 定義測試案例
    tests := []struct {
        name     string // 測試名稱
        password string // 輸入
        want     bool   // 預期結果
    }{
        {"正確密碼", originalPassword, true},
        {"錯誤密碼", "wrongPassword", false},
        {"空密碼", "", false},
    }

    // 執行每個測試案例
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := CheckPassword(tt.password, hash)
            if got != tt.want {
                t.Errorf("CheckPassword(%q) = %v, want %v",
                    tt.password, got, tt.want)
            }
        })
    }
}
```

---

## 常用測試指令

```bash
# 進入容器執行（本專案使用 Docker）
docker exec my-go-app go test -v ./app/utils/...

# 測試整個專案
docker exec my-go-app go test -v ./...

# 只執行特定測試函數
docker exec my-go-app go test -v -run TestHashPassword ./app/utils/...

# 查看測試覆蓋率
docker exec my-go-app go test -cover ./...

# 產生覆蓋率報告
docker exec my-go-app go test -coverprofile=coverage.out ./...
docker exec my-go-app go tool cover -html=coverage.out -o coverage.html
```

---

## testing 套件常用方法

| 方法 | 說明 | 類似 PHPUnit |
|------|------|--------------|
| `t.Error(msg)` | 標記失敗，繼續執行 | `$this->fail()` |
| `t.Errorf(format, args)` | 格式化錯誤訊息 | `$this->fail(sprintf())` |
| `t.Fatal(msg)` | 標記失敗，立即停止 | `$this->fail()` + return |
| `t.Fatalf(format, args)` | 格式化 + 立即停止 | - |
| `t.Log(msg)` | 輸出日誌（-v 時顯示） | `echo` |
| `t.Logf(format, args)` | 格式化日誌 | `printf` |
| `t.Run(name, fn)` | 執行子測試 | Data Provider 的每個 case |
| `t.Skip(msg)` | 跳過此測試 | `$this->markTestSkipped()` |

---

## 測試命名慣例

```go
// 基本格式
func TestFunctionName(t *testing.T) {}

// 測試特定情境
func TestFunctionName_Scenario(t *testing.T) {}

// 範例
func TestHashPassword(t *testing.T) {}
func TestHashPassword_EmptyPassword(t *testing.T) {}
func TestCheckPassword_WrongPassword(t *testing.T) {}
```

---

## 專案現有測試

| 檔案 | 測試內容 | 類型 |
|------|----------|------|
| `app/utils/jwt_test.go` | 密碼加密與驗證 | 單元測試 |
| `app/services/user_service_test.go` | Service 層 CRUD | Mock 測試 |

---

## 與 phpunit.xml 的對比

Go **沒有** 類似 `phpunit.xml` 的官方配置檔，但有替代方案：

| Laravel/PHPUnit | Go 替代方案 |
|-----------------|------------|
| `phpunit.xml` 設定測試資料庫 | 使用 Mock（不碰資料庫） |
| `DB_DATABASE=testing` | 環境變數或 `.env.testing` |
| `RefreshDatabase` trait | Mock 或獨立測試資料庫 |
| `setUp()` / `tearDown()` | `TestMain()` 函數 |

---

## Mock 測試（推薦）

### 為什麼用 Mock？

- **不碰真實資料庫** → 測試快速、安全
- **可控制回傳值** → 方便測試各種情境（成功、失敗、邊界條件）
- **依賴注入** → 本專案已使用 interface，非常適合 Mock

### Mock 原理

```
┌─────────────────┐      ┌─────────────────┐
│   UserService   │ ───▶ │ UserRepository  │ (interface)
└─────────────────┘      └─────────────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │                           │
              ┌─────▼─────┐              ┌──────▼──────┐
              │ 正式環境  │              │  測試環境   │
              │ userRepo  │              │ mockRepo    │
              │ (用 GORM) │              │ (用 map)    │
              └───────────┘              └─────────────┘
```

### 本專案 Mock 範例

參考 `app/services/user_service_test.go`：

```go
// Mock Repository - 實作 UserRepository interface
type mockUserRepository struct {
    users    map[uint]*models.User  // 用 map 模擬資料庫
    emailMap map[string]*models.User
    nextID   uint
}

// 實作 interface 的所有方法
func (m *mockUserRepository) Create(user *models.User) error {
    user.ID = m.nextID
    m.nextID++
    m.users[user.ID] = user
    m.emailMap[user.Email] = user
    return nil
}

func (m *mockUserRepository) FindByID(id uint) (*models.User, error) {
    if user, ok := m.users[id]; ok {
        return user, nil
    }
    return nil, errors.New("record not found")
}

// ... 其他方法

// 測試
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := newMockUserRepository()  // 使用 Mock
    service := NewUserService(mockRepo)  // 注入 Mock

    req := &requests.CreateUserRequest{
        Name:  "張三",
        Email: "test@example.com",
        Age:   25,
    }

    resp, err := service.CreateUser(req)

    if err != nil {
        t.Errorf("不預期的錯誤: %v", err)
    }
    if resp.Name != "張三" {
        t.Errorf("Name 不符")
    }
}
```

---

## 測試環境設定（整合測試用）

如果需要真正的資料庫測試，可以設定測試環境：

### 方式 1：TestMain 初始化

```go
// user_repository_test.go
package repositories

import (
    "os"
    "testing"
)

func TestMain(m *testing.M) {
    // 設定測試環境變數
    os.Setenv("DB_DATABASE", "my_api_test")

    // 初始化測試資料庫
    // bootstrap.InitTestDB()

    // 執行測試
    code := m.Run()

    // 清理
    // bootstrap.CleanupTestDB()

    os.Exit(code)
}
```

### 方式 2：.env.testing

```env
# .env.testing
DB_HOST=localhost
DB_DATABASE=my_api_test
DB_USERNAME=root
DB_PASSWORD=password
```

---

## 常用測試套件

| 套件 | 用途 | 安裝 |
|------|------|------|
| `testify` | 斷言、Mock | `go get github.com/stretchr/testify` |
| `gomock` | Google Mock 框架 | `go get github.com/golang/mock/gomock` |
| `httptest` | HTTP 測試 | 內建 |
| `sqlmock` | 資料庫 Mock | `go get github.com/DATA-DOG/go-sqlmock` |

---

## 測試覆蓋率

```bash
# 查看覆蓋率百分比
docker exec my-go-app go test -cover ./...

# 產生詳細覆蓋率報告
docker exec my-go-app go test -coverprofile=coverage.out ./...

# 在瀏覽器查看（需將檔案複製出來）
docker cp my-go-app:/app/coverage.out ./coverage.out
go tool cover -html=coverage.out
```

---

**最後更新：** 2026-02-02
