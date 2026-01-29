package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(100);not null"`
	Email    string `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password string `json:"-" gorm:"type:varchar(255)"` // json:"-" 隱藏密碼不輸出
	Age      int    `json:"age"`
}
