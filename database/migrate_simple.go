package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"my-api/config"
	"my-api/database/migrations"
)

// RunMigrations 執行所有待執行的 migrations
func RunMigrations() error {
	db, err := getDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	
	// 建立 migrations 記錄表
	if err := createMigrationsTable(db); err != nil {
		return err
	}
	
	// 獲取已執行的 migrations
	executed, err := getExecutedMigrations(db)
	if err != nil {
		return err
	}
	
	// 執行所有未執行的 migrations
	allMigrations := migrations.All()
	hasNew := false
	
	for _, m := range allMigrations {
		if _, exists := executed[m.Version()]; !exists {
			hasNew = true
			log.Printf("🚀 執行 Migration: %s - %s", m.Version(), m.Description())
			
			if err := m.Up(db); err != nil {
				return fmt.Errorf("migration %s 失敗: %v", m.Version(), err)
			}
			
			// 記錄已執行
			if err := recordMigration(db, m.Version(), m.Description()); err != nil {
				return err
			}
		}
	}
	
	if !hasNew {
		log.Println("✓ 資料庫已是最新版本，無需 migration")
	} else {
		log.Println("✅ 所有 Migrations 執行成功！")
	}
	
	return nil
}

// RollbackMigration 回滾最後一個 migration
func RollbackMigration() error {
	db, err := getDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	
	// 獲取最後執行的 migration
	lastVersion, err := getLastMigration(db)
	if err != nil {
		return err
	}
	
	if lastVersion == "" {
		log.Println("⚠️  沒有可以回滾的 migration")
		return nil
	}
	
	// 獲取對應的 migration
	m, exists := migrations.Get(lastVersion)
	if !exists {
		return fmt.Errorf("找不到版本 %s 的 migration", lastVersion)
	}
	
	log.Printf("⏮️  回滾 Migration: %s - %s", m.Version(), m.Description())
	
	// 執行 Down
	if err := m.Down(db); err != nil {
		return fmt.Errorf("rollback %s 失敗: %v", m.Version(), err)
	}
	
	// 刪除記錄
	if err := removeMigrationRecord(db, lastVersion); err != nil {
		return err
	}
	
	log.Println("✅ Rollback 成功！")
	return nil
}

// GetMigrationStatus 獲取 migration 狀態
func GetMigrationStatus() (map[string]bool, error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	
	if err := createMigrationsTable(db); err != nil {
		return nil, err
	}
	
	executed, err := getExecutedMigrations(db)
	if err != nil {
		return nil, err
	}
	
	status := make(map[string]bool)
	for _, m := range migrations.All() {
		_, status[m.Version()] = executed[m.Version()]
	}
	
	return status, nil
}

// === 輔助函數 ===

func getDBConnection() (*sql.DB, error) {
	cfg := config.GlobalConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&multiStatements=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("無法連接資料庫: %v", err)
	}
	
	return db, nil
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			version VARCHAR(14) PRIMARY KEY,
			description VARCHAR(255) NOT NULL,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	
	_, err := db.Exec(query)
	return err
}

func getExecutedMigrations(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query("SELECT version, description FROM migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	executed := make(map[string]string)
	for rows.Next() {
		var version, description string
		if err := rows.Scan(&version, &description); err != nil {
			return nil, err
		}
		executed[version] = description
	}
	
	return executed, nil
}

func getLastMigration(db *sql.DB) (string, error) {
	var version string
	err := db.QueryRow("SELECT version FROM migrations ORDER BY version DESC LIMIT 1").Scan(&version)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return version, err
}

func recordMigration(db *sql.DB, version, description string) error {
	_, err := db.Exec("INSERT INTO migrations (version, description) VALUES (?, ?)", version, description)
	return err
}

func removeMigrationRecord(db *sql.DB, version string) error {
	_, err := db.Exec("DELETE FROM migrations WHERE version = ?", version)
	return err
}
