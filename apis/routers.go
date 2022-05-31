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
	}
	nftAPI.POST("/blockchain/update-block/:block", s.NftLendUpdateBlock)
	nftAPI.POST("/blockchain/:network/scan-block/:block", s.BlockchainScanBlock)
	currencynftAPI := nftAPI.Group("/currencies")
	{
		currencynftAPI.GET("/list", s.GetCurrencies)
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
		collectionnftAPI.POST("/submitted", s.CreateCollectionSubmitted)
	}
	loannftAPI := nftAPI.Group("/loans")
	{
		loannftAPI.GET("/borrower-stats/:address", s.GetBorrowerStats)
		loannftAPI.GET("/listing", s.GetListingLoans)
		loannftAPI.GET("/list", s.GetLoans)
		loannftAPI.GET("/offers", s.GetLoanOffers)
		loannftAPI.GET("/transactions", s.GetLoanTransactions)
		loannftAPI.POST("/create", s.CreateLoan)
		loannftAPI.POST("/offers/create/:loan_id", s.CreateLoanOffer)
		loannftAPI.POST("/near/sync", s.NearUpdateLoan)
	}
	hookInternalnftAPI := nftAPI.Group("/hook/internal")
	{
		hookInternalnftAPI.POST("/solana-instruction", s.LenInternalHookSolanaInstruction)
		hookInternalnftAPI.POST("/near-sync", s.NearSync)
	}
	jobsNftAPI := nftAPI.Group("/jobs")
	jobsNftAPI.Use(s.authorizeJobMiddleware())
	{
		jobsNftAPI.POST("/evm-filter-logs", s.JobEvmNftypawnFilterLogs)
		jobsNftAPI.POST("/update-price", s.JobUpdateCurrencyPrice)
		jobsNftAPI.POST("/email-chedule", s.JobEmailSchedule)
	}
	userNftAPI := nftAPI.Group("/users")
	{
		userNftAPI.GET("/settings", s.UserGetSettings)
		userNftAPI.POST("/settings", s.UserUpdateSetting)
	}
	notiAPI := nftAPI.Group("/notifications")
	{
		notiAPI.GET("/list", s.GetNotifications)
		notiAPI.POST("/seen", s.SeenNotification)
	}
}
