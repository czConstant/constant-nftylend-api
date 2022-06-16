package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type IncentiveProgramStatus string

const (
	IncentiveProgramStatusActived   IncentiveProgramStatus = "actived"
	IncentiveProgramStatusDeactived IncentiveProgramStatus = "deactived"
)

type IncentiveProgram struct {
	gorm.Model
	Network           Network
	CurrencyID        uint
	Currency          *Currency
	Start             *time.Time
	End               *time.Time
	Name              string `gorm:"type:text collate utf8mb4_unicode_ci"`
	Description       string `gorm:"type:text collate utf8mb4_unicode_ci"`
	Status            string
	LoanValidDuration uint `gorm:"default:0"`
	LockDuration      uint `gorm:"default:0"`
}
