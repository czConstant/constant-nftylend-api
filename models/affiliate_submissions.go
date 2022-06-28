package models

import "github.com/jinzhu/gorm"

type AffiliateSubmissionStatus string

const (
	AffiliateSubmissionStatusApproved  AffiliateSubmissionStatus = "approved"
	AffiliateSubmissionStatusSubmitted AffiliateSubmissionStatus = "submitted"
)

type AffiliateSubmission struct {
	gorm.Model
	Network     Network
	UserID      uint
	User        *User
	Contact     string
	FullName    string
	Website     string `gorm:"type:text collate utf8mb4_unicode_ci"`
	Description string `gorm:"type:text collate utf8mb4_unicode_ci"`
	Status      AffiliateSubmissionStatus
}
