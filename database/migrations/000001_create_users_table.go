package migrations

import (
	"database/sql"
	"fmt"
)

// CreateUsersTable - 建立 users 資料表
type CreateUsersTable struct {
	BaseMigration
}

func init() {
	Register(&CreateUsersTable{
		BaseMigration: BaseMigration{
			version:     "000001",
			description: "create_users_table",
		},
	})
}

// Up - 執行 migration
func (m *CreateUsersTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			age INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			INDEX idx_email (email),
			INDEX idx_deleted_at (deleted_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("建立 users 表失敗: %v", err)
	}
	
	fmt.Println("✓ 建立 users 表成功")
	return nil
}

// Down - 回滾 migration
func (m *CreateUsersTable) Down(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS users;`
	
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("刪除 users 表失敗: %v", err)
	}
	
	fmt.Println("✓ 刪除 users 表成功")
	return nil
}
