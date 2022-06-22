package models

import "github.com/jinzhu/gorm"

type CollectionSubmittedStatus string

const (
	CollectionSubmittedStatusApproved CollectionSubmittedStatus = "approved"
)

type CollectionSubmitted struct {
	gorm.Model
	Network         Network
	Name            string
	Description     string
	Creator         string
	ContractAddress string
	TokenSeriesID   string
	ContactInfo     string
	Verified        bool `gorm:"default:0"`
	WhoVerified     string
	Status          CollectionSubmittedStatus
}
