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

| 檔案 | 測試內容 |
|------|----------|
| `app/utils/jwt_test.go` | 密碼加密與驗證 |

---

## 進階主題

### Mock 測試

Go 可以使用介面來實現 Mock：

```go
// 定義介面
type UserRepository interface {
    FindByID(id uint) (*User, error)
}

// 測試用的 Mock
type mockUserRepo struct {
    user *User
    err  error
}

func (m *mockUserRepo) FindByID(id uint) (*User, error) {
    return m.user, m.err
}

// 測試
func TestUserService_GetUser(t *testing.T) {
    mockRepo := &mockUserRepo{
        user: &User{ID: 1, Name: "Test"},
        err:  nil,
    }

    service := NewUserService(mockRepo)
    user, err := service.GetUserByID(1)
    // ...
}
```

### 常用測試套件

| 套件 | 用途 |
|------|------|
| `testify` | 斷言、Mock |
| `gomock` | Google 的 Mock 框架 |
| `httptest` | HTTP 測試（內建） |
| `sqlmock` | 資料庫 Mock |

---

**最後更新：** 2026-02-02
