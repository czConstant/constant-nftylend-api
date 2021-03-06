package apis

import (
	"net/http"
	"time"

	"github.com/czConstant/constant-nftylend-api/configs"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/czConstant/constant-nftylend-api/services"
	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	g    *gin.Engine
	conf *configs.Config
	nls  *services.NftLend
}

func NewServer(
	g *gin.Engine,
	conf *configs.Config,
	nls *services.NftLend,
) *Server {
	return &Server{
		g:    g,
		conf: conf,
		nls:  nls,
	}
}

func (s *Server) Routers() {
	s.g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://*", "https://*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"*"},
		MaxAge:           12 * time.Hour,
	}))
	s.g.Use(s.logApiMiddleware())
	s.g.Use(s.recoveryMiddleware(raven.DefaultClient, false))
	nftAPI := s.g.Group("/api")
	{
		nftAPI.GET("/", func(c *gin.Context) {
			ctxJSON(c, http.StatusOK, &serializers.Resp{Error: nil})
		})
		nftAPI.GET("/configs", s.AppConfigs)
		nftAPI.GET("/moralis/:address/nft", s.MoralisGetNFTs)
		nftAPI.GET("/kitwallet/account/:address/likelyNFTs", s.GetKitwalletAccountNfts)
		nftAPI.GET("/ipfs/:hash", s.GetIpfsInfo)
	}
	nftAPI.POST("/blockchain/update-block/:block", s.NftLendUpdateBlock)
	nftAPI.POST("/blockchain/:network/scan-block/:block", s.BlockchainScanBlock)
	currencynftAPI := nftAPI.Group("/currencies")
	{
		currencynftAPI.GET("/list", s.GetCurrencies)
		currencynftAPI.GET("/pwp-token", s.GetCurrencyPWPToken)
	}
	assetnftAPI := nftAPI.Group("/assets")
	{
		assetnftAPI.GET("/detail/:seo_url", s.GetAssetDetail)
		assetnftAPI.GET("/info", s.GetAssetDetailInfo)
		assetnftAPI.GET("/transactions", s.GetAseetTransactions)
	}
	collectionnftAPI := nftAPI.Group("/collections")
	{
		collectionnftAPI.GET("/list", s.GetCollections)
		collectionnftAPI.GET("/detail/:seo_url", s.GetCollectionDetail)
		collectionnftAPI.GET("/verified", s.GetCollectionAssetVerified)
		collectionnftAPI.POST("/submitted", s.recaptchaV3Middleware(), s.CreateCollectionSubmission)
		collectionnftAPI.GET("/near-whitelist-creators", s.GetNearApprovedCreators)
		collectionnftAPI.GET("/near-whitelist-collections", s.GetNearApprovedCollections)
	}
	loannftAPI := nftAPI.Group("/loans")
	{
		loannftAPI.GET("/borrower-stats", s.GetBorrowerStats)
		loannftAPI.GET("/lender-stats", s.GetLenderStats)
		loannftAPI.GET("/platform-stats", s.GetPlatformStats)
		loannftAPI.GET("/listing", s.GetListingLoans)
		loannftAPI.GET("/list", s.GetLoans)
		loannftAPI.GET("/offers", s.GetLoanOffers)
		loannftAPI.GET("/transactions", s.GetLoanTransactions)
		loannftAPI.POST("/create", s.CreateLoan)
		loannftAPI.POST("/offers/create/:loan_id", s.CreateLoanOffer)
		loannftAPI.POST("/near/sync", s.NearUpdateLoan)
		loannftAPI.GET("/leaderboard", s.GetLeaderBoardAtNow)
	}
	proposalAPI := nftAPI.Group("/proposals")
	{
		proposalAPI.GET("/list", s.GetProposals)
		proposalAPI.POST("/create", s.recaptchaV3Middleware(), s.CreateProposal)
		proposalAPI.GET("/detail/:proposal_id", s.GetProposalDetail)
		proposalAPI.GET("/votes/list/:proposal_id", s.GetProposalVotes)
		proposalAPI.GET("/votes/vote/:proposal_id", s.GetUserProposalVote)
		proposalAPI.POST("/votes/create", s.recaptchaV3Middleware(), s.CreateProposalVote)
	}
	hookInternalnftAPI := nftAPI.Group("/hook/internal")
	{
		hookInternalnftAPI.POST("/solana-instruction", s.LenInternalHookSolanaInstruction)
		hookInternalnftAPI.POST("/near-sync", s.NearSync)
		hookInternalnftAPI.POST("/near-pwp-sync", s.NearPwpSync)
		hookInternalnftAPI.POST("/near-nft-sync", s.NearNftTransfer)
	}
	jobsNftAPI := nftAPI.Group("/jobs")
	jobsNftAPI.Use(s.authorizeJobMiddleware())
	{
		jobsNftAPI.POST("/evm-filter-logs", s.JobEvmNftypawnFilterLogs)
		jobsNftAPI.POST("/incentive-unlock", s.JobIncentiveForUnlock)
		jobsNftAPI.POST("/update-price", s.JobUpdateCurrencyPrice)
		jobsNftAPI.POST("/email-chedule", s.JobEmailSchedule)
		jobsNftAPI.POST("/proposal-status", s.JobProposalStatus)
		jobsNftAPI.POST("/update-stats", s.JobUpdateStats)
	}
	userNftAPI := nftAPI.Group("/users")
	{
		userNftAPI.GET("/settings", s.UserGetSettings)
		userNftAPI.POST("/settings", s.recaptchaV3Middleware(), s.UserUpdateSetting)
		userNftAPI.POST("/connected", s.UserConnected)
		userNftAPI.GET("/stats", s.GetUserStats)
		userNftAPI.GET("/balances/pwp", s.GetUserPWPCurrencyBalance)
		userNftAPI.GET("/balances/near", s.GetUserNEARCurrencyBalance)
		userNftAPI.GET("/balances/transactions", s.GetUserBalanceTransactions)
		userNftAPI.POST("/balances/claim", s.recaptchaV3Middleware(), s.ClaimUserBalance)
	}
	notiAPI := nftAPI.Group("/notifications")
	{
		notiAPI.GET("/list", s.GetNotifications)
		notiAPI.POST("/seen", s.SeenNotification)
	}
	verificationAPI := nftAPI.Group("/verifications")
	{
		verificationAPI.POST("/verify-email", s.recaptchaV3Middleware(), s.UserVerifyEmail)
		verificationAPI.POST("/verify-token", s.UserVerifyEmailToken)
	}
	affiliateAPI := nftAPI.Group("/affiliates")
	{
		affiliateAPI.GET("/stats", s.GetAffiliateStats)
		affiliateAPI.GET("/volumes", s.GetAffiliateVolumes)
		affiliateAPI.GET("/transactions", s.GetAffiliateTransactions)
		affiliateAPI.POST("/submitted", s.recaptchaV3Middleware(), s.CreateAffiliateSubmitted)
	}
}
