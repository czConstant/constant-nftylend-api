package models

import (
	"github.com/jinzhu/gorm"
)

type Notification struct {
	gorm.Model
	Type        NotificationType
	Address     string
	Title       string `gorm:"type:text"`
	Content     string `gorm:"type:text"`
	RedirectURL string `gorm:"type:text"`
}
