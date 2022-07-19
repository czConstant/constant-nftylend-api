package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Leaderboard struct {
	gorm.Model
	Network  Network
	RptDate  *time.Time
	Rewards  string `gorm:"type:text collate utf8mb4_unicode_ci"`
	ImageURL string `gorm:"type:text collate utf8mb4_unicode_ci"`
}
