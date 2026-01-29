package migrations

import (
	"database/sql"
	"fmt"
	"strings"
)

// AddPasswordToUsers - 新增 password 欄位到 users 資料表
type AddPasswordToUsers struct {
	BaseMigration
}

func init() {
	Register(&AddPasswordToUsers{
		BaseMigration: BaseMigration{
			version:     "000002",
			description: "add_password_to_users",
		},
	})
}

// Up - 執行 migration
func (m *AddPasswordToUsers) Up(db *sql.DB) error {
	// 直接嘗試新增欄位，如果已存在會報錯
	query := `ALTER TABLE users ADD COLUMN password VARCHAR(255) DEFAULT '' AFTER email;`

	_, err := db.Exec(query)
	if err != nil {
		// 如果欄位已存在，MySQL 會報 "Duplicate column name" 錯誤
		if strings.Contains(err.Error(), "Duplicate column") {
			fmt.Println("→ password 欄位已存在，跳過")
			return nil
		}
		return fmt.Errorf("新增 password 欄位失敗: %v", err)
	}

	fmt.Println("✓ 新增 password 欄位成功")
	return nil
}

// Down - 回滾 migration
func (m *AddPasswordToUsers) Down(db *sql.DB) error {
	query := `ALTER TABLE users DROP COLUMN password;`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("刪除 password 欄位失敗: %v", err)
	}

	fmt.Println("✓ 刪除 password 欄位成功")
	return nil
}
