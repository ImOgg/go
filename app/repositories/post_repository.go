package repositories

import (
	"gorm.io/gorm"
	"my-api/app/models"
)

// PostRepository - 文章資料存取層介面
type PostRepository interface {
	Create(post *models.Post) error
	FindAll() ([]models.Post, error)
	FindByID(id uint) (*models.Post, error)
	Update(post *models.Post) error
	Delete(id uint) error
	FindByUserID(userID uint) ([]models.Post, error)
}

// postRepository - 實作 PostRepository 介面
type postRepository struct {
	db *gorm.DB
}

// NewPostRepository - 建立新的 PostRepository 實例
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

// Create - 新增文章
func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

// FindAll - 取得所有文章
func (r *postRepository) FindAll() ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Preload("User").Find(&posts).Error
	return posts, err
}

// FindByID - 根據 ID 查詢文章
func (r *postRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("User").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// Update - 更新文章
func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

// Delete - 刪除文章（軟刪除）
func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

// FindByUserID - 根據使用者 ID 查詢文章
func (r *postRepository) FindByUserID(userID uint) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.Where("user_id = ?", userID).Preload("User").Find(&posts).Error
	return posts, err
}
