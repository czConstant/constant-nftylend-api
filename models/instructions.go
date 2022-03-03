package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Instruction struct {
	gorm.Model
	BlockNumber      uint64
	BlockTime        *time.Time
	TransactionHash  string
	TransactionIndex uint
	InstructionIndex uint
	Program          string
	Instruction      string
	Data             string `gorm:"type:text"`
	Status           string `gorm:"default:0"`
}
