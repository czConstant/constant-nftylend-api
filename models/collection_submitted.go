package models

import "github.com/jinzhu/gorm"

type CollectionSubmittedStatus string

type CollectionSubmitted struct {
	gorm.Model
	Network         Network
	Name            string
	Description     string
	Creator         string
	ContractAddress string
	ContactInfo     string
	Verified        bool `gorm:"default:0"`
	WhoVerified     string
	Status          CollectionSubmittedStatus
}
