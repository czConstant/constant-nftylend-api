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
	nftAPI := s.g.Group("/nfty-lend-api")
	{
		nftAPI.GET("/", func(c *gin.Context) {
			ctxJSON(c, http.StatusOK, &serializers.Resp{Error: nil})
		})
	}
	nftAPI.POST("/blockchain/update-block/:block", s.NftLendUpdateBlock)
	currencynftAPI := nftAPI.Group("/currencies")
	{
		currencynftAPI.GET("/list", s.GetCurrencies)
	}
	assetnftAPI := nftAPI.Group("/assets")
	{
		assetnftAPI.GET("/detail/:seo_url", s.GetAssetDetail)
		assetnftAPI.GET("/transactions", s.GetAseetTransactions)
	}
	collectionnftAPI := nftAPI.Group("/collections")
	{
		collectionnftAPI.GET("/list", s.GetCollections)
		collectionnftAPI.GET("/detail/:seo_url", s.GetCollectionDetail)
		collectionnftAPI.GET("/verified", s.GetCollectionAssetVerified)
	}
	loannftAPI := nftAPI.Group("/loans")
	{
		loannftAPI.GET("/listing", s.GetListingLoans)
		loannftAPI.GET("/list", s.GetLoans)
		loannftAPI.GET("/offers", s.GetLoanOffers)
		loannftAPI.GET("/transactions", s.GetLoanTransactions)
	}
	hookInternalnftAPI := nftAPI.Group("/hook/internal")
	{
		hookInternalnftAPI.POST("/solana-instruction", s.LenInternalHookSolanaInstruction)
	}
}
