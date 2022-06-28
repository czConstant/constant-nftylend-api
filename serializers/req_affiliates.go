package serializers

type AffiliateSubmittedReq struct {
	SignatureTimestampReq
	Contact     string `json:"contact"`
	FullName    string `json:"full_name"`
	Website     string `json:"website"`
	Description string `json:"description"`
}
