package models

import (
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/helpers"
	"github.com/jinzhu/gorm"
)

type NotificationType string

type NotificationTemplate struct {
	gorm.Model
	Type        NotificationType
	Title       string `gorm:"type:text"`
	Content     string `gorm:"type:text"`
	RedirectURL string `gorm:"type:text"`
	ImageURL    string `gorm:"type:text"`
	Enabled     bool   `gorm:"default:0"`
}

func (m *NotificationTemplate) Execute(network Network, userID uint, data map[string]interface{}) (*Notification, error) {
	title, err := helpers.GenerateTemplateContent(m.Title, data)
	if err != nil {
		return nil, errs.NewError(err)
	}
	content, err := helpers.GenerateTemplateContent(m.Content, data)
	if err != nil {
		return nil, errs.NewError(err)
	}
	redirectURL, err := helpers.GenerateTemplateContent(m.RedirectURL, data)
	if err != nil {
		return nil, errs.NewError(err)
	}
	imageURL, err := helpers.GenerateTemplateContent(m.ImageURL, data)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return &Notification{
		Network:     network,
		UserID:      userID,
		Type:        m.Type,
		Title:       title,
		Content:     content,
		RedirectURL: redirectURL,
		ImageURL:    imageURL,
	}, nil
}
