package requests

// CreatePostRequest - 建立文章請求驗證
type CreatePostRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Content     string `json:"content" binding:"required,min=10"`
	Description string `json:"description" binding:"omitempty,max=500"`
	UserID      uint   `json:"user_id" binding:"required,gt=0"`
}

// UpdatePostRequest - 更新文章請求驗證
type UpdatePostRequest struct {
	Title       string `json:"title" binding:"omitempty,min=3,max=255"`
	Content     string `json:"content" binding:"omitempty,min=10"`
	Description string `json:"description" binding:"omitempty,max=500"`
}
