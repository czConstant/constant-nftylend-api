package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type IncentiveProgramType string
type IncentiveProgramStatus string

const (
	IncentiveProgramTypeIncentive IncentiveProgramType = "incentive"
	IncentiveProgramTypeAffiliate IncentiveProgramType = "affiliate"
	IncentiveProgramTypeReferral  IncentiveProgramType = "referral"

	IncentiveProgramStatusActived   IncentiveProgramStatus = "actived"
	IncentiveProgramStatusDeactived IncentiveProgramStatus = "deactived"
)

type IncentiveProgram struct {
	gorm.Model
	Network     Network
	Type        IncentiveProgramType `gorm:"default:'incentive'"`
	CurrencyID  uint
	Currency    *Currency
	Start       *time.Time
	End         *time.Time
	Name        string `gorm:"type:text collate utf8mb4_unicode_ci"`
	Description string `gorm:"type:text collate utf8mb4_unicode_ci"`
	Status      string
}
