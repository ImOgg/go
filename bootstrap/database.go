package bootstrap

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"my-api/config"
)

var DB *gorm.DB

// InitDB 初始化資料庫連接
func InitDB() {
	cfg := config.GlobalConfig.Database
	var dialector gorm.Dialector

	// 根據配置選擇資料庫類型
	switch cfg.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
		dialector = mysql.Open(dsn)

	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)
		dialector = postgres.Open(dsn)

	default:
		log.Fatal("不支援的資料庫類型:", cfg.Type)
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		log.Fatal("無法連接到資料庫:", err)
	}

	fmt.Printf("資料庫連接成功！(類型: %s)\n", cfg.Type)
}
