package models

import (
	"github.com/jinzhu/gorm"
)

type Notification struct {
	gorm.Model
	UserID      uint
	Network     Network
	Type        NotificationType
	Address     string
	Title       string `gorm:"type:text"`
	Content     string `gorm:"type:text"`
	RedirectURL string `gorm:"type:text"`
	ImageURL    string `gorm:"type:text"`
}
