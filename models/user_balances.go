package models

import (
	"math/big"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type UserBalance struct {
	gorm.Model
	Network        Network
	UserID         uint
	User           *User
	CurrencyID     uint
	Currency       *Currency
	Balance        numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	LockedBalance  numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	ClaimedBalance numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
}

func (m *UserBalance) GetAvaiableBalance() *big.Float {
	return SubBigFloats(&m.Balance.Float, &m.LockedBalance.Float)
}
