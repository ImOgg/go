# Migration 系統使用指南

## 📖 什麼是 Migration？

Migration 是資料庫的版本控制系統，類似 Git 對程式碼的管理。

### 與 Laravel 對比

| Laravel | 本專案 (Go) |
|---------|-------------|
| `php artisan make:migration` | `go run cmd/migrate/main.go make` |
| `php artisan migrate` | `go run cmd/migrate/main.go migrate` |
| `php artisan migrate:rollback` | `go run cmd/migrate/main.go rollback` |
| `php artisan migrate:status` | `go run cmd/migrate/main.go status` |

---

## 🚀 快速開始

### 1. 執行 Migration

**在 Docker 容器內：**
```bash
# 進入容器
docker exec -it my-go-app bash

# 執行所有待執行的 migrations
go run cmd/migrate/main.go migrate
```

**輸出範例：**
```
🚀 執行 Migration: 000001 - create_users_table
✓ 建立 users 表成功
✅ 所有 Migrations 執行成功！
```

### 2. 查看狀態

```bash
go run cmd/migrate/main.go status
```

**輸出範例：**
```
📊 Migration 狀態:

  ✅ [000001] create_users_table - 已執行
  ⏳ [000002] add_phone_to_users - 待執行
```

### 3. 回滾

```bash
go run cmd/migrate/main.go rollback
```

**輸出範例：**
```
⏮️  回滾 Migration: 000001 - create_users_table
✓ 刪除 users 表成功
✅ Rollback 成功！
```

---

## 📁 檔案結構

```
database/
  ├── migrations/
  │   ├── migration.go                      # Migration 介面定義
  │   ├── registry.go                       # 註冊器（管理所有 migrations）
  │   └── 000001_create_users_table.go     # 實際的 migration（一個檔案包含 Up 和 Down）
  │
  └── migrate_simple.go                     # Migration 執行引擎

cmd/
  └── migrate/
      └── main.go                           # Migration 命令行工具

main.go                                     # 啟動時自動執行 migration（可選）
```

---

## ✍️ 建立新的 Migration

### 方式 1：手動建立（推薦）

在 `database/migrations/` 目錄下建立新檔案：

**檔案名稱：** `000002_add_phone_to_users.go`

```go
package migrations

import (
	"database/sql"
	"fmt"
)

// AddPhoneToUsers - 新增 phone 欄位到 users 表
type AddPhoneToUsers struct {
	BaseMigration
}

func init() {
	Register(&AddPhoneToUsers{
		BaseMigration: BaseMigration{
			version:     "000002",
			description: "add_phone_to_users",
		},
	})
}

// Up - 執行 migration（新增欄位）
func (m *AddPhoneToUsers) Up(db *sql.DB) error {
	query := `
		ALTER TABLE users 
		ADD COLUMN phone VARCHAR(20) AFTER email,
		ADD INDEX idx_phone (phone);
	`
	
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("新增 phone 欄位失敗: %v", err)
	}
	
	fmt.Println("✓ 新增 phone 欄位成功")
	return nil
}

// Down - 回滾 migration（移除欄位）
func (m *AddPhoneToUsers) Down(db *sql.DB) error {
	query := `
		ALTER TABLE users 
		DROP INDEX idx_phone,
		DROP COLUMN phone;
	`
	
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("移除 phone 欄位失敗: %v", err)
	}
	
	fmt.Println("✓ 移除 phone 欄位成功")
	return nil
}
```

### 方式 2：使用命令生成模板（待實作）

```bash
go run cmd/migrate/main.go make add_phone_to_users
```

---

## 📝 Migration 範例

### 1. 建立資料表

```go
// Up
func (m *CreateProductsTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS products (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			stock INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_name (name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	
	_, err := db.Exec(query)
	return err
}

// Down
func (m *CreateProductsTable) Down(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS products;")
	return err
}
```

### 2. 新增欄位

```go
// Up
func (m *AddAvatarToUsers) Up(db *sql.DB) error {
	query := `ALTER TABLE users ADD COLUMN avatar VARCHAR(255) AFTER email;`
	_, err := db.Exec(query)
	return err
}

// Down
func (m *AddAvatarToUsers) Down(db *sql.DB) error {
	query := `ALTER TABLE users DROP COLUMN avatar;`
	_, err := db.Exec(query)
	return err
}
```

### 3. 修改欄位類型

```go
// Up
func (m *ChangeEmailLength) Up(db *sql.DB) error {
	query := `ALTER TABLE users MODIFY COLUMN email VARCHAR(320) NOT NULL;`
	_, err := db.Exec(query)
	return err
}

// Down
func (m *ChangeEmailLength) Down(db *sql.DB) error {
	query := `ALTER TABLE users MODIFY COLUMN email VARCHAR(255) NOT NULL;`
	_, err := db.Exec(query)
	return err
}
```

### 4. 新增索引

```go
// Up
func (m *AddIndexToUsers) Up(db *sql.DB) error {
	query := `CREATE INDEX idx_created_at ON users(created_at);`
	_, err := db.Exec(query)
	return err
}

// Down
func (m *AddIndexToUsers) Down(db *sql.DB) error {
	query := `DROP INDEX idx_created_at ON users;`
	_, err := db.Exec(query)
	return err
}
```

### 5. 插入初始資料

```go
// Up
func (m *SeedDefaultUsers) Up(db *sql.DB) error {
	query := `
		INSERT INTO users (name, email, age) VALUES 
		('Admin', 'admin@example.com', 30),
		('Test User', 'test@example.com', 25);
	`
	_, err := db.Exec(query)
	return err
}

// Down
func (m *SeedDefaultUsers) Down(db *sql.DB) error {
	query := `DELETE FROM users WHERE email IN ('admin@example.com', 'test@example.com');`
	_, err := db.Exec(query)
	return err
}
```

---

## 🎯 重要概念

### 1. 版本號規則

- 格式：`000001`, `000002`, `000003`...
- **必須遞增**，不能跳號
- **不能重複**

### 2. Migration 結構

```go
type YourMigration struct {
	BaseMigration  // 繼承基礎功能
}

func init() {
	Register(&YourMigration{  // 註冊到系統
		BaseMigration: BaseMigration{
			version:     "000001",        // 版本號
			description: "create_users",  // 描述
		},
	})
}

func (m *YourMigration) Up(db *sql.DB) error {
	// 執行變更
}

func (m *YourMigration) Down(db *sql.DB) error {
	// 回滾變更
}
```

### 3. 執行順序

Migration 按版本號**由小到大**執行：
```
000001 → 000002 → 000003
```

Rollback 按版本號**由大到小**回滾：
```
000003 → 000002 → 000001
```

---

## 📊 migrations 記錄表

系統會自動建立 `migrations` 表來追蹤執行狀態：

```sql
CREATE TABLE migrations (
    version VARCHAR(14) PRIMARY KEY,
    description VARCHAR(255) NOT NULL,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**查看已執行的 migrations：**
```sql
SELECT * FROM migrations ORDER BY version;
```

**結果範例：**
```
+--------+--------------------+---------------------+
| version| description        | executed_at         |
+--------+--------------------+---------------------+
| 000001 | create_users_table | 2026-01-28 10:30:15 |
| 000002 | add_phone_to_users | 2026-01-28 11:20:45 |
+--------+--------------------+---------------------+
```

---

## ⚙️ 整合方式

### 方式 1：手動執行（開發時推薦）

在容器內手動執行：
```bash
go run cmd/migrate/main.go migrate
```

### 方式 2：自動執行（生產環境推薦）

在 `main.go` 中已設定，每次啟動自動執行：

```go
func main() {
	config.LoadConfig()
	database.InitDB()
	
	// 自動執行 migrations
	if err := database.RunMigrations(); err != nil {
		log.Println("⚠️  Migration 警告:", err)
	}
	
	// ... 啟動服務
}
```

**優點：**
- 部署時自動更新資料庫
- 不需要手動執行 migrate 命令

**缺點：**
- 啟動時間稍長
- 如果 migration 失敗，服務仍會啟動

---

## 🔒 最佳實踐

### 1. 永遠提供 Down 方法

每個 Up 都要有對應的 Down，確保可以回滾：

```go
// ✅ 好的做法
func (m *Migration) Up(db *sql.DB) error {
	// 新增欄位
}

func (m *Migration) Down(db *sql.DB) error {
	// 刪除欄位（與 Up 相反）
}
```

### 2. 測試 Migration

在開發環境先測試：

```bash
# 執行
go run cmd/migrate/main.go migrate

# 確認資料表正確
# mysql -u root -p mydb

# 測試回滾
go run cmd/migrate/main.go rollback

# 再次執行確保可重複
go run cmd/migrate/main.go migrate
```

### 3. 不要修改已執行的 Migration

如果需要改動，建立新的 migration：

```go
// ❌ 不要修改 000001_create_users_table.go

// ✅ 建立新的 migration
// 000003_modify_users_table.go
```

### 4. 使用事務（重要變更時）

```go
func (m *Migration) Up(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	
	// 執行多個操作
	if _, err := tx.Exec("ALTER TABLE..."); err != nil {
		tx.Rollback()
		return err
	}
	
	if _, err := tx.Exec("CREATE INDEX..."); err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit()
}
```

### 5. 版本號命名

```
000001_create_users_table.go
000002_add_phone_to_users.go
000003_create_products_table.go
000004_add_index_to_users.go
```

使用描述性名稱，一看就知道做什麼。

---

## 🐛 常見問題

### Q1: Migration 執行失敗怎麼辦？

1. 查看錯誤訊息
2. 檢查 SQL 語法
3. 檢查資料庫連線
4. 手動修復資料庫
5. 更新 `migrations` 表的記錄

### Q2: 如何重新執行某個 Migration？

```sql
-- 刪除記錄
DELETE FROM migrations WHERE version = '000001';
```

然後重新執行：
```bash
go run cmd/migrate/main.go migrate
```

### Q3: 多人協作時版本號衝突怎麼辦？

1. Git pull 最新代碼
2. 重新命名你的 migration 版本號
3. 更新 `version` 欄位
4. 提交代碼

### Q4: 可以跳過某個 Migration 嗎？

可以手動插入記錄：

```sql
INSERT INTO migrations (version, description) 
VALUES ('000002', 'skipped_migration');
```

---

## 🆚 與其他方案對比

| 方案 | 優點 | 缺點 |
|------|------|------|
| **本專案（Go Migration）** | ✅ 類似 Laravel<br>✅ 可 Rollback<br>✅ 版本控制<br>✅ 一個檔案包含 Up/Down | ⚠️ 需要手動建立檔案 |
| **GORM AutoMigrate** | ✅ 自動建表<br>✅ 簡單快速 | ❌ 不能刪除欄位<br>❌ 不能 Rollback<br>❌ 無版本控制 |
| **golang-migrate (SQL)** | ✅ 成熟穩定 | ❌ 需要兩個檔案（.up.sql + .down.sql）<br>❌ 外部依賴 |

---

## 🎓 學習資源

- [GORM 官方文檔](https://gorm.io/)
- [database/sql 套件](https://pkg.go.dev/database/sql)
- [MySQL 文檔](https://dev.mysql.com/doc/)

---

## 📚 相關文檔

- [專案設置](setup.md)
- [常用命令](commands.md)
- [資料庫操作](database.md)
- [Controller 結構](controller-structure.md)

---

**最後更新：** 2026-01-28
