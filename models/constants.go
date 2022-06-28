package models

const (
	EMAIL_BORROWER_OFFER_NEW       = "nftpawn_borrower_offer_new"
	EMAIL_BORROWER_LOAN_REMIND7    = "nftpawn_borrower_loan_remind7"
	EMAIL_BORROWER_LOAN_REMIND3    = "nftpawn_borrower_loan_remind3"
	EMAIL_BORROWER_LOAN_REMIND1    = "nftpawn_borrower_loan_remind1"
	EMAIL_BORROWER_LOAN_STARTED    = "nftpawn_borrower_loan_started"
	EMAIL_BORROWER_LOAN_LIQUIDATED = "nftpawn_borrower_loan_liquidated"
	EMAIL_LENDER_OFFER_STARTED     = "nftpawn_lender_offer_started"
	EMAIL_LENDER_LOAN_REPAID       = "nftpawn_lender_loan_repaid"
	EMAIL_LENDER_LOAN_LIQUIDATED   = "nftpawn_lender_loan_liquidated"
	EMAIL_USER_VERIFY_EMAIL        = "nftpawn_user_verify_email"
	EMAIL_AFFILIATE_SUBMISSION     = "nftpawn_affiliate_submission"
)

type EmailQueue struct {
	EmailType string
	ObjectID  uint
}
