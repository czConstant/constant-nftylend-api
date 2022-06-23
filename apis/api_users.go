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

func (s *Server) UserConnected(c *gin.Context) {
	ctx := s.requestContext(c)
	var req serializers.UserConnectedReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	err := s.validateTimestampWithSignature(
		ctx,
		req.Network,
		req.Address,
		req.Signature,
		req.Timestamp,
	)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	_, err = s.nls.UserConnected(ctx, req.Network, req.Address, req.ReferrerCode)
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
