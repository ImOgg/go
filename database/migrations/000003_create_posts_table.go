package migrations

import (
	"database/sql"
	"fmt"
)

// CreatePostsTable - 建立 posts 資料表
type CreatePostsTable struct {
	BaseMigration
}

func init() {
	Register(&CreatePostsTable{
		BaseMigration: BaseMigration{
			version:     "000003",
			description: "create_posts_table",
		},
	})
}

// Up - 執行 migration
func (m *CreatePostsTable) Up(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS posts (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content LONGTEXT NOT NULL,
			description VARCHAR(500),
			user_id BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			INDEX idx_user_id (user_id),
			INDEX idx_deleted_at (deleted_at),
			CONSTRAINT fk_posts_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("建立 posts 表失敗: %v", err)
	}

	fmt.Println("✓ 建立 posts 表成功")
	return nil
}

// Down - 回滾 migration
func (m *CreatePostsTable) Down(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS posts;`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("刪除 posts 表失敗: %v", err)
	}

	fmt.Println("✓ 刪除 posts 表成功")
	return nil
}
