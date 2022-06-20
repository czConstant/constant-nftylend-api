package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) UserGetSettings(c *gin.Context) {
	ctx := s.requestContext(c)
	user, err := s.nls.UserGetSettings(ctx, models.Network(s.stringFromContextQuery(c, "network")), s.stringFromContextQuery(c, "address"))
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewUserResp(user)})
}

func (s *Server) UserUpdateSetting(c *gin.Context) {
	ctx := s.requestContext(c)
	var req serializers.UpdateUserSettingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	_, err := s.nls.UserUpdateSetting(ctx, &req)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}

func (s *Server) GetUserStats(c *gin.Context) {
	ctx := s.requestContext(c)
	borrowStats, lendStats, err := s.nls.GetUserStats(ctx, models.Network(s.stringFromContextQuery(c, "network")), s.stringFromContextQuery(c, "address"))
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: map[string]interface{}{
		"borrow_stats": borrowStats,
		"lend_stats":   lendStats,
	}})
}

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
	err := s.nls.ClaimUserBalance(ctx, &req)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}
