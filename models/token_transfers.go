package models

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type TokenTransfer struct {
	gorm.Model
	Network       Network
	Address       string
	FromAddress   string
	ToAddress     string
	Amount        numeric.BigFloat
	Hash          string
	LogIndex      int64
	TransferredAt *time.Time
}
