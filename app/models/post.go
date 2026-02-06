package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title       string `json:"title" gorm:"type:varchar(255);not null"`
	Content     string `json:"content" gorm:"type:longtext;not null"`
	Description string `json:"description" gorm:"type:varchar(500)"`
	UserID      uint   `json:"user_id" gorm:"not null;index"`
	User        User   `json:"user" gorm:"foreignKey:UserID"`
}
