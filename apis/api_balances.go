package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetUserPWPCurrencyBalance(c *gin.Context) {
	ctx := s.requestContext(c)
	userBalance, err := s.nls.GetUserCurrencyBalance(
		ctx,
		models.Network(s.stringFromContextQuery(c, "network")),
		s.stringFromContextQuery(c, "address"),
		models.SymbolPWPToken,
	)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewUserBalanceResp(userBalance)})
}

func (s *Server) GetUserNEARCurrencyBalance(c *gin.Context) {
	ctx := s.requestContext(c)
	userBalance, err := s.nls.GetUserCurrencyBalance(
		ctx,
		models.Network(s.stringFromContextQuery(c, "network")),
		s.stringFromContextQuery(c, "address"),
		models.SymbolNEARToken,
	)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewUserBalanceResp(userBalance)})
}

func (s *Server) GetUserBalanceTransactions(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	currencyID, _ := s.uintFromContextQuery(c, "currency_id")
	userBalanceTnxs, count, err := s.nls.GetUserBalanceTransactions(
		ctx,
		models.Network(s.stringFromContextQuery(c, "network")),
		s.stringFromContextQuery(c, "address"),
		currencyID,
		s.stringFromContextQuery(c, "symbol"),
		page,
		limit,
	)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewUserBalanceTransactionRespArr(userBalanceTnxs), Count: &count})
}

func (s *Server) ClaimUserBalance(c *gin.Context) {
	ctx := s.requestContext(c)
	var req serializers.ClaimUserBalanceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	// err := s.validateTimestampWithSignature(
	// 	ctx,
	// 	req.Network,
	// 	req.Address,
	// 	req.Signature,
	// 	req.Timestamp,
	// )
	// if err != nil {
	// 	ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
	// 	return
	// }
	err := s.nls.ClaimUserBalance(ctx, &req)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}
