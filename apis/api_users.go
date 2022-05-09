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
	user, err := s.nls.UserGetSettings(ctx, s.stringFromContextQuery(c, "address"), models.Network(s.stringFromContextQuery(c, "network")))
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewUserResp(user)})
}

func (s *Server) UserSettingEmail(c *gin.Context) {
	ctx := s.requestContext(c)
	var req struct {
		Address string         `json:"address"`
		Network models.Network `json:"network"`
		Email   string         `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	_, err := s.nls.UserSettingEmail(ctx, req.Address, req.Network, req.Email)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}
