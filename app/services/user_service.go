package services

import (
	"errors"
	"my-api/app/models"
	"my-api/app/repositories"
	"my-api/app/requests"
	"my-api/app/responses"
)

// UserService - 使用者業務邏輯層介面
type UserService interface {
	CreateUser(req *requests.CreateUserRequest) (*responses.UserResponse, error)
	GetAllUsers() ([]responses.UserResponse, error)
	GetUserByID(id uint) (*responses.UserResponse, error)
	UpdateUser(id uint, req *requests.UpdateUserRequest) (*responses.UserResponse, error)
	DeleteUser(id uint) error
}

// userService - 實作 UserService 介面
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService - 建立新的 UserService 實例
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser - 新增使用者（含業務邏輯）
func (s *userService) CreateUser(req *requests.CreateUserRequest) (*responses.UserResponse, error) {
	// 檢查 Email 是否已存在
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("電子郵件已被使用")
	}

	// 建立使用者
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 轉換為 Response DTO
	return responses.NewUserResponse(user), nil
}

// GetAllUsers - 取得所有使用者
func (s *userService) GetAllUsers() ([]responses.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// 轉換為 Response DTO 列表
	var userResponses []responses.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, *responses.NewUserResponse(&user))
	}

	return userResponses, nil
}

// GetUserByID - 根據 ID 取得使用者
func (s *userService) GetUserByID(id uint) (*responses.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("使用者不存在")
	}

	return responses.NewUserResponse(user), nil
}

// UpdateUser - 更新使用者
func (s *userService) UpdateUser(id uint, req *requests.UpdateUserRequest) (*responses.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("使用者不存在")
	}

	// 只更新有提供的欄位
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// 檢查新 Email 是否已被其他人使用
		existingUser, _ := s.userRepo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("電子郵件已被使用")
		}
		user.Email = req.Email
	}
	if req.Age > 0 {
		user.Age = req.Age
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return responses.NewUserResponse(user), nil
}

// DeleteUser - 刪除使用者
func (s *userService) DeleteUser(id uint) error {
	// 檢查使用者是否存在
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("使用者不存在")
	}

	return s.userRepo.Delete(id)
}
