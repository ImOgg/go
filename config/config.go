package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

type AppConfig struct {
	Env  string
	Port string
}

type DatabaseConfig struct {
	Type     string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string // For PostgreSQL
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       string
}

// GlobalConfig 全域配置變數
var GlobalConfig *Config

// 載入配置
func LoadConfig() {
	// 載入 .env 檔案
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: 找不到 .env 檔案，使用環境變數")
	}

	GlobalConfig = &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "mysql"),
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "test"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "127.0.0.1"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnv("REDIS_DB", "0"),
		},
	}
}

// 獲取環境變數，如果不存在則返回預設值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
