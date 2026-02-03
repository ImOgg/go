package bootstrap

import (
	"my-api/app/pkg/logger"
	"my-api/config"
)

// Log 全域 Logger 實例
var Log *logger.Logger

// InitLogger 初始化 Logger
func InitLogger() {
	cfg := config.GlobalConfig.Log

	// 建立 Logger
	Log = logger.New(cfg)

	// 設定為全域 Logger
	logger.SetGlobal(Log)

	Log.Info("Logger 初始化成功", map[string]interface{}{
		"level":  cfg.Level,
		"format": cfg.Format,
		"output": cfg.Output,
	})
}
