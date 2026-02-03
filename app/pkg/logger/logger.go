package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"my-api/config"
)

// Logger 封裝 zerolog.Logger，提供 Laravel 風格的 API
type Logger struct {
	zerolog.Logger
}

// 全域 Logger 實例
var globalLogger *Logger

// New 建立新的 Logger
func New(cfg config.LogConfig) *Logger {
	// 設定日誌等級
	level := parseLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	// 設定時間格式
	zerolog.TimeFieldFormat = time.RFC3339

	// 建立 writer
	var writer io.Writer

	// Console writer (彩色輸出)
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
		NoColor:    false,
	}

	// File writer (使用 lumberjack 做日誌輪替)
	// 確保目錄存在
	logDir := filepath.Dir(cfg.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// 如果無法建立目錄，使用 stdout
		writer = consoleWriter
	}

	fileWriter := &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}

	// 根據環境和配置選擇輸出
	switch cfg.Output {
	case "stdout":
		if cfg.Format == "json" {
			writer = os.Stdout
		} else {
			writer = consoleWriter
		}
	case "file":
		writer = fileWriter
	case "both":
		if cfg.Format == "json" {
			writer = io.MultiWriter(os.Stdout, fileWriter)
		} else {
			writer = io.MultiWriter(consoleWriter, fileWriter)
		}
	default:
		writer = consoleWriter
	}

	logger := zerolog.New(writer).With().
		Timestamp().
		Logger()

	return &Logger{Logger: logger}
}

// parseLevel 解析日誌等級字串
func parseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

// SetGlobal 設定全域 Logger
func SetGlobal(l *Logger) {
	globalLogger = l
}

// Global 取得全域 Logger
func Global() *Logger {
	return globalLogger
}

// ===== Laravel 風格的便捷方法 =====

// Debug 記錄 debug 等級日誌
func (l *Logger) Debug(message string, context ...map[string]interface{}) {
	event := l.Logger.Debug()
	if len(context) > 0 {
		event = event.Fields(context[0])
	}
	event.Msg(message)
}

// Info 記錄 info 等級日誌
func (l *Logger) Info(message string, context ...map[string]interface{}) {
	event := l.Logger.Info()
	if len(context) > 0 {
		event = event.Fields(context[0])
	}
	event.Msg(message)
}

// Warning 記錄 warning 等級日誌
func (l *Logger) Warning(message string, context ...map[string]interface{}) {
	event := l.Logger.Warn()
	if len(context) > 0 {
		event = event.Fields(context[0])
	}
	event.Msg(message)
}

// Error 記錄 error 等級日誌
func (l *Logger) Error(message string, context ...map[string]interface{}) {
	event := l.Logger.Error()
	if len(context) > 0 {
		event = event.Fields(context[0])
	}
	event.Msg(message)
}

// Fatal 記錄 fatal 等級日誌並結束程式
func (l *Logger) Fatal(message string, context ...map[string]interface{}) {
	event := l.Logger.Fatal()
	if len(context) > 0 {
		event = event.Fields(context[0])
	}
	event.Msg(message)
}

// WithContext 建立帶有預設欄位的子 Logger
func (l *Logger) WithContext(ctx map[string]interface{}) *Logger {
	return &Logger{
		Logger: l.Logger.With().Fields(ctx).Logger(),
	}
}

// WithRequestID 建立帶有 Request ID 的子 Logger
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		Logger: l.Logger.With().Str("request_id", requestID).Logger(),
	}
}

// WithError 建立帶有錯誤資訊的子 Logger
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		Logger: l.Logger.With().Err(err).Logger(),
	}
}
