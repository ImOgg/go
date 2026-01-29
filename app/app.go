package app

import (
	"gorm.io/gorm"
	"my-api/app/repositories"
	"my-api/app/services"
)

// App - 應用程式容器（類似 Laravel 的 Service Container）
type App struct {
	DB *gorm.DB

	// Repositories
	UserRepository repositories.UserRepository

	// Services
	UserService services.UserService
	AuthService services.AuthService
}

// NewApp - 建立新的應用程式容器
func NewApp(db *gorm.DB) *App {
	app := &App{
		DB: db,
	}

	// 初始化 Repositories
	app.UserRepository = repositories.NewUserRepository(db)

	// 初始化 Services（注入 Repository 依賴）
	app.UserService = services.NewUserService(app.UserRepository)
	app.AuthService = services.NewAuthService(app.UserRepository)

	return app
}
