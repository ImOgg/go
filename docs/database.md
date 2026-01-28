# 資料庫連接設定

## 安裝依賴

```bash
# GORM 核心
go get -u gorm.io/gorm

# 資料庫驅動（按需安裝）
go get -u gorm.io/driver/mysql      # GORM MySQL 驅動
go get -u gorm.io/driver/postgres   # GORM PostgreSQL 驅動

# 原生 database/sql 驅動
go get -u github.com/go-sql-driver/mysql  # 原生 MySQL 驅動
go get -u github.com/lib/pq               # 原生 PostgreSQL 驅動

# Redis
go get -u github.com/redis/go-redis/v9

# 環境變數管理
go get -u github.com/joho/godotenv
```

## 兩種資料庫操作方式

### 🚀 方式一：GORM（ORM 框架）

**優點：**
- 自動 migration
- 簡潔的 CRUD 操作
- 關聯關係處理
- 軟刪除支援

**適合：**快速開發、標準 CRUD、關聯查詢

**使用：**
```go
// 在 main.go 中
database.InitDB()

// 查詢範例
var users []models.User
database.DB.Where("age > ?", 18).Find(&users)
```

### 📘 方式二：原生 database/sql

**優點：**
- 完全控制 SQL
- 性能更好（複雜查詢）
- 更靈活
- 無額外抽象層

**適合：**複雜查詢、性能優化、需要特定 SQL 功能

**使用：**
```go
// 在 main.go 中
database.InitRawDB()
defer database.CloseRawDB()

// 查詢範例
rows, err := database.SqlDB.Query("SELECT * FROM users WHERE age > ?", 18)
for rows.Next() {
    var user User
    rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age)
}
```

## API 路由範例

### GORM 方式
- `GET /api/gorm/users` - 取得所有使用者
- `POST /api/gorm/users` - 新增使用者

### 原生 SQL 方式
- `GET /api/sql/users` - 取得所有使用者
- `POST /api/sql/users` - 新增使用者

## 環境配置

### 1. 複製 .env.example 為 .env

```bash
cp .env.example .env
```

### 2. 編輯 .env 設定

```env
# 應用程式設定
APP_ENV=development
APP_PORT=8080

# 資料庫類型: mysql 或 postgres
DB_TYPE=mysql

# MySQL 設定
DB_HOST=host.docker.internal
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=go

# Redis 設定
REDIS_HOST=host.docker.internal
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 3. 切換資料庫

**使用 MySQL：**
```env
DB_TYPE=mysql
DB_PORT=3306
```

**使用 PostgreSQL：**
```env
DB_TYPE=postgres
DB_PORT=5432
DB_SSLMODE=disable
```

## 資料庫連接配置

### 連接字串格式

```
用戶名:密碼@tcp(主機:埠)/資料庫名?charset=utf8mb4&parseTime=True&loc=Local
```

### Docker 環境下連接主機 MySQL

在 Docker 容器中運行時，使用 `host.docker.internal` 來連接主機的 MySQL：

```go
dsn := "root:password@tcp(host.docker.internal:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
```

### 本機環境連接

在本機直接運行時，使用 `127.0.0.1` 或 `localhost`：

```go
dsn := "root:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
```

## 資料庫初始化

在 `database/db.go` 中：

```go
package database

import (
    "fmt"
    "log"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    dsn := "root:password@tcp(host.docker.internal:3306)/go?charset=utf8mb4&parseTime=True&loc=Local"
    
    var err error
    DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    
    if err != nil {
        log.Fatal("無法連接到資料庫:", err)
    }
    
    fmt.Println("資料庫連接成功！")
}
```

## 在 main.go 中使用

```go
package main

import (
    "github.com/gin-gonic/gin"
    "my-api/database"
    "my-api/routes"
)

func main() {
    // 初始化資料庫
    database.InitDB()
    
    r := gin.Default()
    routes.InitRoutes(r)
    r.Run(":8080")
}
```

## 使用資料庫

在任何地方都可以使用 `database.DB` 來操作資料庫：

```go
import "my-api/database"

// 查詢範例
var users []User
database.DB.Find(&users)

// 新增資料
user := User{Name: "John", Email: "john@example.com"}
database.DB.Create(&user)

// 更新資料
database.DB.Model(&user).Update("Name", "John Doe")

// 刪除資料
database.DB.Delete(&user)
```

## ORM 與資料表遷移 (AutoMigrate)

### 定義 Model

在 `models/user.go` 中：

```go
package models

import "gorm.io/gorm"

type User struct {
    gorm.Model  // 包含 ID, CreatedAt, UpdatedAt, DeletedAt
    Name  string `json:"name" gorm:"type:varchar(100);not null"`
    Email string `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
    Age   int    `json:"age"`
}
```

### GORM 標籤說明

- `gorm:"type:varchar(100)"` - 指定欄位類型
- `gorm:"not null"` - 設為必填
- `gorm:"uniqueIndex"` - 建立唯一索引
- `gorm:"default:0"` - 設定預設值
- `json:"name"` - JSON 序列化時的欄位名稱

### 自動遷移

在 `database/db.go` 的 `InitDB()` 中加入：

```go
import "my-api/models"

func InitDB() {
    // ... 資料庫連接 ...
    
    // 自動遷移
    err = DB.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatal("資料表遷移失敗:", err)
    }
    
    fmt.Println("資料表遷移完成！")
}
```

### 遷移多個 Model

```go
DB.AutoMigrate(
    &models.User{},
    &models.Product{},
    &models.Order{},
)
```

### CRUD 操作範例

```go
// Create - 新增
user := models.User{Name: "John", Email: "john@example.com", Age: 25}
database.DB.Create(&user)

// Read - 查詢
var user models.User
database.DB.First(&user, 1)  // 根據 ID 查詢
database.DB.First(&user, "email = ?", "john@example.com")  // 根據條件查詢

// 查詢多筆
var users []models.User
database.DB.Find(&users)
database.DB.Where("age > ?", 18).Find(&users)

// Update - 更新
database.DB.Model(&user).Update("Name", "John Doe")
database.DB.Model(&user).Updates(models.User{Name: "Jane", Age: 30})

// Delete - 刪除（軟刪除）
database.DB.Delete(&user)

// 永久刪除
database.DB.Unscoped().Delete(&user)
```

## 常見問題

### 連接被拒絕 (connection refused)

- **Docker 環境**：確保使用 `host.docker.internal` 而非 `127.0.0.1`
- **本機環境**：確認 MySQL 服務正在運行
- 檢查埠號是否正確（預設 3306）

### 認證失敗

- 確認用戶名和密碼正確
- 檢查 MySQL 用戶權限：`GRANT ALL PRIVILEGES ON dbname.* TO 'username'@'%';`

### 資料庫不存在

- 先在 MySQL 中創建資料庫：`CREATE DATABASE dbname;`
