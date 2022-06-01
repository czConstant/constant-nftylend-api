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
	Name              string `gorm:"type:text"`
	Description       string `gorm:"type:text"`
	Status            string
	LoanValidDuration uint `gorm:"default:0"`
}
