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
	Title       string `gorm:"type:text"`
	Content     string `gorm:"type:text"`
	RedirectURL string `gorm:"type:text"`
	ImageURL    string `gorm:"type:text"`
}
