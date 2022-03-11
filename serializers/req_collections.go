package serializers

import "github.com/czConstant/constant-nftylend-api/models"

type CollectionSubmittedReq struct {
	Network         models.Chain `json:"network"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	Creator         string       `json:"creator"`
	ContractAddress string       `json:"contract_address"`
	ContactInfo     string       `json:"contact_info"`
	Verified        bool         `json:"verified"`
	WhoVerified     string       `json:"who_verified"`
}
