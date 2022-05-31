package serializers

import "github.com/czConstant/constant-nftylend-api/models"

type SeenNotificationReq struct {
	Address    string         `json:"address"`
	Network    models.Network `json:"network"`
	SeenNotiID uint           `json:"seen_noti_id"`
}
