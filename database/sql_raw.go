package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"my-api/config"
)

var SqlDB *sql.DB

// 使用原生 database/sql 初始化連接
func InitRawDB() {
	cfg := config.GlobalConfig.Database
	var dsn string
	var driverName string

	switch cfg.Type {
	case "mysql":
		driverName = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	case "postgres":
		driverName = "postgres"
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	default:
		log.Fatal("不支援的資料庫類型:", cfg.Type)
	}

	var err error
	SqlDB, err = sql.Open(driverName, dsn)
	if err != nil {
		log.Fatal("無法開啟資料庫連接:", err)
	}

	// 測試連接
	err = SqlDB.Ping()
	if err != nil {
		log.Fatal("無法連接到資料庫:", err)
	}

	// 設定連接池
	SqlDB.SetMaxOpenConns(25)
	SqlDB.SetMaxIdleConns(5)

	fmt.Printf("原生 SQL 資料庫連接成功！(類型: %s)\n", cfg.Type)
}

// 關閉資料庫連接
func CloseRawDB() {
	if SqlDB != nil {
		SqlDB.Close()
	}
}
