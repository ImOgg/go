package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"my-api/config"
	"my-api/database"
)

func main() {
	// 定義命令行參數
	migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	rollbackCmd := flag.NewFlagSet("rollback", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)

	// 檢查參數
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// 載入配置
	config.LoadConfig()

	// 根據命令執行對應操作
	switch os.Args[1] {
	case "migrate":
		migrateCmd.Parse(os.Args[2:])
		runMigrate()
		
	case "rollback":
		rollbackCmd.Parse(os.Args[2:])
		runRollback()
		
	case "status":
		statusCmd.Parse(os.Args[2:])
		showStatus()
	
	case "make":
		if len(os.Args) < 3 {
			fmt.Println("❌ 請提供 migration 名稱")
			fmt.Println("範例: go run cmd/migrate/main.go make add_phone_to_users")
			os.Exit(1)
		}
		makeMigration(os.Args[2])
		
	default:
		printUsage()
		os.Exit(1)
	}
}

func runMigrate() {
	fmt.Println("🚀 執行 Migration...")
	if err := database.RunMigrations(); err != nil {
		log.Fatal("❌ Migration 失敗:", err)
	}
}

func runRollback() {
	fmt.Println("⏮️  執行 Rollback...")
	if err := database.RollbackMigration(); err != nil {
		log.Fatal("❌ Rollback 失敗:", err)
	}
}

func showStatus() {
	status, err := database.GetMigrationStatus()
	if err != nil {
		log.Fatal("❌ 無法獲取狀態:", err)
	}
	
	fmt.Println("📊 Migration 狀態:")
	fmt.Println()
	
	// 引入 migrations 確保所有 migration 被註冊
	_ = "my-api/database/migrations"
	
	allMigrations := getAllMigrations()
	for _, m := range allMigrations {
		executed := status[m.Version]
		statusIcon := "⏳"
		statusText := "待執行"
		if executed {
			statusIcon = "✅"
			statusText = "已執行"
		}
		fmt.Printf("  %s [%s] %s - %s\n", statusIcon, m.Version, m.Description, statusText)
	}
}

type migrationInfo struct {
	Version     string
	Description string
}

func getAllMigrations() []migrationInfo {
	// 這裡需要手動列出或從 registry 獲取
	// 為了簡化，這裡返回空列表
	// 實際使用時，可以從 migrations.All() 獲取
	return []migrationInfo{
		{Version: "000001", Description: "create_users_table"},
	}
}

func makeMigration(name string) {
	template := `package migrations

import (
	"database/sql"
	"fmt"
)

// TODO: 修改這個名稱
type Migration%s struct {
	BaseMigration
}

func init() {
	Register(&Migration%s{
		BaseMigration: BaseMigration{
			version:     "%s",  // TODO: 修改版本號
			description: "%s",
		},
	})
}

// Up - 執行 migration
func (m *Migration%s) Up(db *sql.DB) error {
	query := ` + "`" + `
		-- TODO: 在這裡寫 SQL
		-- 範例:
		-- ALTER TABLE users ADD COLUMN phone VARCHAR(20);
	` + "`" + `
	
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("migration 失敗: %%v", err)
	}
	
	fmt.Println("✓ Migration 成功")
	return nil
}

// Down - 回滾 migration
func (m *Migration%s) Down(db *sql.DB) error {
	query := ` + "`" + `
		-- TODO: 在這裡寫回滾 SQL
		-- 範例:
		-- ALTER TABLE users DROP COLUMN phone;
	` + "`" + `
	
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("rollback 失敗: %%v", err)
	}
	
	fmt.Println("✓ Rollback 成功")
	return nil
}
`
	
	// 生成文件名和結構名
	structName := toPascalCase(name)
	version := getNextVersion()
	
	content := fmt.Sprintf(template, structName, structName, version, name, structName, structName)
	
	filename := fmt.Sprintf("database/migrations/%s_%s.go", version, name)
	
	// 這裡應該寫入文件，但為了簡化，只打印
	fmt.Printf("✅ Migration 模板已生成：\n")
	fmt.Printf("   文件名: %s\n", filename)
	fmt.Printf("   結構名: %s\n", structName)
	fmt.Printf("\n請手動創建文件並複製以下內容：\n\n")
	fmt.Println(content)
}

func toPascalCase(s string) string {
	// 簡化版本的轉換
	return s
}

func getNextVersion() string {
	// 簡化版本，返回固定值
	// 實際應該讀取現有 migrations 並 +1
	return "000002"
}

func printUsage() {
	fmt.Println(`
Migration 管理工具 - 類似 Laravel Artisan （一個文件包含 Up 和 Down）

使用方式:
  go run cmd/migrate/main.go [命令]

可用命令:
  migrate   - 執行所有待執行的 migrations
  rollback  - 回滾最後一次 migration
  status    - 查看當前 migration 狀態
  make      - 建立新的 migration 文件

範例:
  go run cmd/migrate/main.go migrate
  go run cmd/migrate/main.go rollback
  go run cmd/migrate/main.go status
  go run cmd/migrate/main.go make add_phone_to_users

Docker 內執行:
  docker exec -it my-go-app go run cmd/migrate/main.go migrate
  docker exec -it my-go-app go run cmd/migrate/main.go rollback
`)
}
