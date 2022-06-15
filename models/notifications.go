package models

import (
	"github.com/jinzhu/gorm"
)

type Notification struct {
	gorm.Model
	Network     Network
	UserID      uint
	User        *User
	Type        NotificationType
	Title       string `gorm:"type:text collate utf8mb4_unicode_ci"`
	Content     string `gorm:"type:text collate utf8mb4_unicode_ci"`
	RedirectURL string `gorm:"type:text collate utf8mb4_unicode_ci"`
	ImageURL    string `gorm:"type:text collate utf8mb4_unicode_ci"`
}
