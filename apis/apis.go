package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) AppConfigs(c *gin.Context) {
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: gin.H{
		"program_id":               s.conf.Contract.ProgramID,
		"matic_nftypawn_address":   s.conf.Contract.MaticNftypawnAddress,
		"matic_nftypawn_admin_fee": 100,
		"avax_nftypawn_address":    s.conf.Contract.AvaxNftypawnAddress,
		"avax_nftypawn_admin_fee":  100,
		"bsc_nftypawn_address":     s.conf.Contract.BscNftypawnAddress,
		"bsc_nftypawn_admin_fee":   100,
		"boba_nftypawn_address":    s.conf.Contract.BobaNftypawnAddress,
		"boba_nftypawn_admin_fee":  100,
		"near_nftypawn_address":    s.conf.Contract.NearNftypawnAddress,
		"near_nftypawn_admin_fee":  100,
		"proposals":                s.conf.Proposals,
	}})
}

func (s *Server) MoralisGetNFTs(c *gin.Context) {
	limit, _ := s.uintFromContextQuery(c, "limit")
	rs, err := s.nls.MoralisGetNFTs(s.requestContext(c), models.MoralisNetworkMap[s.conf.Env][s.stringFromContextQuery(c, "network")], s.stringFromContextParam(c, "address"), s.stringFromContextQuery(c, "cursor"), int(limit))
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, rs)
}

func (s *Server) GetCurrencies(c *gin.Context) {
	ctx := s.requestContext(c)
	currencies, err := s.nls.GetCurrencies(ctx, models.Network(s.stringFromContextQuery(c, "network")))
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewCurrencyRespArr(currencies)})
}

func (s *Server) GetCurrencyPWPToken(c *gin.Context) {
	ctx := s.requestContext(c)
	m, err := s.nls.GetCurrencyBySymbol(ctx, models.SymbolPWPToken, models.NetworkNEAR)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewCurrencyResp(m)})
}
