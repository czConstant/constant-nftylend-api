package apis

import (
	"net/http"
	"time"

	"github.com/czConstant/constant-nftlend-api/configs"
	"github.com/czConstant/constant-nftlend-api/serializers"
	"github.com/czConstant/constant-nftlend-api/services"
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
	}
	nftAPI.POST("/blockchain/update-block/:block", s.LendNftLendUpdateBlock)
	currencynftAPI := nftAPI.Group("/currencies")
	{
		currencynftAPI.GET("/list", s.LendGetCurrencies)
	}
	assetnftAPI := nftAPI.Group("/assets")
	{
		assetnftAPI.GET("/detail/:seo_url", s.LendGetAssetDetail)
	}
	collectionnftAPI := nftAPI.Group("/collections")
	{
		collectionnftAPI.GET("/list", s.LendGetCollections)
		collectionnftAPI.GET("/detail/:seo_url", s.LendGetCollectionDetail)
	}
	loannftAPI := nftAPI.Group("/loans")
	{
		loannftAPI.GET("/listing", s.LendGetListingLoans)
		loannftAPI.GET("/list", s.LendGetLoans)
		loannftAPI.GET("/offers", s.LendGetLoanOffers)
		loannftAPI.GET("/transactions", s.LendGetLoanTransactions)
	}
	hookInternalnftAPI := nftAPI.Group("/hook/internal")
	{
		hookInternalnftAPI.POST("/solana-instruction", s.LenInternalHookSolanaInstruction)
	}
}
