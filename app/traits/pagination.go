package traits

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Pagination - 分頁資料結構
type Pagination struct {
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

// PaginationParams - 分頁參數
type PaginationParams struct {
	Page    int
	PerPage int
}

// GetPaginationParams - 從請求中取得分頁參數
func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	return PaginationParams{
		Page:    page,
		PerPage: perPage,
	}
}

// Paginate - GORM 分頁查詢
func Paginate(db *gorm.DB, page, perPage int, data interface{}) (*Pagination, error) {
	var total int64
	
	// 計算總數
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 計算 offset
	offset := (page - 1) * perPage

	// 查詢資料
	if err := db.Offset(offset).Limit(perPage).Find(data).Error; err != nil {
		return nil, err
	}

	// 計算總頁數
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &Pagination{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
		Data:       data,
	}, nil
}
