package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
)

type NotificationResp struct {
	ID          uint                    `json:"id"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
	Network     models.Network          `json:"network"`
	Type        models.NotificationType `json:"type"`
	Address     string                  `json:"address"`
	Title       string                  `json:"title"`
	Content     string                  `json:"content"`
	RedirectURL string                  `json:"redirect_url"`
	ImageURL    string                  `json:"image_url"`
}

func NewNotificationResp(m *models.Notification) *NotificationResp {
	if m == nil {
		return nil
	}
	resp := &NotificationResp{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Network:     m.Network,
		Type:        m.Type,
		Address:     m.Address,
		Title:       m.Title,
		Content:     m.Content,
		RedirectURL: m.RedirectURL,
		ImageURL:    m.ImageURL,
	}
	return resp
}

func NewNotificationRespArr(arr []*models.Notification) []*NotificationResp {
	resps := []*NotificationResp{}
	for _, m := range arr {
		resps = append(resps, NewNotificationResp(m))
	}
	return resps
}
