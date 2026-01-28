package repositories

import (
	"gorm.io/gorm"
	"my-api/app/models"
)

// UserRepository - 使用者資料存取層介面
type UserRepository interface {
	Create(user *models.User) error
	FindAll() ([]models.User, error)
	FindByID(id uint) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	FindByEmail(email string) (*models.User, error)
}

// userRepository - 實作 UserRepository 介面
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository - 建立新的 UserRepository 實例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create - 新增使用者
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindAll - 取得所有使用者
func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

// FindByID - 根據 ID 查詢使用者
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update - 更新使用者
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete - 刪除使用者（軟刪除）
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// FindByEmail - 根據 Email 查詢使用者
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
