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

func (s *Server) GetUserPWPTokenBalance(c *gin.Context) {
	ctx := s.requestContext(c)
	userBalance, err := s.nls.GetUserPWPTokenBalance(ctx, models.Network(s.stringFromContextQuery(c, "network")), s.stringFromContextQuery(c, "address"))
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewUserBalanceResp(userBalance)})
}

func (s *Server) WithdrawUserBalance(c *gin.Context) {
	ctx := s.requestContext(c)
	var req serializers.WithdrawUserBalanceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	err := s.nls.WithdrawUserBalance(ctx, &req)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}
