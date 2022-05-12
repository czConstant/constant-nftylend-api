package models

const (
	EMAIL_BORROWER_OFFER_NEW       = "nftpawn_borrower_offer_new"
	EMAIL_BORROWER_REMIND_PAYBACK  = "nftpawn_borrower_loan_remind"
	EMAIL_BORROWER_LOAN_STARTED    = "nftpawn_borrower_loan_started"
	EMAIL_BORROWER_LOAN_LIQUIDATED = "nftpawn_borrower_loan_liquidated"
	EMAIL_LENDER_OFFER_STARTED     = "nftpawn_lender_offer_started"
	EMAIL_LENDER_LOAN_REPAID       = "nftpawn_lender_loan_repaid"
	EMAIL_LENDER_LOAN_LIQUIDATED   = "nftpawn_lender_loan_liquidated"
)

type EmailQueue struct {
	EmailType string
	ObjectID  uint
}
