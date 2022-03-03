package serializers

type Resp struct {
	Result      interface{} `json:"result"`
	Error       error       `json:"error"`
	Count       *uint       `json:"count,omitempty"`
	CountUnread *uint       `json:"count_unread,omitempty"`
}

// type NftTokenResp struct {
// 	ID        uint             `json:"id"`
// 	MintHash  string           `json:"mint_hash,omitempty"`
// 	IpfsHash  string           `json:"ipfs_hash,omitempty"`
// 	Status    models.NftStatus `json:"status"`
// 	Receiver  string           `json:"receiver,omitempty"`
// 	CreatedAt time.Time        `json:"created_at"`
// 	Type      models.NftType   `json:"type"`
// 	ExtraData string           `json:"extra_data"`
// 	Reward    *NftRewardResp   `json:"reward"`
// }

// func NewNftTokenResp(m *models.NftTokenUser) *NftTokenResp {
// 	if m == nil {
// 		return nil
// 	}
// 	return &NftTokenResp{
// 		ID:        m.ID,
// 		CreatedAt: m.CreatedAt,
// 		MintHash:  m.MintHash,
// 		IpfsHash:  m.IpfsHash,
// 		Status:    m.Status,
// 		Receiver:  m.ReceiverAddress,
// 		Type:      m.Type,
// 		ExtraData: m.ExtraData,
// 		Reward:    NewNftRewardResp(m.Reward),
// 	}
// }

// func NewNftTokenRespArr(arr []*models.NftTokenUser) []*NftTokenResp {
// 	res := []*NftTokenResp{}
// 	for _, m := range arr {
// 		res = append(res, NewNftTokenResp(m))
// 	}
// 	return res
// }
