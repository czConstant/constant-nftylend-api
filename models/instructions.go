package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Instruction struct {
	gorm.Model
	Network          Network `gorm:"default:'SOL'"`
	BlockNumber      uint64
	BlockTime        *time.Time
	TransactionHash  string
	TransactionIndex uint
	InstructionIndex uint
	Program          string
	Instruction      string
	Data             string `gorm:"type:text collate utf8mb4_unicode_ci"`
	Status           string `gorm:"default:0"`
}
