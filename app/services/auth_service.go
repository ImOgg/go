package services

import (
	"errors"

	"my-api/app/models"
	"my-api/app/repositories"
	"my-api/app/requests"
	"my-api/app/responses"
	"my-api/app/utils"
	"my-api/config"
)

// AuthService - 認證業務邏輯層介面
type AuthService interface {
	Register(req *requests.RegisterRequest) (*responses.AuthResponse, error)
	Login(req *requests.LoginRequest) (*responses.AuthResponse, error)
	GetCurrentUser(userID uint) (*responses.UserResponse, error)
}

// authService - 實作 AuthService 介面
type authService struct {
	userRepo repositories.UserRepository
}

// NewAuthService - 建立新的 AuthService 實例
func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

// Register - 使用者註冊
func (s *authService) Register(req *requests.RegisterRequest) (*responses.AuthResponse, error) {
	// 檢查 Email 是否已存在
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("電子郵件已被使用")
	}

	// 加密密碼
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密碼加密失敗")
	}

	// 建立使用者
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Age:      req.Age,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 產生 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("Token 產生失敗")
	}

	expiresIn := config.GlobalConfig.JWT.ExpiryHours * 3600
	return responses.NewAuthResponse(token, expiresIn, user), nil
}

// Login - 使用者登入
func (s *authService) Login(req *requests.LoginRequest) (*responses.AuthResponse, error) {
	// 查詢使用者
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("帳號或密碼錯誤")
	}

	// 驗證密碼
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("帳號或密碼錯誤")
	}

	// 產生 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.New("Token 產生失敗")
	}

	expiresIn := config.GlobalConfig.JWT.ExpiryHours * 3600
	return responses.NewAuthResponse(token, expiresIn, user), nil
}

// GetCurrentUser - 取得當前用戶資訊
func (s *authService) GetCurrentUser(userID uint) (*responses.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("使用者不存在")
	}

	return responses.NewUserResponse(user), nil
}
