package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type VerificationType string
type VerificationStatus string

const (
	VerificationTypeEmail VerificationType = "email"

	VerificationStatusVerifying VerificationStatus = "verifying"
	VerificationStatusVerified  VerificationStatus = "verified"
	VerificationStatusCancelled VerificationStatus = "cancelled"
)

type Verification struct {
	gorm.Model
	Network   Network
	UserID    uint
	Type      VerificationType
	Email     string
	Token     string
	ExpiredAt *time.Time
	Status    VerificationStatus
}
