package models

import "gorm.io/gorm"

// User 使用者模型
//
// Struct Tag 說明（前綴區分不同套件使用）：
//
// 【json: 套件】控制 API JSON 輸出（類似 Laravel 的 $hidden / $visible）
//   - json:"name"      → 輸出為 "name" 欄位
//   - json:"-"         → 隱藏此欄位，不輸出
//   - json:",omitempty" → 若為空值則不輸出
//
// 【gorm: 套件】控制資料庫設定（類似 Laravel 的 Migration）
//   - gorm:"type:varchar(100)"  → 欄位類型 VARCHAR(100)
//   - gorm:"not null"           → 不可為空（NOT NULL）
//   - gorm:"uniqueIndex"        → 唯一索引（UNIQUE INDEX）
//   - gorm:"default:'value'"    → 預設值
//   - gorm:"primaryKey"         → 主鍵
//
// 【其他常見 Tag】
//   - validate:"required"  → 驗證必填（go-validator 套件）
//   - binding:"required"   → Gin 綁定驗證
type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(100);not null"`         // 輸出為 "name"
	Email    string `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"` // 輸出為 "email"
	Password string `json:"-" gorm:"type:varchar(255)"`                     // "-" 表示隱藏，不輸出到 JSON
	Age      int    `json:"age"`                                            // 輸出為 "age"
}
