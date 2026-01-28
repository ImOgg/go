# Controller 結構組織筆記

## 📁 目錄結構

### 重構前（不推薦）
```
controllers/
  ├── hello.go    // 包含 HelloHandler, GetUserByName, Search, TestHandler
  ├── user.go     // 包含 GetUsersGORM, CreateUserGORM, GetUsersSQL, CreateUserSQL
  └── test.go
```

**問題：**
- 所有函數都在 `controllers` 命名空間，難以區分
- 不知道去哪個文件修改特定功能
- 函數命名容易衝突

### 重構後（推薦）✅
```
controllers/
  ├── user/
  │   └── user.go      // 使用者相關功能
  ├── hello/
  │   └── hello.go     // Hello 相關功能
  └── test/
      └── test.go      // 測試相關功能
```

**優點：**
- 功能分類清晰
- 類似 Laravel 的 Controller 組織方式
- IDE 自動提示會顯示包名
- 避免命名衝突

---

## 🔧 實作方式

### 1. 建立子包結構

每個子包都是一個資料夾，包含一個或多個 `.go` 文件：

**controllers/user/user.go**
```go
package user  // 注意：package 名稱是資料夾名

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// 函數名稱首字母大寫 = Public（外部可訪問）
func GetUsersGORM(c *gin.Context) {
    // ...
}

func CreateUserGORM(c *gin.Context) {
    // ...
}
```

### 2. 在路由中引入

**routes/router.go**
```go
package routes

import (
    "github.com/gin-gonic/gin"
    "my-api/controllers/user"    // 引入 user 子包
    "my-api/controllers/hello"   // 引入 hello 子包
    "my-api/controllers/test"    // 引入 test 子包
)

func InitRoutes(r *gin.Engine) {
    // 使用時：包名.函數名
    r.GET("/api/users", user.GetUsersGORM)
    r.GET("/hello", hello.Handler)
    r.GET("/test", test.Handler)
}
```

---

## 🆚 與 Laravel 的對比

### Laravel 寫法
```php
<?php
namespace App\Http\Controllers;

class UserController extends Controller
{
    public function index()
    {
        // ...
    }
    
    public function store()
    {
        // ...
    }
}
```

**路由：**
```php
use App\Http\Controllers\UserController;

Route::get('/users', [UserController::class, 'index']);
Route::post('/users', [UserController::class, 'store']);
```

### Go 寫法
```go
package user

import "github.com/gin-gonic/gin"

// 注意：Go 沒有 class，只有 package 和 function
func GetList(c *gin.Context) {
    // ...
}

func Create(c *gin.Context) {
    // ...
}
```

**路由：**
```go
import "my-api/controllers/user"

r.GET("/users", user.GetList)
r.POST("/users", user.Create)
```

---

## 📌 重要概念

### 1. Package（包）vs Class（類）

| Laravel | Go |
|---------|-----|
| `UserController` class | `user` package |
| `$this->method()` | `user.Method()` |
| 類方法 | 獨立函數 |

### 2. Public vs Private

**Go 的規則非常簡單：**
- **首字母大寫** = Public（外部可訪問）
  ```go
  func GetUsers() {}      // ✅ 可以被其他包使用
  ```
  
- **首字母小寫** = Private（僅包內使用）
  ```go
  func validateUser() {}  // ❌ 只能在同一個 package 內使用
  ```

**Laravel：**
```php
public function index() {}     // public
private function helper() {}   // private
protected function validate()  // protected
```

### 3. 引入方式

**Go：**
```go
import (
    "my-api/controllers/user"    // 引入整個包
    "my-api/controllers/hello"
)

// 使用：包名.函數名
user.GetList()
hello.Handler()
```

**Laravel：**
```php
use App\Http\Controllers\UserController;  // 引入單個類
use App\Http\Controllers\HelloController;

// 使用：類名::方法
UserController::index();
```

---

## 🎯 最佳實踐

### 1. 命名規範

```go
// ✅ 好的命名
package user

func GetList(c *gin.Context) {}      // 簡潔，因為已經在 user 包內
func Create(c *gin.Context) {}
func GetByID(c *gin.Context) {}

// ❌ 避免重複
func UserGetList() {}  // 不需要，已經在 user 包內了
```

### 2. 文件組織

```
controllers/
  ├── user/
  │   ├── user.go           // 基本 CRUD
  │   ├── user_gorm.go      // GORM 相關（可選）
  │   └── user_validator.go // 驗證邏輯（可選）
  ├── product/
  │   └── product.go
  └── order/
      └── order.go
```

### 3. 路由組織

```go
func InitRoutes(r *gin.Engine) {
    // 按功能分組
    userRoutes := r.Group("/users")
    {
        userRoutes.GET("", user.GetList)
        userRoutes.POST("", user.Create)
        userRoutes.GET("/:id", user.GetByID)
        userRoutes.PUT("/:id", user.Update)
        userRoutes.DELETE("/:id", user.Delete)
    }
    
    productRoutes := r.Group("/products")
    {
        productRoutes.GET("", product.GetList)
        productRoutes.POST("", product.Create)
    }
}
```

---

## 🔍 常見問題

### Q1: 為什麼不用 struct 模擬 class？

可以這樣做，但不符合 Go 的慣用寫法：

```go
// 可以但不推薦
type UserController struct {}

func (uc *UserController) GetList(c *gin.Context) {}
```

**Go 的哲學：**
- 簡單、直接
- 組合優於繼承
- 函數式編程

### Q2: 如何共享邏輯？

**建立共用函數或中間件：**

```go
// controllers/common/validator.go
package common

func ValidateEmail(email string) bool {
    // 驗證邏輯
}
```

**使用：**
```go
package user

import "my-api/controllers/common"

func Create(c *gin.Context) {
    if !common.ValidateEmail(email) {
        // ...
    }
}
```

### Q3: 包名和資料夾名必須一樣嗎？

**不一定，但強烈建議一致：**

```go
// 資料夾：controllers/user/
package user  // ✅ 推薦：包名與資料夾名一致

package userController  // ❌ 可以但不推薦
```

---

## 📚 相關文檔

- [專案設置](setup.md)
- [常用命令](commands.md)
- [資料庫操作](database.md)
- [常見問題](troubleshooting.md)

---

## 🚀 快速參考

### 建立新的 Controller

1. **建立資料夾和文件：**
   ```bash
   mkdir controllers/product
   touch controllers/product/product.go
   ```

2. **編寫代碼：**
   ```go
   package product
   
   import "github.com/gin-gonic/gin"
   
   func GetList(c *gin.Context) {
       // ...
   }
   ```

3. **在路由中引入：**
   ```go
   import "my-api/controllers/product"
   
   r.GET("/products", product.GetList)
   ```

### 測試路由

```bash
# 啟動服務
go run main.go

# 測試端點
curl http://localhost:8080/test
curl http://localhost:8080/users/張三
curl http://localhost:8080/hello
curl http://localhost:8080/search?keyword=golang
```

---

**最後更新：** 2026-01-28
